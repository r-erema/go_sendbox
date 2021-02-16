package interface2struct

import (
	"reflect"
)

type WrongValueError struct{}

func (e *WrongValueError) Error() string {
	return "Wrong value"
}

func i2s(data interface{}, out interface{}) error {

	outValue := reflect.ValueOf(out)
	if outValue.Kind() != reflect.Ptr {
		return &WrongValueError{}
	}

	outValuePtr := outValue.Elem()

	switch outValuePtr.Type().Kind() {
	case reflect.Struct:
		if source, ok := data.(map[string]interface{}); ok {
			return fillStructure(outValuePtr, source)
		} else {
			return &WrongValueError{}
		}
	case reflect.Slice:
		if source, ok := data.([]interface{}); ok {
			return fillSlice(outValuePtr, source)
		} else {
			return &WrongValueError{}
		}
	}
	return nil
}

func fillSlice(outValue reflect.Value, source []interface{}) error {
	slice := reflect.MakeSlice(outValue.Type(), len(source), len(source))
	for i, value := range source {
		elemValue := slice.Index(i)
		_ = fillStructure(elemValue, value.(map[string]interface{}))
	}
	outValue.Set(slice)
	return nil
}

func fillStructure(outValue reflect.Value, source map[string]interface{}) error {
	for name, value := range source {
		for i := 0; i < outValue.NumField(); i++ {
			fieldType := outValue.Type().Field(i)
			fieldValue := outValue.Field(i)
			if fieldType.Name == name {
				switch fieldType.Type.Kind() {
				case reflect.Int:
					if value, ok := value.(float64); ok {
						fieldValue.SetInt(int64(value))
					} else {
						return &WrongValueError{}
					}

				case reflect.String:
					if value, ok := value.(string); ok {
						fieldValue.SetString(value)
					} else {
						return &WrongValueError{}
					}
				case reflect.Bool:
					if value, ok := value.(bool); ok {
						fieldValue.SetBool(value)
					} else {
						return &WrongValueError{}
					}
				case reflect.Struct:
					if value, ok := value.(map[string]interface{}); ok {
						_ = fillStructure(fieldValue, value)
					} else {
						return &WrongValueError{}
					}
				case reflect.Slice:
					if value, ok := value.([]interface{}); ok {
						_ = fillSlice(fieldValue, value)
					} else {
						return &WrongValueError{}
					}
				}
			}
		}
	}
	return nil
}
