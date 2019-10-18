package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	Field    string
	Type     string
	Position int
}

var (
	WordPressConfig = Config{
		Patterns: []ConfigPattern{
			{
				TableName: "wp_users",
				Fields: []PatternField{
					{
						Field:    "user_login",
						Type:     "username",
						Position: 2,
					},
					{
						Field:    "user_pass",
						Type:     "password",
						Position: 3,
					},
					{
						Field:    "user_nicename",
						Type:     "username",
						Position: 4,
					},
					{
						Field:    "user_email",
						Type:     "email",
						Position: 5,
					},
					{
						Field:    "user_url",
						Type:     "url",
						Position: 6,
					},
					{
						Field:    "display_name",
						Type:     "name",
						Position: 10,
					},
				},
			},
			// TODO Hmm, usermeta is going to be a challenge because there's only one
			// colum we want to change, but it requires knowledge of another field to
			// trigger a modification 🤔
			// {
			// TableName: "wp_usermeta",
			// Fields: []MultiPatternField{
			// {
			// Field:    "meta_value",
			// Type:     "multi",
			// Position: 4,
			// Patterns: []MultiPattern{
			// {
			// }
			// },
			// },
			// },
			// },
		},
	}
)

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
			fmt.Fprintln(os.Stderr, err.Error())
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
		fmt.Fprintf(os.Stderr, "Failed parsing line with error: %v\n", err)
		return line
	}

	// TODO Detect if line matches pattern
	processed, err := applyConfigToParsedLine(parsed, WordPressConfig)
	// TODO make modifications

	// TODO Return changes
	recompiled, err := recompileStatementToSQL(processed)
	if err != nil {
		// TODO Add line number to log
		fmt.Fprintf(os.Stderr, "Failed recompiling line with error: %v\n", err)
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

	// We have to use a switch here to avoid compile-time errors
	switch s := stmt.(type) {
	case *sqlparser.Insert:
		modified, err := applyConfigToInserts(s, config)
		if err != nil {
			return nil, err
		}
		return modified, nil

	default:
		// ignore all other statements
		return stmt, nil
	}
}

func applyConfigToInserts(stmt *sqlparser.Insert, config Config) (*sqlparser.Insert, error) {

	if values, ok := stmt.Rows.(sqlparser.Values); ok {
		newValues, err := modifyValues(values, config)
		if err != nil {
			return nil, err
		}

		stmt.Rows = newValues
	}

	return stmt, nil
}

// TODO we're gonna have to figure out how to retain types if we ever want to
// mask number-based fields
func modifyValues(values sqlparser.Values, config Config) (sqlparser.Values, error) {
	fmt.Printf("%+#v\n", values)
	for _, row := range values {
		for _, n := range row {
			if value, ok := n.(*sqlparser.SQLVal); ok {
				switch value.Type {
				case sqlparser.IntVal, sqlparser.StrVal:
					fmt.Printf("%+#v\n", string(value.Val))
				default:
				}

				// fmt.Printf("%+#v\n", string(value.Val))
			}
		}
	}
	return values, nil
}

// TODO update to have query include bound variables
// TODO add replacements to bound variables
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
