package env

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	tagEnv     = "env"
	tagDefault = "default"
)

func Parse[T any](pth string) (*T, error) {
	if pth == "" {
		pth = ".env"
	}

	if err := Load(pth); err != nil {
		return nil, err
	}

	dst := new(T)
	dstVal := reflect.ValueOf(dst).Elem()

	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Type().Field(i)

		if field.Type.Kind() != reflect.Struct {
			continue
		}

		for j := 0; j < field.Type.NumField(); j++ {
			dstTag := field.Type.Field(j).Tag.Get(tagEnv)
			srcVal := os.Getenv(dstTag)

			fieldType := field.Type.Field(j)
			fieldVal := dstVal.Field(i).Field(j)

			if srcVal == "" {
				srcVal = fieldType.Tag.Get(tagDefault)
			}

			if srcVal == "" {
				return nil, fmt.Errorf("no value for field: %s", field.Type.Field(j).Name)
			}

			if err := setFieldValue(fieldType.Type, fieldVal, srcVal); err != nil {
				return nil, err
			}
		}
	}

	return dst, nil
}

func Load(pth string) error {
	if pth == "" {
		pth = ".env"
	}

	env, err := os.Open(pth)
	if err != nil {
		return fmt.Errorf("opening environment file: %w", err)
	}

	defer func() {
		if err := env.Close(); err != nil {
			log.Printf("closing environment file: %v", err)
		}
	}()

	buf := bufio.NewScanner(env)
	buf.Split(bufio.ScanLines)

	const numElements = 2

	for buf.Scan() {
		pair := strings.Split(buf.Text(), "=")
		if len(pair) != numElements {
			continue
		}

		if err := os.Setenv(pair[0], pair[1]); err != nil {
			return fmt.Errorf("setting environment variable: %w", err)
		}
	}

	return nil
}

func setFieldValue(fieldType reflect.Type, fieldVal reflect.Value, val string) error {
	switch fieldType.Kind() {
	case reflect.Int:
		integer, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		fieldVal.SetInt(int64(integer))

	case reflect.Float64:
		float, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		fieldVal.SetFloat(float)

	case reflect.TypeOf(time.Duration(0)).Kind():
		duration, err := time.ParseDuration(val)
		if err != nil {
			return err
		}
		fieldVal.Set(reflect.ValueOf(duration))

	default:
		fieldVal.SetString(val)
	}

	return nil
}
