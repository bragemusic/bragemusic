package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func generateConfigDocs(cfg any) (string, error) {
	t := reflect.TypeOf(cfg).Elem()
	v := reflect.ValueOf(cfg).Elem()

	if t.Kind() != reflect.Struct {
		return "", errors.New("input must be a struct")
	}

	strs, err := reflectDocsStruct(t, v, nil, nil)
	if err != nil {
		return "", err
	}

	return strings.Join(strs, "\n"), nil
}

func reflectDocsStruct(t reflect.Type, v reflect.Value, p []string, res []string) ([]string, error) {
	if len(p) > 0 {
		res = append(res, fmt.Sprintf("### %s", strings.Join(p, ".")))
		res = append(res, "| Parameter | Type | ENV Variable | Descripion |")
		res = append(res, "| --------- | ---- | ------------ | ---------- |")
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		fres, err := reflectDocsField(field, value, p)
		if err != nil {
			return nil, err
		}
		res = append(res, fres...)
	}
	return res, nil
}

func reflectDocsField(field reflect.StructField, value reflect.Value, p []string) ([]string, error) {
	if field.Tag.Get("toml") == "" {
		return nil, fmt.Errorf("field '%s' does not have a toml tag", field.Name)
	}

	p = append(p, strings.ToUpper(field.Tag.Get("toml")))
	if value.Kind() == reflect.Map {
		mapKeys := value.MapKeys()
		if len(mapKeys) > 0 {
			p = append(p, mapKeys[0].String())
			return reflectDocsStruct(value.MapIndex(mapKeys[0]).Type(), value.MapIndex(mapKeys[0]), p, nil)
		}
	}
	if value.Kind() == reflect.Struct {
		return reflectDocsStruct(field.Type, value, p, nil)
	}

	envName := fmt.Sprintf("%s_%s", envPrefix, strings.Join(p, "_"))

	res := []string{}
	res = append(res, fmt.Sprintf("| %s | %s | %s | %s |", field.Tag.Get("toml"), value.Kind(), envName, field.Tag.Get("desc")))

	return res, nil
}
