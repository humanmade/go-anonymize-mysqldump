package main

import (
	"fmt"
	"syreclabs.com/go/faker"
	"github.com/xwb1989/sqlparser"
	"math/rand"
	"strings"
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

func generateCustomTitle(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Lorem().Word()))
}

func generateIPv4(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().IpV4Address()))
}

func generateCustomPhoneNumber(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(fmt.Sprintf("0%v", faker.Number().Number(rand.Intn(12-7)+7))))
}

func generateCustomUserAgent(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9"))
}

func generateLibPrefix(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Name().Prefix()))
}

func generateLibCompanyName(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Company().Name()))
}

func generateLibStreet(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().StreetName()))
}

func generateCustomStreet(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().StreetName()))
}

func generateLibBuildingNumber(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().BuildingNumber()))
}

func generateLibAdditionalAddress(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().SecondaryAddress()))
}

func generateLibCity(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().City()))
}

func generateLibZip(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().ZipCode()))
}

func generateLibState(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().State()))
}

func generateLibCountry(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Address().Country()))
}

func generateCustomCountry(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte("Deutschland"))
}

func generateLibParagraph(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Lorem().Paragraph(3)))
}

func generateCustomSocialNetwork(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.RandomChoice([]string{"Google", "Facebook"})))
}

func generateCustomSocialId(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte("REDACTED"))
}

func generateCustomSocialToken(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte("REDACTED"))
}

func generateLibInternetUser(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().UserName()))
}

func generateCustomUniqueUser(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(fmt.Sprintf("%v.%v", faker.RandomString(16), faker.Internet().SafeEmail())))
}

func generateCustomPassword(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(strings.ToLower(faker.Internet().Password(64, 64))))
}

func generateCustomRecoverToken(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().Password(30, 30)))
}

func generateCustomUserToken(value *sqlparser.SQLVal) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(faker.Internet().Password(8, 20)))
}

func generateCustomString(value *sqlparser.SQLVal, string string) *sqlparser.SQLVal {
	return sqlparser.NewStrVal([]byte(string))
}