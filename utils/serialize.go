package utils

import (
	"fmt"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
  "errors"
)

// json编码
func JsonEncode(data interface{}) (string, error) {
	a, err := json.Marshal(data)
	return string(a), err
}

// json解码
func JsonDecode(data string) (interface{}, error) {
	dataByte := []byte(data)
	var dat interface{}

	err := json.Unmarshal(dataByte, &dat)
	return dat, err
}

func Encode(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// 设置Struct中某个属性值
func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)
	if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
			return errors.New("Provided value type didn't match obj field type")
	}
	structFieldValue.Set(val)
	return nil
}