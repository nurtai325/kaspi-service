package config

import (
	"fmt"
	"testing"
)

func TestNewConf(t *testing.T) {
	vars := map[string]string{
		"DATABASE_URL": "database_url",
		"PORT":         "8080",
	}
	conf, err := newConf(vars)
	if err != nil {
		t.Error(err)
	}
	if conf == nil {
		t.Errorf("config is nil")
		t.FailNow()
	}
	actual := conf.DATABASE_URL
	expected := vars["DATABASE_URL"]
	if actual != expected {
		t.Errorf("DATABASE_URL expected: %s, actual: %s", expected, actual)
	}
	actual = conf.PORT
	expected = vars["PORT"]
	if actual != expected {
		t.Errorf("PORT expected: %s, actual: %s", expected, actual)
	}
}

func TestParse(t *testing.T) {
	DATABASE_URL := "database_url"
	PORT := "8080"

	testEnvFile := ""
	testEnvFile += fmt.Sprintf("DATABASE_URL=%s\n", DATABASE_URL)
	testEnvFile += fmt.Sprintf("PORT=%s\n", PORT)

	err := parse([]byte(testEnvFile))
	if err != nil {
		t.Error(err)
	}
}
