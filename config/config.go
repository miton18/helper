package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func LoadConfig[T any]() (*T, error) {
	cfg := new(T)
	st := reflect.TypeOf(*cfg)
	vt := reflect.ValueOf(cfg).Elem()

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		fieldV := vt.FieldByIndex(field.Index)

		if !fieldV.IsValid() {
			return nil, fmt.Errorf("field is invalid: %+v", fieldV)
		}
		if !fieldV.CanSet() {
			return nil, fmt.Errorf("filed cannot be set: %+v", fieldV)
		}

		// apply default
		if defaultValue, ok := field.Tag.Lookup("default"); ok {
			switch field.Type.Kind() {
			case reflect.String:
				fieldV.SetString(defaultValue)
			case reflect.Int:
				i, err := strconv.ParseInt(defaultValue, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid default: %w", err)
				}
				fieldV.SetInt(i)
			}
		}

		// load from env
		envKey, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}

		envValue, ok := os.LookupEnv(envKey)
		if !ok {
			continue
		}

		if un, ok := fieldV.Addr().Interface().(EnvUnmarshaller); ok {
			//logrus.Debugf("pointer unmarshaller")
			if err := un.UnmarshalText(envValue); err != nil {
				return nil, fmt.Errorf("failed to parse value '%+v': %w", envValue, err)
			}
		} else if ok := field.Type.Implements(EnvUnmarshallerKind); ok {
			//logrus.Debugf("canInterface:%+v", fieldV.CanInterface())

			unmarshaller, ok := fieldV.Interface().(EnvUnmarshaller)
			if !ok {
				//logrus.Debugf("NOT OK") // must not append
			} else {
				//logrus.Debugf("value unmarshaller")
				if err := unmarshaller.UnmarshalText(envValue); err != nil {
					return nil, fmt.Errorf("failed to parse value '%+v': %w", envValue, err)
				}
			}

		} else {

			switch field.Type.Kind() {
			case reflect.String:
				fieldV.SetString(envValue)
			case reflect.Int:
				i, err := strconv.ParseInt(envValue, 10, 0)
				if err != nil {
					return nil, fmt.Errorf("cannot '%s' as int: %w", envValue, err)
				}
				fieldV.SetInt(int64(i))
			// case:TextUnmarshaler // idee

			default:
				return nil, fmt.Errorf("unsupported type: %+v", field)
			}
		}
	}

	return cfg, nil
}

type EnvUnmarshaller interface {
	UnmarshalText(string) error
}

var EnvUnmarshallerKind = reflect.TypeFor[EnvUnmarshaller]()
