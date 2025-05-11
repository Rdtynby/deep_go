package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"omitempty,age"`
	Married bool   `properties:"married"`
}

func Serialize(person interface{}) string {
	personType := reflect.TypeOf(person)
	personValue := reflect.ValueOf(person)

	var result []string

outer:
	for i := 0; i < personType.NumField(); i++ {
		props := strings.Split(personType.Field(i).Tag.Get("properties"), ",")
		field := personValue.Field(i)

		isZero := reflect.Zero(field.Type()).Interface() == field.Interface()

		nameIndex := 0

		for index, prop := range props {
			if prop == "omitempty" {
				if isZero {
					continue outer
				} else {
					if index == nameIndex {
						nameIndex++
					}
				}
			}
		}

		result = append(result, props[nameIndex]+"="+fmt.Sprintf("%v", field.Interface()))
	}

	return strings.Join(result, "\n")
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
