package main

import (
	"github.com/xwb1989/sqlparser"
	"syreclabs.com/go/faker"
)

func generateUsername(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().UserName()))
}

func generatePassword(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	// TODO encrypt this value
	return sqlparser.NewStrVal([]byte(faker.Internet().Password(8, 14)))
}

func generateEmail(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().SafeEmail()))
}

func generateURL(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().Url()))
}

func generateName(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Name().Name()))
}

func generateFirstName(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Name().FirstName()))
}

func generateLastName(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Name().LastName()))
}

func generateParagraph(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Lorem().Sentence(3)))
}

func generateIPv4(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().IpV4Address()))
}
