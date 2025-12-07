package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const envPrefix = "BRAGE"

func AddFromEnv(cfg any) error {
	t := reflect.TypeOf(cfg).Elem()
	v := reflect.ValueOf(cfg).Elem()

	if t.Kind() != reflect.Struct {
		return errors.New("input must be a struct")
	}

	return reflectStruct(t, v, nil)
}

func reflectStruct(t reflect.Type, v reflect.Value, p []string) error {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		reflectField(field, value, p)
	}
	return nil
}

func reflectField(field reflect.StructField, value reflect.Value, p []string) error {
	if field.Tag.Get("toml") == "" {
		return fmt.Errorf("field '%s' does not have a toml tag", field.Name)
	}

	p = append(p, strings.ToUpper(field.Tag.Get("toml")))

	if value.Kind() == reflect.Struct {
		return reflectStruct(field.Type, value, p)
	}

	envName := fmt.Sprintf("%s_%s", envPrefix, strings.Join(p, "_"))
	envVal := os.Getenv(envName)

	if envVal == "" {
		return nil
	}

	switch value.Kind() {
	case reflect.String:
		value.SetString(envVal)
	case reflect.Int:
		val, err := strconv.Atoi(envVal)
		if err != nil {
			return err
		}
		value.SetInt(int64(val))
	case reflect.Float32:
		val, err := strconv.ParseFloat(envVal, 32)
		if err != nil {
			return err
		}
		value.SetFloat(val)
	default:
		return fmt.Errorf("field kind '%s' is not supported ", value.Kind())
	}

	return nil
}
