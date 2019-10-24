package main

import (
	"bytes"
	"syreclabs.com/go/faker"
	"testing"
)

var (
	jsonConfig     Config
	multilineQuery = `INSERT INTO wp_usermeta VALUES
	(1,1,'first_name','John'),(2,1,'last_name','Doe'),
	(3,1,'foobar','bazquz'),
	(4,1,'nickname','Jim'),
	(5,1,'description','Lorum ipsum.');
`
	multilineQueryRecompiled = "insert into wp_usermeta values (1, 1, 'first_name', 'Nat'), (2, 1, 'last_name', 'Hermiston'), (3, 1, 'foobar', 'bazquz'), (4, 1, 'nickname', 'Treva'), (5, 1, 'description', 'Enim odio nihil.');\n"
	commentsQuery            = "INSERT INTO `wp_comments` VALUES (1,1,'A WordPress Commenter','wapuu@wordpress.example','https://wordpress.org/','','2019-06-12 00:59:19','2019-06-12 00:59:19','Hi, this is a comment.\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\nCommenter avatars come from <a href=\"https://gravatar.com\">Gravatar</a>.',0,'1','','',0,0);\n"
	// Don't forget to escape \ because it'll translate to a newline and not pass
	// the comparison test
	commentsQueryRecompiled = "insert into wp_comments values (1, 1, 'sam_harvey', 'jillian@example.com', 'http://balistreriwiegand.name/sunny', '', '2019-06-12 00:59:19', '2019-06-12 00:59:19', 'Hi, this is a comment.\\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\\nCommenter avatars come from <a href=\\\"https://gravatar.com\\\">Gravatar</a>.', 0, '1', '', '', 0, 0);\n"
	usersQuery              = "INSERT INTO `wp_users` VALUES (1,'username','user_pass','username','hosting@humanmade.com','','2019-06-12 00:59:19','',0,'username'),(2,'username','user_pass','username','hosting@humanmade.com','http://notreal.com/username','2019-06-12 00:59:19','',0,'username');\n"
	usersQueryRecompiled    = "insert into wp_users values (1, 'fatima.fisher', 'abOSwkVS', 'lillian', 'grover@example.net', '', '2019-06-12 00:59:19', '', 0, 'Retta Bailey'), (2, 'juwan.kassulke', 'zgtEQA3nm4Wlro', 'evalyn', 'camilla.hilll@example.org', 'http://dickensmurphy.info/ophelia', '2019-06-12 00:59:19', '', 0, 'Rick Fahey III');\n"
	userMetaQuery           = "INSERT INTO `wp_usermeta` VALUES (1,1,'first_name','John'),(2,1,'last_name','Doe'),(3,1,'foobar','bazquz'),(4,1,'nickname','Jim'),(5,1,'description','Lorum ipsum.'),(6,2,'first_name','Janet'),(7,2,'last_name','Doe'),(8,2,'foobar','bazquz'),(9,2,'nickname','Jane'),(10,2,'description','Lorum ipsum.');\n"
	userMetaQueryRecompiled = "insert into wp_usermeta values (1, 1, 'first_name', 'Ed'), (2, 1, 'last_name', 'Koelpin'), (3, 1, 'foobar', 'bazquz'), (4, 1, 'nickname', 'Watson'), (5, 1, 'description', 'Qui voluptatum est.'), (6, 2, 'first_name', 'Olen'), (7, 2, 'last_name', 'Williamson'), (8, 2, 'foobar', 'bazquz'), (9, 2, 'nickname', 'Kamren'), (10, 2, 'description', 'Eveniet repellat in.');\n"
)

func init() {
	faker.Seed(432)
	jsonConfig = readConfigFile("./config.example.json")
}

func BenchmarkProcessLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processLine(usersQuery, jsonConfig)
		processLine(userMetaQuery, jsonConfig)
		processLine(commentsQuery, jsonConfig)
	}
}

func TestProcessFile(t *testing.T) {
	input := bytes.NewBufferString(multilineQuery)

	lines := setupAndProcessInput(jsonConfig, input)

	var result string
	for line := range lines {
		result = <-line
	}

	if result != multilineQueryRecompiled {
		t.Error("\nExpected:\n", multilineQueryRecompiled, "\nActual:\n", result)
	}
}

func TestApplyConfigToQuery(t *testing.T) {

	var tests = []struct {
		testName string
		line     string
		wants    string
	}{
		{
			testName: "users query",
			line:     usersQuery,
			wants:    usersQueryRecompiled,
		},
		{
			testName: "usermeta query",
			line:     userMetaQuery,
			wants:    userMetaQueryRecompiled,
		},
		{
			testName: "comments query",
			line:     commentsQuery,
			wants:    commentsQueryRecompiled,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			line := processLine(test.line, jsonConfig)
			if line != test.wants {
				t.Error("\nExpected:\n", test.wants, "\nActual:\n", line)
			}
		})
	}
}
