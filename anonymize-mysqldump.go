package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xwb1989/sqlparser"
	"io"
	"os"
	"strings"
	"sync"
)

type Config struct {
	// add something about tables and fields to process
	Patterns []ConfigPattern
}

type ConfigPattern struct {
	TableName string
	Fields    []PatternField
}

type PatternField struct {
	Field       string
	Position    int
	Type        string
	Constraints []PatternFieldConstraint
}

type PatternFieldConstraint struct {
	Field    string
	Position int
	Value    string
}

// TODO Should we add validation that enforces only one ConfigPattern per Table?
var (
	WordPressConfig = Config{
		Patterns: []ConfigPattern{
			{
				TableName: "wp_users",
				Fields: []PatternField{
					{
						Field:       "user_login",
						Type:        "username",
						Position:    2,
						Constraints: nil,
					},
					{
						Field:       "user_pass",
						Type:        "password",
						Position:    3,
						Constraints: nil,
					},
					{
						Field:       "user_nicename",
						Type:        "username",
						Position:    4,
						Constraints: nil,
					},
					{
						Field:       "user_email",
						Type:        "email",
						Position:    5,
						Constraints: nil,
					},
					{
						Field:       "user_url",
						Type:        "url",
						Position:    6,
						Constraints: nil,
					},
					{
						Field:       "display_name",
						Type:        "name",
						Position:    10,
						Constraints: nil,
					},
				},
			},
			{
				TableName: "wp_usermeta",
				Fields: []PatternField{
					{
						Field:    "meta_value",
						Position: 4,
						Type:     "firstName",
						Constraints: []PatternFieldConstraint{
							{
								Field:    "meta_key",
								Position: 3,
								Value:    "first_name",
							},
						},
					},
					{
						Field:    "meta_value",
						Position: 4,
						Type:     "lastName",
						Constraints: []PatternFieldConstraint{
							{
								Field:    "meta_key",
								Position: 3,
								Value:    "last_name",
							},
						},
					},
					{
						Field:    "meta_value",
						Position: 4,
						Type:     "firstName",
						Constraints: []PatternFieldConstraint{
							{
								Field:    "meta_key",
								Position: 3,
								Value:    "nickname",
							},
						},
					},
					{
						Field:    "meta_value",
						Position: 4,
						Type:     "paragraph",
						Constraints: []PatternFieldConstraint{
							{
								Field:    "meta_key",
								Position: 3,
								Value:    "description",
							},
						},
					},
				},
			},
			{
				TableName: "wp_comments",
				Fields: []PatternField{
					{
						Field:       "comment_author",
						Type:        "username",
						Position:    3,
						Constraints: nil,
					},
					{
						Field:       "comment_author_email",
						Type:        "email",
						Position:    4,
						Constraints: nil,
					},
					{
						Field:       "comment_author_url",
						Type:        "url",
						Position:    5,
						Constraints: nil,
					},
					{
						Field:       "comment_author_IP",
						Type:        "ipv4",
						Position:    6,
						Constraints: nil,
					},
				},
			},
		},
	}
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
	// if len(os.Args) < 2 {
	// fmt.Fprintln(os.Stderr, "Usage: anonymize-mysqldump <config>")
	// os.Exit(1)
	// return
	// }

	// config := loadConfiguration(os.Args[1])
	// fmt.Println(config)

	var wg sync.WaitGroup
	lines := make(chan chan string, 10)

	wg.Add(1)
	go processFile(&wg, lines)

	go func() {
		wg.Wait()
		close(lines)
	}()

	for line := range lines {
		fmt.Print(<-line)
	}
}

func loadConfiguration(jsonConfig string) Config {
	var decoded Config
	jsonReader := strings.NewReader(jsonConfig)
	jsonParser := json.NewDecoder(jsonReader)
	jsonParser.Decode(&decoded)
	return decoded
}

func processFile(wg *sync.WaitGroup, lines chan chan string) {
	defer wg.Done()

	r := bufio.NewReaderSize(os.Stdin, 2*1024*1024)
	for {
		line, err := r.ReadString('\n')

		if err == io.EOF {
			break
		}

		if err != nil {
			logrus.Error(err.Error())
			break
		}

		wg.Add(1)
		ch := make(chan string)
		lines <- ch

		go func(line string) {
			defer wg.Done()
			line = processLine(line)
			ch <- line
		}(line)
	}
}

func processLine(line string) string {

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
	processed, err := applyConfigToParsedLine(parsed, WordPressConfig)
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
		}).Debug("Debuging constraint obediance: ")
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
	return string(bytes) + ";", nil
}
