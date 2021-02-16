package example1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestInterfaceToReflection(t *testing.T) {

	var x = 3.4
	xValue := reflect.ValueOf(x)
	xType := reflect.TypeOf(x)
	xTypeFromValue := xValue.Type()
	assert.Equal(t, xType, xTypeFromValue)
	assert.Equal(t, xType.String(), "float64")
	assert.Equal(t, xType.Kind(), reflect.Float64)
	assert.Equal(t, xValue.Kind(), reflect.Float64)
	assert.Equal(t, xValue.Float(), 3.4)

	var x2 uint8 = 'x'
	assert.Equal(t, uint8(reflect.ValueOf(x2).Uint()), uint8(120))

	type customInt int
	var x3 customInt = 7
	assert.Equal(t, reflect.ValueOf(x3).Kind(), reflect.Int)

}

func TestReflectionToInterface(t *testing.T) {
	var y = 5.99
	yValue := reflect.ValueOf(y)
	assert.Equal(t, yValue.Interface(), 5.99)
	assert.Equal(t, yValue.Interface().(float64), 5.99)
}

func TestModifyReflectionObject(t *testing.T) {
	var z = 19.364
	zValue := reflect.ValueOf(z)
	assert.False(t, zValue.CanSet())

	zValuePointer := reflect.ValueOf(&z)
	assert.True(t, zValuePointer.Elem().CanSet())
	zValuePointer.Elem().SetFloat(15.55)

	assert.Equal(t, z, 15.55)
}

func TestReflectionStruct(t *testing.T) {
	type T struct {
		A int
		B string
	}

	obj := T{23, "ski-doo"}
	objValue := reflect.ValueOf(&obj).Elem()

	typeOfObj := objValue.Type()
	for i := 0; i < objValue.NumField(); i++ {
		f := objValue.Field(i)
		fieldName := typeOfObj.Field(i).Name
		switch fieldName {
		case "A":
			assert.Equal(t, fmt.Sprintf("%s %s = %v", typeOfObj.Field(i).Name, f.Type(), f.Interface()), "A int = 23")
		case "B":
			assert.Equal(t, fmt.Sprintf("%s %s = %v", typeOfObj.Field(i).Name, f.Type(), f.Interface()), "B string = ski-doo")
		default:
			t.Error()
		}
	}

	objValue.Field(0).SetInt(11)
	objValue.Field(1).SetString("Sunset Strip")
	assert.Equal(t, obj.A, 11)
	assert.Equal(t, obj.B, "Sunset Strip")

}
