package env

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestParseTo(t *testing.T) {
	type config struct {
		Home    string `env:"HOME"`
		Default string `default:"default"`
		Empty   int
		Nested  struct{ Value string }
	}

	cases := []struct {
		name    string
		vars    map[string]string
		want    *config
		wantErr error
	}{
		{
			name:    "All environment variables present",
			vars:    map[string]string{"HOME": "/home/test", "DEFAULT": "default", "EMPTY": "10", "NESTED_VALUE": "nested"},
			want:    &config{Home: "/home/test", Default: "default", Empty: 10, Nested: struct{ Value string }{"nested"}},
			wantErr: nil,
		},
		{
			name:    "No environment variables, except default present",
			vars:    nil,
			want:    new(config),
			wantErr: fmt.Errorf("no environment variables")},
		{
			name:    "Some environment variables present",
			vars:    map[string]string{"HOME": "/home/test", "EMPTY": "10"},
			want:    &config{Home: "/home/test", Empty: 10, Default: "default"},
			wantErr: errors.New("no value for field: Value"),
		},
		{
			name:    "Environment variable set to empty string",
			vars:    map[string]string{"HOME": "", "DEFAULT": "", "EMPTY": "0", "NESTED_VALUE": ""},
			want:    &config{Home: "", Default: "", Empty: 0},
			wantErr: errors.New("no environment variables"),
		},
		{
			name:    "Environment variable set to invalid value",
			vars:    map[string]string{"HOME": "/home/test", "DEFAULT": "default", "EMPTY": "invalid"},
			want:    nil,
			wantErr: fmt.Errorf("parsing integer: strconv.ParseInt: parsing \"invalid\": invalid syntax"),
		},
		{
			name:    "Nested struct environment variable set",
			vars:    map[string]string{"HOME": "/home/test", "DEFAULT": "default", "EMPTY": "10", "NESTED_VALUE": "nested"},
			want:    &config{Home: "/home/test", Default: "default", Empty: 10, Nested: struct{ Value string }{"nested"}},
			wantErr: nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			teardown, err := setupEnv(t, tt.vars)
			require.NoError(t, err)
			defer teardown()

			var cfg config
			err = parseTo(&cfg, "")

			if tt.wantErr != nil {
				assert.Error(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, &cfg)
		})
	}
}

func setupEnv(t *testing.T, vars map[string]string) (func(), error) {
	if vars == nil {
		return func() {}, nil
	}

	for k, v := range vars {
		if err := os.Setenv(k, v); err != nil {
			return func() {}, fmt.Errorf("setupEnv: setting %s[%s]: %v", k, v, err)
		}
	}

	teardown := func() {
		for k := range vars {
			if err := os.Unsetenv(k); err != nil {
				t.Errorf("setEnv: unsetting %s: %v", k, err)
			}
		}
	}

	return teardown, nil
}

func TestLoadEnv(t *testing.T) {
	cases := []struct {
		name    string
		vars    map[string]string
		wantErr error
	}{
		{
			name:    "Valid environment file",
			vars:    map[string]string{"HOME": "/home/test", "DEFAULT": "default", "EMPTY": "10", "NESTED_VALUE": "nested"},
			wantErr: nil,
		},
		{
			name:    "Empty variables",
			vars:    map[string]string{"": ""},
			wantErr: errors.New("setting []: setenv: The parameter is incorrect"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			file, teardown, err := createEnvFile(t, tt.vars)
			require.NoError(t, err)
			defer teardown()

			teardown, err = setupEnv(t, tt.vars)
			if tt.wantErr != nil {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			defer teardown()

			err = loadEnv(file.Name())
			if tt.wantErr != nil {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			for k, v := range tt.vars {
				assert.Equal(t, v, os.Getenv(k))
			}
		})
	}
}

func createEnvFile(t *testing.T, vars map[string]string) (*os.File, func(), error) {
	file, err := os.CreateTemp(os.TempDir(), "test.env")
	if err != nil {
		return nil, func() {}, err
	}

	if vars == nil {
		return file, func() {}, nil
	}

	for key, value := range vars {
		_, err := file.WriteString(fmt.Sprintf("%s=%s\n", key, value))
		if err != nil {
			return nil, nil, err
		}
	}

	teardown := func() {
		err := file.Close()
		if err != nil {
			t.Error("createEnvFile: close:", err)
		}

		err = os.RemoveAll(file.Name())
		if err != nil {
			t.Error("createEnvFile: teardown:", err)
		}
	}

	return file, teardown, nil
}

func TestSetFieldValue(t *testing.T) {
	cases := []struct {
		name    string
		field   interface{}
		have    string
		want    interface{}
		wantErr error
	}{
		{"Duration", time.Duration(0), "1h", time.Hour, nil},
		{"Invalid Duration", time.Duration(0), "invalid", nil, errors.New("parsing duration: time: invalid duration \"invalid\"")},
		{"Int", int(0), "123", 123, nil},
		{"Invalid Int", int(0), "invalid", nil, errors.New("parsing integer: strconv.ParseInt: parsing \"invalid\": invalid syntax")},
		{"Float", float64(0), "1.23", float64(1.23), nil},
		{"Invalid Float", float64(0), "invalid", nil, errors.New("parsing float: strconv.ParseFloat: parsing \"invalid\": invalid syntax")},
		{"Bool", bool(false), "true", true, nil},
		{"Invalid Bool", bool(false), "invalid", nil, errors.New("parsing bool: strconv.ParseBool: parsing \"invalid\": invalid syntax")},
		{"String", string(""), "test", "test", nil},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			fieldVal := reflect.New(reflect.TypeOf(tt.field)).Elem()
			err := setFieldValue(reflect.TypeOf(tt.field), fieldVal, tt.have)

			if tt.wantErr != nil {
				assert.Error(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, fieldVal.Interface())
		})
	}
}

func TestCamelToSnake(t *testing.T) {
	cases := []struct {
		name string
		have string
		want string
	}{
		{name: "Pascal case", have: "PascalCase", want: "PASCAL_CASE"},
		{name: "Camel case", have: "camelCase", want: "CAMEL_CASE"},
		{name: "Snake case", have: "snake_case", want: "SNAKE_CASE"},
		{name: "Lowercase", have: "lowercase", want: "LOWERCASE"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			result := camelToSnake(tt.have)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestParseLine(t *testing.T) {
	cases := []struct {
		name      string
		have      string
		wantKey   string
		wantValue string
	}{
		{name: "Line with key-value pair", have: "KEY=value", wantKey: "KEY", wantValue: "value"},
		{name: "Line with empty key", have: "EMPTY_KEY=", wantKey: "EMPTY_KEY", wantValue: ""},
		{name: "Line with empty value", have: "=EMPTY_VALUE", wantKey: "", wantValue: ""},
		{name: "Line without equal sign", have: "NO_EQUAL_SIGN", wantKey: "", wantValue: ""},
		{name: "Another have with key-value pair", have: "ANOTHER_CASE=another_value", wantKey: "ANOTHER_CASE", wantValue: "another_value"},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			key, value := parseLine(tt.have)
			assert.Equal(t, tt.wantKey, key)
			assert.Equal(t, tt.wantValue, value)
		})
	}
}
