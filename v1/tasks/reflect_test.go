package tasks_test

import (
	"testing"

	"github.com/anzersy/machinery/v1/tasks"
)

func TestReflectValue(t *testing.T) {
	value, err := tasks.ReflectValue("bool", interface{}(false))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "bool" {
		t.Errorf("type is %v, want bool", value.Type().String())
	}

	value, err = tasks.ReflectValue("int", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "int" {
		t.Errorf("type is %v, want int", value.Type().String())
	}

	value, err = tasks.ReflectValue("int8", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "int8" {
		t.Errorf("type is %v, want int8", value.Type().String())
	}

	value, err = tasks.ReflectValue("int16", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "int16" {
		t.Errorf("type is %v, want int16", value.Type().String())
	}

	value, err = tasks.ReflectValue("int32", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "int32" {
		t.Errorf("type is %v, want int32", value.Type().String())
	}

	value, err = tasks.ReflectValue("int64", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "int64" {
		t.Errorf("type is %v, want int64", value.Type().String())
	}

	value, err = tasks.ReflectValue("uint", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "uint" {
		t.Errorf("type is %v, want uint", value.Type().String())
	}

	value, err = tasks.ReflectValue("uint8", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "uint8" {
		t.Errorf("type is %v, want uint8", value.Type().String())
	}

	value, err = tasks.ReflectValue("uint16", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "uint16" {
		t.Errorf("type is %v, want uint16", value.Type().String())
	}

	value, err = tasks.ReflectValue("uint32", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "uint32" {
		t.Errorf("type is %v, want uint32", value.Type().String())
	}

	value, err = tasks.ReflectValue("uint64", interface{}(float64(1)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "uint64" {
		t.Errorf("type is %v, want uint64", value.Type().String())
	}

	value, err = tasks.ReflectValue("float32", interface{}(float64(0.5)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "float32" {
		t.Errorf("type is %v, want float32", value.Type().String())
	}

	value, err = tasks.ReflectValue("float64", interface{}(float64(0.5)))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "float64" {
		t.Errorf("type is %v, want float64", value.Type().String())
	}

	value, err = tasks.ReflectValue("string", interface{}("123"))
	if err != nil {
		t.Error(err)
	}
	if value.Type().String() != "string" {
		t.Errorf("type is %v, want string", value.Type().String())
	}
}
