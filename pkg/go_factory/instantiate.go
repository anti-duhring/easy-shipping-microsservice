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
		currentField := structType.Field(i)
		fieldName := currentField.Name

		// Get the field by name
		field := structValue.FieldByName(fieldName)

		// Check if the field exists
		if !field.IsValid() {
			return fmt.Errorf("field %s does not exist in the struct", fieldName)
		}

		// Check if the field can be set (exported)
		if !field.CanSet() {
			return fmt.Errorf("field %s cannot be set", fieldName)
		}

		if defaultValues != nil {
			defaultFieldValue := reflect.ValueOf(defaultValues).Elem().FieldByName(fieldName)

			if defaultFieldValue.IsValid() {
				field.Set(defaultFieldValue)

				break
			}

		}

		setNewValue, ok := typesMap[currentField.Type.Name()]

		if !ok {
			return fmt.Errorf("type %s not supported", currentField.Type.Name())
		}

		setNewValue(&field)
	}

	return nil
}
