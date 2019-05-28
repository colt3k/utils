package store

import (
	"bytes"
	"strconv"
)

//StringSet store map of string[bool]
type StringSet struct {
	stringSet map[string]bool
}

//NewStringSet create a new StringSet
func NewStringSet() *StringSet {
	return &StringSet{make(map[string]bool)}
}

//Add add a new string to our set
func (set *StringSet) Add(i string) bool {
	_, found := set.stringSet[i]
	set.stringSet[i] = true
	return !found //False if it existed already
}

//Contains check if string exists in set
func (set *StringSet) Contains(i string) bool {
	_, found := set.stringSet[i]
	return found //true if it existed already
}

//Remove remove string from set
func (set *StringSet) Remove(i string) {
	delete(set.stringSet, i)
}

//Size determine size of set
func (set *StringSet) Size() int {
	return len(set.stringSet)
}

func (set *StringSet) ToString() string {
	var byt bytes.Buffer
	for i, d := range set.stringSet {
		byt.WriteString("\"")
		byt.WriteString(i)
		byt.WriteString(":")
		byt.WriteString(strconv.FormatBool(d))
		byt.WriteString("\",")
	}
	return byt.String()
}
