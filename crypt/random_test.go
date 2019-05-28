package crypt

import log "github.com/colt3k/nglog/ng"

func ExampleGenSalt() {

	var dataset = []int{10, 20, 30, 5, 20}

	var data []byte
	for i, d := range dataset {
		if data != nil {
			data = GenSalt(data, d)
		} else {
			data = GenSalt(nil, d)
		}
		show(data, i, d)
		if i > 0 {
			data = nil
		}
	}

	/*
		Output:

	*/
}

func show(data []byte, idx, size int) {
	log.Println(idx, "DataSize", len(string(data)), "should be:", size)
}
