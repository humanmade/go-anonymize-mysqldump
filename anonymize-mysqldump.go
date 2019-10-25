package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/sirupsen/logrus"
	"github.com/xwb1989/sqlparser"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type Config struct {
	Patterns []ConfigPattern `json:"patterns"`
}

type ConfigPattern struct {
	TableName string         `json:"tableName"`
	Fields    []PatternField `json:"fields"`
}

type PatternField struct {
	Field       string                   `json:"field"`
	Position    int                      `json:"position"`
	Type        string                   `json:"type"`
	Constraints []PatternFieldConstraint `json:"constraints"`
}

type PatternFieldConstraint struct {
	Field    string `json:"field"`
	Position int    `json:"position"`
	Value    string `json:"value"`
}

var (
	transformationFunctionMap = map[string]func(*sqlparser.SQLVal) *sqlparser.SQLVal{
		"username":  generateUsername,
		"password":  generatePassword,
		"email":     generateEmail,
		"url":       generateURL,
		"name":      generateName,
		"firstName": generateFirstName,
		"lastName":  generateLastName,
		"paragraph": generateParagraph,
		"ipv4":      generateIPv4,
	}
)

// Many thanks to https://stackoverflow.com/a/47515580/1454045
func init() {
	lvl, ok := os.LookupEnv("LOG_LEVEL")
	// LOG_LEVEL not set, let's default to info
	if !ok {
		lvl = "info"
	}
	// parse string, this is built-in feature of logrus
	ll, err := logrus.ParseLevel(lvl)
	if err != nil {
		ll = logrus.InfoLevel
	}
	// set global log level
	logrus.SetLevel(ll)
}

func main() {
	config := parseArgs()

	lines := setupAndProcessInput(config, os.Stdin)

	for line := range lines {
		fmt.Print(<-line)
	}
}

func setupAndProcessInput(config Config, input io.Reader) chan chan string {
	var wg sync.WaitGroup
	lines := make(chan chan string, 10)

	wg.Add(1)
	go processInput(&wg, input, lines, config)

	go func() {
		wg.Wait()
		close(lines)
	}()

	return lines
}

func parseArgs() Config {
	parser := argparse.NewParser("anonymize-mysqldump", "Reads SQL from STDIN and replaces content for anonymity based on the provided config.")
	configFilePath := parser.String("c", "config", &argparse.Options{Required: true, Help: "Path to config.json"})

	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	return readConfigFile(*configFilePath)
}

func readConfigFile(filepath string) Config {
	jsonConfig, err := ioutil.ReadFile(filepath)
	if err != nil {
		logrus.Fatal(err)
	}

	var decoded Config
	jsonReader := strings.NewReader(string(jsonConfig))
	jsonParser := json.NewDecoder(jsonReader)
	jsonParser.Decode(&decoded)
	return decoded
}

func processInput(wg *sync.WaitGroup, input io.Reader, lines chan chan string, config Config) {
	defer wg.Done()

	r := bufio.NewReaderSize(input, 2*1024*1024)
	var nextLine string
	insertStarted := false
	continueLooping := true
	for continueLooping {
		line, err := r.ReadString('\n')

		if err == io.EOF {
			// continueLooping is used because line might be populated even when we've
			// reached the end of the file, so we set a boolean once the last line is
			// being processed to end the loop.
			continueLooping = false
		} else if err != nil {
			// log any other errors and break
			logrus.Error(err.Error())
			break
		}

		// If the line is shorter than 6 characters, which is the shortest line for
		// an insert query, let's skip processing it
		if len(line) < 6 {

			// TODO I'd love to clean this up so we don't make ch in three different
			// places, but that's a task for another day
			ch := make(chan string)
			lines <- ch
			ch <- line
			//ch <- line + "\n"
			continue
		}

		// Test if this is an INSERT query. We'll use this to determine if we need
		// to concatenate lines together if they're spread apart multiple lines
		// instead of on a single line
		maybeInsert := strings.ToUpper(line[:6]) == "INSERT"
		if maybeInsert {
			insertStarted = true
		}

		line = strings.TrimSpace(line)
		// Now that we've detected this is an INSERT query, let's append the lines
		// together to form a single line in the event this spans multiple lines
		if insertStarted {
			nextLine += line
		} else {
			// When it's not an insert query, let's add this line and move on without
			// processing it
			// TODO clean this up too
			ch := make(chan string)
			lines <- ch
			ch <- line + "\n"
			continue
		}

		lastCharacter := line[len(line)-1:]
		if lastCharacter == ";" {
			insertStarted = false
		} else {
			// If we haven't reached a query terminator and and insert query has
			// begun, let's move on to the next line
			continue
		}

		// Now let's actually process the line!
		wg.Add(1)
		ch := make(chan string)
		lines <- ch
		go func(line string) {
			defer wg.Done()
			line = processLine(line, config)
			ch <- line
		}(nextLine)

		// Now let's reset nextLine to empty so that it doesn't continue
		// appending lines forever
		nextLine = ""
	}

}

func processLine(line string, config Config) string {

	parsed, err := parseLine(line)
	if err != nil {
		// TODO Add line number to log
		logrus.WithFields(logrus.Fields{
			"error": err,
			"line":  line,
		}).Error("Failed parsing line with error: ")
		return line
	}

	// TODO Detect if line matches pattern
	processed, err := applyConfigToParsedLine(parsed, config)
	// TODO make modifications

	// TODO Return changes
	recompiled, err := recompileStatementToSQL(processed)
	if err != nil {
		// TODO Add line number to log
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed recompiling line with error: ")
		return line
	}
	return recompiled
}

func parseLine(line string) (sqlparser.Statement, error) {
	stmt, err := sqlparser.Parse(line)
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func applyConfigToParsedLine(stmt sqlparser.Statement, config Config) (sqlparser.Statement, error) {

	insert, isInsertStatement := stmt.(*sqlparser.Insert)
	if !isInsertStatement {
		// Let's skip other statements as we only want to process inserts.
		return stmt, nil
	}

	modified, err := applyConfigToInserts(insert, config)
	if err != nil {
		// TODO Log error and move on
		return stmt, nil
	}
	return modified, nil
}

func applyConfigToInserts(stmt *sqlparser.Insert, config Config) (*sqlparser.Insert, error) {

	values, isValuesSlice := stmt.Rows.(sqlparser.Values)
	if !isValuesSlice {
		// This _should_ have type Values, but if it doesn't, let's skip it
		// TODO Perhaps worth logging when this happens?
		return stmt, nil
	}

	// Iterate over the specified configs and see if this statement matches any
	// of the desired changes
	// TODO make this use goroutines
	for _, pattern := range config.Patterns {
		if stmt.Table.Name.String() != pattern.TableName {
			// Config is not for this table, move onto next available config
			continue
		}

		// Ok, now it's time to make some modifications
		newValues, err := modifyValues(values, pattern)
		if err != nil {
			// TODO Perhaps worth logging when this happens?
			return stmt, nil
		}
		stmt.Rows = newValues
	}

	return stmt, nil
}

// TODO we're gonna have to figure out how to retain types if we ever want to
// mask number-based fields
func modifyValues(values sqlparser.Values, pattern ConfigPattern) (sqlparser.Values, error) {

	// TODO make this use goroutines
	for row := range values {
		// TODO make this use goroutines
		for _, fieldPattern := range pattern.Fields {
			// Position is 1 indexed instead of 0, so let's subtract 1 in order to get
			// it to line up with the value inside the ValTuple inside of values.Values
			valTupleIndex := fieldPattern.Position - 1
			value := values[row][valTupleIndex].(*sqlparser.SQLVal)

			// Skip transformation if transforming function doesn't exist
			if transformationFunctionMap[fieldPattern.Type] == nil {
				// TODO in the event a transformation function isn't correctly defined,
				// should we actually exit? Should we exit or fail softly whenever
				// something goes wrong in general?
				logrus.WithFields(logrus.Fields{
					"type":  fieldPattern.Type,
					"field": fieldPattern.Field,
				}).Error("Failed applying transformation type for field")
				continue
			}

			// Skipping applying a transformation because field is empty
			if len(value.Val) == 0 {
				continue
			}

			// Skip this PatternField if none of its constraints match
			if fieldPattern.Constraints != nil && !rowObeysConstraints(fieldPattern.Constraints, values[row]) {
				continue
			}

			values[row][valTupleIndex] = transformationFunctionMap[fieldPattern.Type](value)
		}

	}

	// values[0][0] = sqlparser.NewStrVal([]byte("Foobar"))
	return values, nil
}

func rowObeysConstraints(constraints []PatternFieldConstraint, row sqlparser.ValTuple) bool {
	for _, constraint := range constraints {
		valTupleIndex := constraint.Position - 1
		value := row[valTupleIndex].(*sqlparser.SQLVal)

		parsedValue := convertSQLValToString(value)
		logrus.WithFields(logrus.Fields{
			"parsedValue":      parsedValue,
			"constraint.value": constraint.Value,
		}).Trace("Debuging constraint obediance: ")
		if parsedValue != constraint.Value {
			return false
		}
	}
	return true
}

func convertSQLValToString(value *sqlparser.SQLVal) string {
	buf := sqlparser.NewTrackedBuffer(nil)
	buf.Myprintf("%s", []byte(value.Val))
	pq := buf.ParsedQuery()

	bytes, err := pq.GenerateQuery(nil, nil)
	if err != nil {
		return ""
	}
	return string(bytes)
}
func recompileStatementToSQL(stmt sqlparser.Statement) (string, error) {
	// TODO Potentially replace with BuildParsedQuery
	buf := sqlparser.NewTrackedBuffer(nil)
	buf.Myprintf("%v", stmt)
	pq := buf.ParsedQuery()

	bytes, err := pq.GenerateQuery(nil, nil)
	if err != nil {
		return "", err
	}
	return string(bytes) + ";\n", nil
}
