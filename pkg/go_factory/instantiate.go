package gofactory

import (
	"fmt"
	"reflect"

	"github.com/anti-duhring/easy-shipping-microsservice/pkg/go_factory/generate"
)

var typesMap = map[string]func(field *reflect.Value){
	"string": func(field *reflect.Value) {
		val := generate.String()
		field.SetString(val)
	},
	"float64": func(field *reflect.Value) {
		val := generate.Float64()
		field.SetFloat(val)
	},
	"Time": func(field *reflect.Value) {
		val := generate.Date()
		field.Set(reflect.ValueOf(val))
	},
}

func Instantiate(s interface{}, defaultValues interface{}) error {
	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("not a pointer to a struct")
	}

	// Dereference the pointer and get the struct value
	structValue := val.Elem()
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		fieldName := structType.Field(i).Name
		field := structValue.FieldByName(fieldName)
		defaultValue := reflect.ValueOf(defaultValues).Elem().FieldByName(fieldName)

		// Check if the field exists
		if !field.IsValid() {
			return fmt.Errorf("field %s does not exist in the struct", fieldName)
		}

		// Check if the field can be set (exported)
		if !field.CanSet() {
			return fmt.Errorf("field %s cannot be set", fieldName)
		}

		if defaultValue.IsValid() {
			field.Set(defaultValue)

			break
		}

		setField(field.Type().String(), fieldName, &field)
	}

	return nil
}

type SetFieldParams struct {
	FieldType string
	FieldName string
	Field     *reflect.Value
}

func setField(fieldType string, fieldName string, field *reflect.Value) error {
	setNewValue, ok := typesMap[fieldType]

	if !ok {
		return fmt.Errorf("type %s not supported", fieldType)
	}

	setNewValue(field)

	return nil
}
