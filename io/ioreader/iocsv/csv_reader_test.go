package iocsv

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

var content string

func TestMain(m *testing.M) {
	content = `
***********************************************************************************
username=jc|firstname=Joe|lastname=Coe|suffix=Jr|prefix=Mr.|dob=1/1/2008|

***********************************************************************************
username=jdoe|firstname=John|lastname=Doe|suffix=|prefix=Mr.|dob=1/1/1900|

***********************************************************************************
username=jadoe|firstname=Jane|lastname=Doe|suffix=|prefix=Mrs.|dob=1/1/1800|
`
	m.Run()
}
func TestReadCSV(t *testing.T) {

	data, err := ReadCSVFromFile("./test.csv", true, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("Rows", len(data.Rows))
	var buff bytes.Buffer
	for _, row := range data.Rows {
		buff.WriteString("\n***********************************************************************************\n")
		for _, col := range row.Cols {
			buff.WriteString(col.Name + "=" + col.Data + "|")
		}
		buff.WriteString("\n")
	}
	fmt.Println(buff.String())
	if buff.String() != content {
		t.Fail()
	}
}
