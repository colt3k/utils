package env

import (
	"os"
	"strings"
)

var (
	std = New()
)

type Environment struct {
}

func New() *Environment {
	return new(Environment)
}

func Find(key string) string {
	return std.Find(key)
}
func Prefix(key string, caseSensitive bool) map[string]string {
	return std.Prefix(key, caseSensitive)
}
func Suffix(key string, caseSensitive bool) map[string]string {
	return std.Suffix(key, caseSensitive)
}
func Includes(key string, caseSensitive bool) map[string]string {
	return std.Includes(key, caseSensitive)
}
func All() map[string]string {
	return std.All()
}
func Add(key, val string) error {
	return std.Add(key, val)
}

func (e *Environment) Find(key string) string {
	tmp := os.Getenv(key)
	return strings.TrimSpace(tmp)
}

func (e *Environment) Prefix(key string, caseSensitive bool) map[string]string {

	if !caseSensitive {
		key = strings.ToLower(key)
	}
	data := make(map[string]string)
	for _, e := range os.Environ() {

		pair := strings.Split(e, "=")
		if strings.HasPrefix(caseChange(pair[0], caseSensitive), key) {
			data[pair[0]] = pair[1]
		}
	}
	return data
}
func (e *Environment) Suffix(key string, caseSensitive bool) map[string]string {
	if !caseSensitive {
		key = strings.ToLower(key)
	}
	data := make(map[string]string)
	for _, e := range os.Environ() {

		pair := strings.Split(e, "=")
		if strings.HasSuffix(caseChange(pair[0], caseSensitive), key) {
			data[pair[0]] = pair[1]
		}
	}
	return data
}
func caseChange(key string, caseSensitive bool) string {
	if !caseSensitive {
		return strings.ToLower(key)
	}
	return key
}
func (e *Environment) Includes(key string, caseSensitive bool) map[string]string {
	if !caseSensitive {
		key = strings.ToLower(key)
	}
	data := make(map[string]string)
	for _, e := range os.Environ() {

		pair := strings.Split(e, "=")
		if strings.Index(caseChange(pair[0], caseSensitive), key) > -1 {
			data[pair[0]] = pair[1]
		}
	}
	return data
}

func (e *Environment) All() map[string]string {

	data := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		data[pair[0]] = pair[1]
	}

	return data
}

func (e *Environment) Add(key, val string) error {
	return os.Setenv(key, val)
}
