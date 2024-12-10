package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"unicode/utf8"
)

var (
	ErrParsing  = errors.New("error parsing env file")
	envFileName = ".env"
)

type Config struct {
	PORT              string
	TOKEN             string
	POSTGRES_HOST     string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_DB       string
	POSTGRES_PORT     string
}

var conf *Config

func New() Config {
	if conf == nil {
		contents, err := os.ReadFile(envFileName)
		if err != nil {
			log.Panic(errors.Join(ErrParsing, err))
		}
		r, _ := utf8.DecodeLastRune(contents)
		if r != '\n' {
			contents = append(contents, byte('\n'))
		}
		err = parse(contents)
		if err != nil {
			log.Panic(errors.Join(ErrParsing, err))
		}
	}
	return *conf
}

func parse(envFile []byte) error {
	vars := make(map[string]string)
	cursor, line, pos := 0, 0, 0
	key, val := "", ""
	onKey := true
	for {
		r, size := utf8.DecodeRune(envFile[cursor:])
		if r == utf8.RuneError && size == 0 {
			break
		}
		if r == utf8.RuneError && size == 1 {
			return fmt.Errorf("invalid character at line %d character %d", line, pos)
		}
		cursor += size
		pos += size
		if r == '\n' {
			line++
			vars[key] = val
			key, val = "", ""
			onKey = true
			pos = 0
			continue
		}
		if r == '=' {
			onKey = false
			continue
		}
		if onKey {
			key += string(r)
		} else {
			val += string(r)
		}
	}
	var err error = nil
	conf, err = newConf(vars)
	if err != nil {
		return err
	}
	return nil
}

func newConf(vars map[string]string) (*Config, error) {
	confTyp := reflect.TypeOf(Config{})
	confValue := reflect.New(confTyp)
	confValue = confValue.Elem()
	confNumField := confValue.NumField()

	for i := 0; i < confNumField; i++ {
		confField := confTyp.Field(i)
		key := confField.Name
		val, ok := vars[key]
		if !ok || key == "" {
			return nil, fmt.Errorf("variable %q is not set", key)
		}
		confValue.FieldByName(key).SetString(val)
	}

	v := confValue.Interface().(Config)
	return &v, nil
}
