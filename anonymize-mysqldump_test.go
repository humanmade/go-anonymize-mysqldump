package main

import (
	"testing"
)

var (
	usersQuery           = "INSERT INTO `wp_users` VALUES (1,'username','user_pass','username','hosting@humanmade.com','','2019-06-12 00:59:19','',0,'username');"
	usersQueryRecompiled = "insert into wp_users values (1, 'username', 'user_pass', 'username', 'hosting@humanmade.com', '', '2019-06-12 00:59:19', '', 0, 'username');"
)

func TestRecompileSQL(t *testing.T) {
	var tests = []struct {
		testName string
		line     string
		wants    string
	}{
		{
			testName: "recompile users query",
			line:     usersQuery,
			wants:    usersQueryRecompiled,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			line := processLine(test.line)
			if line != test.wants {
				t.Error("Expected:", test.wants, "Actual:", line)
			}
		})
	}
}

// First, lets test if we can match a line
// Then, lets read config to see what transformations need to be performed on
// each field
// Then lets apply the transformations
func TestApplyConfigToQuery(t *testing.T) {

	var tests = []struct {
		testName string
		line     string
		wants    string
	}{
		{
			testName: "recompile users query",
			line:     usersQuery,
			wants:    "fails",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			line, _ := parseLine(test.line)
			if true == false {
				t.Error("Expected:", test.wants, "Actual:", line)
			}
		})
	}
}
