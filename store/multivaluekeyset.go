package store

import (
	log "github.com/colt3k/nglog/ng"
)

type ValueSet struct {
	Vals []interface{}
}
type MVKeySetMap struct {
	Sets map[interface{}]ValueSet
}

func NewMVKeySet() *MVKeySetMap {
	return &MVKeySetMap{make(map[interface{}]ValueSet)}
}

func (set *MVKeySetMap) Add(k, v interface{}) bool {
	st, found := set.Sets[k]
	//Not found create new array and add our value
	if !found {
		ar := make([]interface{}, 0)
		ar = append(ar, v)
		set.Sets[k] = ValueSet{ar}
	} else {
		st.Vals = append(st.Vals, v)
		set.Sets[k] = st
	}

	return true //False if it existed already
}

func (set *MVKeySetMap) ContainsKey(k interface{}) bool {
	_, found := set.Sets[k]
	return found //true if it existed already
}
func (set *MVKeySetMap) ContainsVal(k, v interface{}) bool {
	m := set.Sets[k]
	for _, d := range m.Vals {
		if d == v {
			return true
		}
	}
	return false
}

func (set *MVKeySetMap) RemoveKey(k interface{}) {
	log.Logf(log.DEBUG, "Removing Key %v", k)
	delete(set.Sets, k)
}
func (set *MVKeySetMap) RemoveVal(k, v interface{}) {
	m := set.Sets[k]
	ar := make([]interface{}, 0)

	for _, d := range m.Vals {
		if d != v {
			ar = append(ar, d)
		}
	}
	if len(ar) > 0 {
		set.Sets[k] = ValueSet{ar}
	}
}

func (set *MVKeySetMap) SizeKeys() int {
	return len(set.Sets)
}
func (set *MVKeySetMap) SizeVals(k interface{}) int {
	return len(set.Sets[k].Vals)
}
