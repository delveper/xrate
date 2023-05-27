package config

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func LoadVars() (err error) {
	env, err := os.Open(".env")
	if err != nil {
		return fmt.Errorf("error during opening environment file: %w", err)
	}

	defer func() {
		if err = env.Close(); err != nil {
			err = fmt.Errorf("error during closing environment file: %w", err)
		}
	}()

	buf := bufio.NewScanner(env)
	buf.Split(bufio.ScanLines)

	for buf.Scan() {
		if keyVal := strings.Split(buf.Text(), "="); len(keyVal) > 1 {
			if err := os.Setenv(keyVal[0], keyVal[1]); err != nil {
				return fmt.Errorf("error during setting environment variable: %w", err)
			}
		}
	}

	return nil
}

func ParseVars[T any]() (*T, error) {
	if err := LoadVars(); err != nil {
		return nil, err
	}

	dst := new(T)
	val := reflect.ValueOf(dst).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		if field.Type.Kind() == reflect.Struct {
			for j := 0; j < field.Type.NumField(); j++ {
				envTag := field.Type.Field(j).Tag.Get("env")
				if envTag != "" {
					envVal := os.Getenv(envTag)
					if envVal == "" {
						continue
					}
					switch field.Type.Field(j).Type.Kind() {
					case reflect.Int:
						integer, err := strconv.Atoi(envVal)
						if err != nil {
							return nil, err
						}
						val.Field(i).Field(j).SetInt(int64(integer))

					case reflect.Float64:
						float, err := strconv.ParseFloat(envVal, 64)
						if err != nil {
							return nil, err
						}
						val.Field(i).Field(j).SetFloat(float)

					case reflect.TypeOf(time.Duration(0)).Kind():
						duration, err := time.ParseDuration(envVal)
						if err != nil {
							return nil, err
						}
						val.Field(i).Field(j).Set(reflect.ValueOf(duration))

					default:
						val.Field(i).Field(j).SetString(envVal)
					}
				}
			}
		}
	}

	return dst, nil
}
