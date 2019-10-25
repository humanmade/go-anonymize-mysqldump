package main

import (
	"bytes"
	"syreclabs.com/go/faker"
	"testing"
)

var (
	jsonConfig         Config
	dropAndCreateTable = "DROP TABLE IF EXISTS `wp_options`;\n" +
		"/*!40101 SET @saved_cs_client     = @@character_set_client */;\n" +
		"/*!40101 SET character_set_client = utf8 */;\n" +
		"CREATE TABLE `wp_options` (\n" +
		"`option_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,\n" +
		"`option_name` varchar(191) NOT NULL DEFAULT '',\n" +
		"`option_value` longtext NOT NULL,\n" +
		"`autoload` varchar(20) NOT NULL DEFAULT 'yes',\n" +
		"PRIMARY KEY (`option_id`),\n" +
		"UNIQUE KEY `option_name` (`option_name`)\n" +
		") ENGINE=InnoDB AUTO_INCREMENT=123 DEFAULT CHARSET=utf8mb4;\n" +
		"/*!40101 SET character_set_client = @saved_cs_client */;"

	// Don't forget to escape \ because it'll translate to a newline and not pass
	// the comparison test
	multilineQuery = `INSERT INTO wp_usermeta VALUES
	(1,1,'first_name','John'),(2,1,'last_name','Doe'),
	(3,1,'foobar','bazquz'),
	(4,1,'nickname','Jim'),
	(5,1,'description','Lorum ipsum.');`
	multilineQueryRecompiled = "insert into wp_usermeta values (1, 1, 'first_name', 'Jazmyn'), (2, 1, 'last_name', 'Reynolds'), (3, 1, 'foobar', 'bazquz'), (4, 1, 'nickname', 'Sherman'), (5, 1, 'description', 'Vel at et.');\n"
	commentsQuery            = "INSERT INTO `wp_comments` VALUES (1,1,'A WordPress Commenter','wapuu@wordpress.example','https://wordpress.org/','','2019-06-12 00:59:19','2019-06-12 00:59:19','Hi, this is a comment.\\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\\nCommenter avatars come from <a href=\\\"https://gravatar.com\\\">Gravatar</a>.',0,'1','','',0,0);\n"
	commentsQueryRecompiled  = "insert into wp_comments values (1, 1, 'kamren.ohara', 'michele_barton@example.net', 'http://ebert.com/korey_keeling', '', '2019-06-12 00:59:19', '2019-06-12 00:59:19', 'Hi, this is a comment.\\nTo get started with moderating, editing, and deleting comments, please visit the Comments screen in the dashboard.\\nCommenter avatars come from <a href=\\\"https://gravatar.com\\\">Gravatar</a>.', 0, '1', '', '', 0, 0);\n"
	usersQuery               = "INSERT INTO `wp_users` VALUES (1,'username','user_pass','username','hosting@humanmade.com','','2019-06-12 00:59:19','',0,'username'),(2,'username','user_pass','username','hosting@humanmade.com','http://notreal.com/username','2019-06-12 00:59:19','',0,'username');\n"
	usersQueryRecompiled     = "insert into wp_users values (1, 'treva_cremin', 'NjaK5HeMAMuv', 'hailey', 'bernice.heaney@example.net', '', '2019-06-12 00:59:19', '', 0, 'Kylie Rice'), (2, 'eduardo', 'J3JRQ4XoIxXX6A', 'albert.okeefe', 'brooke.hayes@example.net', 'http://pfannerstill.net/brando', '2019-06-12 00:59:19', '', 0, 'Ardella Jenkins PhD');\n"
	userMetaQuery            = "INSERT INTO `wp_usermeta` VALUES (1,1,'first_name','John'),(2,1,'last_name','Doe'),(3,1,'foobar','bazquz'),(4,1,'nickname','Jim'),(5,1,'description','Lorum ipsum.'),(6,2,'first_name','Janet'),(7,2,'last_name','Doe'),(8,2,'foobar','bazquz'),(9,2,'nickname','Jane'),(10,2,'description','Lorum ipsum.');\n"
	userMetaQueryRecompiled  = "insert into wp_usermeta values (1, 1, 'first_name', 'Stephania'), (2, 1, 'last_name', 'Hamill'), (3, 1, 'foobar', 'bazquz'), (4, 1, 'nickname', 'Noah'), (5, 1, 'description', 'Dolorum nostrum alias.'), (6, 2, 'first_name', 'Ed'), (7, 2, 'last_name', 'Koelpin'), (8, 2, 'foobar', 'bazquz'), (9, 2, 'nickname', 'Watson'), (10, 2, 'description', 'Qui voluptatum est.');\n"
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

func TestSetupAndProcessInput(t *testing.T) {

	var tests = []struct {
		testName string
		query    string
		wants    string
	}{
		{
			testName: "users query",
			query:    usersQuery,
			wants:    usersQueryRecompiled,
		},
		{
			testName: "usermeta query",
			query:    userMetaQuery,
			wants:    userMetaQueryRecompiled,
		},
		{
			testName: "comments query",
			query:    commentsQuery,
			wants:    commentsQueryRecompiled,
		},
		{
			testName: "multiline query",
			query:    multilineQuery,
			wants:    multilineQueryRecompiled,
		},
		{
			testName: "table creation",
			query:    dropAndCreateTable,
			wants:    dropAndCreateTable + "\n",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {

			input := bytes.NewBufferString(test.query)
			lines := setupAndProcessInput(jsonConfig, input)

			var result string
			for line := range lines {
				result += <-line
			}

			if result != test.wants {
				t.Error("\nExpected:\n", test.wants, "\nActual:\n", result)
			}
		})
	}
}
