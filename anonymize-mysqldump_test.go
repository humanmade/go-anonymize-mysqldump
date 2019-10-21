package main

import (
	"fmt"
	"testing"
)

var (
	usersQuery              = "INSERT INTO `wp_users` VALUES (1,'username','user_pass','username','hosting@humanmade.com','','2019-06-12 00:59:19','',0,'username'),(2,'username','user_pass','username','hosting@humanmade.com','http://notreal.com/username','2019-06-12 00:59:19','',0,'username');"
	usersQueryRecompiled    = "insert into wp_users values (1, 'foobar', 'foobar', 'foobar', 'foobar@example.com', '', '2019-06-12 00:59:19', '', 0, 'Ashley Jones'), (2, 'foobar', 'foobar', 'foobar', 'foobar@example.com', 'https://example.com', '2019-06-12 00:59:19', '', 0, 'Ashley Jones');"
	userMetaQuery           = "INSERT INTO `wp_usermeta` VALUES (1,1,'first_name','John'),(2,1,'last_name','Doe'),(3,1,'foobar','bazquz'),(4,1,'nickname','Jim'),(5,1,'description','Lorum ipsum.'),(6,2,'first_name','Janet'),(7,2,'last_name','Doe'),(8,2,'foobar','bazquz'),(9,2,'nickname','Jane'),(10,2,'description','Lorum ipsum.');"
	userMetaQueryRecompiled = "insert into wp_usermeta values (1, 1, 'first_name', 'Ashley'), (2, 1, 'last_name', 'Jones'), (3, 1, 'foobar', 'bazquz'), (4, 1, 'nickname', 'Ashley'), (5, 1, 'description', 'Foo bar baz quz.'), (6, 2, 'first_name', 'Ashley'), (7, 2, 'last_name', 'Jones'), (8, 2, 'foobar', 'bazquz'), (9, 2, 'nickname', 'Ashley'), (10, 2, 'description', 'Foo bar baz quz.');"
	commentsQuery           = "INSERT INTO `wp_comments` VALUES (1,1,'A WordPress Commenter','wapuu@wordpress.example','https://wordpress.org/','','2019-06-12 00:59:19','2019-06-12 00:59:19','Hi, this is a comment.\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\nCommenter avatars come from <a href=\"https://gravatar.com\">Gravatar</a>.',0,'1','','',0,0);"
	commentsQueryRecompiled = "insert into wp_comments values (1, 1, 'foobar', 'foobar@example.com', 'https://example.com', '', '2019-06-12 00:59:19', '2019-06-12 00:59:19', 'Hi, this is a comment.\\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\\nCommenter avatars come from <a href=\\\"https://gravatar.com\\\">Gravatar</a>.', 0, '1', '', '', 0, 0);"
)

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
