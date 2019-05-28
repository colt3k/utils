package osut

import "fmt"

func ExampleOS() {

	fmt.Println("Win", Windows())
	fmt.Println("Mac", Mac())
	fmt.Println("Linux", Linux())
	fmt.Println("Android", Android())

	/*
			 Output:
			 Win false
		Mac true
		Linux false
		Android false
	*/
}
