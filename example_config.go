package main

var ExampleWordPressConfig = Config{
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
