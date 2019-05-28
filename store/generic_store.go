package store

//Put whatever you want in this to move it around
type FormatStore []interface{}

func (f *FormatStore) Put(elem interface{}) {
	*f = append(*f, elem)
}

//Get gets an element from the container.
func (f *FormatStore) Get() interface{} {
	elem := (*f)[0]
	*f = (*f)[1:]
	return elem
}

//Put in calling code to ensure correct type
//func assertExample() {
//
//	intContainer := &Container{}
//	intContainer.Put(7)
//	intContainer.Put(42)
//
//	elem, ok := intContainer.Get().(int) // assert that the actual type is int
//	if !ok {
//		fmt.Println("Unable to read an int from intContainer")
//	}
//	fmt.Printf("assertExample: %d (%T)\n", elem, elem)
//}
