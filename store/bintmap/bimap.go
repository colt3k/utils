package bintmap

type BiMap interface {
	PutIntKey(int, string)
	GetByKey(key int) (val string, exists bool)
	GetByValue(val string) (key int, exists bool)
	DeleteKey(key int) BiMap
	DeleteValue(val string) BiMap
}
type BiMapInt struct {
	a map[int]string
	b map[string]int
}

func NewIntBiMap() BiMap {
	t := new(BiMapInt)
	t.a = make(map[int]string)
	t.b = make(map[string]int)
	return t
}
func NewIntBiMapInitd(init map[int]string) BiMap {
	t := new(BiMapInt)
	t.a = init
	t.b = make(map[string]int)
	for k,v := range t.a {
		t.b[v]=k
	}
	return t
}

func (b *BiMapInt) PutIntKey(key int, val string) {
	b.a[key] = val
	b.b[val] = key
}
func (b *BiMapInt) GetByKey(key int) (val string, exists bool) {
	val, exists = b.a[key]
	return
}
func (b *BiMapInt) GetByValue(val string) (key int, exists bool) {
	key, exists = b.b[val]
	return
}
func (b *BiMapInt) Len() int {
	return len(b.a)
}

func (b *BiMapInt) DeleteKey(key int) BiMap {
	value, exists := b.a[key]
	if exists {
		delete(b.a, key)
		delete(b.b, value)
	}
	return b
}

func (b *BiMapInt) DeleteValue(val string) BiMap {
	key, exists := b.b[val]
	if exists {
		delete(b.a, key)
		delete(b.b, val)
	}
	return b
}