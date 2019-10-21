package main

import (
	"fmt"
	"syreclabs.com/go/faker"
	"testing"
)

var (
	usersQuery              = "INSERT INTO `wp_users` VALUES (1,'username','user_pass','username','hosting@humanmade.com','','2019-06-12 00:59:19','',0,'username'),(2,'username','user_pass','username','hosting@humanmade.com','http://notreal.com/username','2019-06-12 00:59:19','',0,'username');"
	usersQueryRecompiled    = "insert into wp_users values (1, 'treva_cremin', 'NjaK5HeMAMuv', 'hailey', 'bernice.heaney@example.net', '', '2019-06-12 00:59:19', '', 0, 'Kylie Rice'), (2, 'eduardo', 'J3JRQ4XoIxXX6A', 'albert.okeefe', 'brooke.hayes@example.net', 'http://pfannerstill.net/brando', '2019-06-12 00:59:19', '', 0, 'Ardella Jenkins PhD');"
	userMetaQuery           = "INSERT INTO `wp_usermeta` VALUES (1,1,'first_name','John'),(2,1,'last_name','Doe'),(3,1,'foobar','bazquz'),(4,1,'nickname','Jim'),(5,1,'description','Lorum ipsum.'),(6,2,'first_name','Janet'),(7,2,'last_name','Doe'),(8,2,'foobar','bazquz'),(9,2,'nickname','Jane'),(10,2,'description','Lorum ipsum.');"
	userMetaQueryRecompiled = "insert into wp_usermeta values (1, 1, 'first_name', 'Stephania'), (2, 1, 'last_name', 'Hamill'), (3, 1, 'foobar', 'bazquz'), (4, 1, 'nickname', 'Noah'), (5, 1, 'description', 'Dolorum nostrum alias.'), (6, 2, 'first_name', 'Ed'), (7, 2, 'last_name', 'Koelpin'), (8, 2, 'foobar', 'bazquz'), (9, 2, 'nickname', 'Watson'), (10, 2, 'description', 'Qui voluptatum est.');"
	commentsQuery           = "INSERT INTO `wp_comments` VALUES (1,1,'A WordPress Commenter','wapuu@wordpress.example','https://wordpress.org/','','2019-06-12 00:59:19','2019-06-12 00:59:19','Hi, this is a comment.\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\nCommenter avatars come from <a href=\"https://gravatar.com\">Gravatar</a>.',0,'1','','',0,0);"
	// Don't forget to escape \ because it'll translate to a newline and not pass
	// the comparison test
	commentsQueryRecompiled = "insert into wp_comments values (1, 1, 'kamren.ohara', 'michele_barton@example.net', 'http://ebert.com/korey_keeling', '', '2019-06-12 00:59:19', '2019-06-12 00:59:19', 'Hi, this is a comment.\\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\\nCommenter avatars come from <a href=\\\"https://gravatar.com\\\">Gravatar</a>.', 0, '1', '', '', 0, 0);"
)

func init() {
	faker.Seed(432)
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
			line := processLine(test.line)
			fmt.Printf("%+#v\n", line)
			if line != test.wants {
				t.Error("\nExpected:\n", test.wants, "\nActual:\n", line)
			}
		})
	}
}
