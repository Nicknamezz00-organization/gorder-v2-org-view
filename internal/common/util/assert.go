package util

import (
	"reflect"

	"github.com/pkg/errors"
)

func AssertNotEmpty(fields ...any) error {
	for _, field := range fields {
		if isEmpty(field) {
			return errors.New("assert: not empty failed")
		}
	}
	return nil
}

func isEmpty(object interface{}) bool {

	// get nil case out of the way
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return isEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}
