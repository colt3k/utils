package iocsv

import (
	"bytes"

	ers "github.com/colt3k/nglog/ers/bserr"
	log "github.com/colt3k/nglog/ng"
)

func ExampleReadCSV() {
	data, err := ReadCSV("./test.csv", true, nil)
	ers.StopErr(err)
	log.Println("Rows", len(data.Rows))
	var buff bytes.Buffer
	for _, row := range data.Rows {
		buff.WriteString("\n***********************************************************************************\n")
		for _, col := range row.Cols {
			buff.WriteString(col.Name + "=" + col.Data + "|")
		}
		buff.WriteString("\n")
	}
	log.Println(buff.String())

	/*
			Output:
			***********************************************************************************
		username=jc|firstname=Joe|lastname=Coe|suffix=Jr|prefix=Mr.|dob=1/1/2008|

		***********************************************************************************
		username=jdoe|firstname=John|lastname=Doe|suffix=|prefix=Mr.|dob=1/1/1900|

		***********************************************************************************
		username=jadoe|firstname=Jane|lastname=Doe|suffix=|prefix=Mrs.|dob=1/1/1800|

	*/
}
