package env

import (
	"fmt"
	"testing"

	log "github.com/colt3k/nglog/ng"
)

func setup() {
	tmp := New()
	tmp.Add("CT_SETUP", "One")
	tmp.Add("CT_setUP", "Two")
	tmp.Add("CT_Test1", "3")
	tmp.Add("cT_Test2", "4")
	tmp.Add("CT_TestCT", "5")
}

func TestEnvironment_All(t *testing.T) {
	tmp := New()
	myMap := tmp.All()
	var count int
	for k, v := range myMap {
		count++
		log.Println(k, ":", v)
	}
	if count <= 0 {
		t.Failed()
	}
}

func TestEnvironment_Prefix(t *testing.T) {
	setup()
	tmp := New()
	myMap := tmp.Prefix("ct", true)
	var count int
	for k, v := range myMap {
		count++
		fmt.Println("Find Lower Prefix Sensitive", k, ":", v)
	}
	log.Println("Count: shouldn't have found any ", count)
	if count > 0 {
		t.FailNow()
	}

	fmt.Println("")
	count = 0
	myMap = tmp.Prefix("CT", true)
	for k, v := range myMap {
		count++
		fmt.Println("Find CAPS Prefix Sensitive", k, ":", v)
	}
	log.Println("Count: should have found four ", count)
	if count < 0 {
		t.FailNow()
	}

	fmt.Println("")
	count = 0
	myMap = tmp.Prefix("ct", false)
	for k, v := range myMap {
		count++
		fmt.Println("Find Lower Prefix Non Sensitive", k, ":", v)
	}
	log.Println("Count: should have found five: ", count)
	if count != 5 {
		t.Error("should have found at least one")
		t.FailNow()
	}
}

func TestEnvironment_Suffix(t *testing.T) {
	setup()
	tmp := New()
	myMap := tmp.Suffix("ct", true)
	var count int
	for k, v := range myMap {
		count++
		fmt.Println("Find Lower Suffix Sensitive", k, ":", v)
	}
	log.Println("Count: shouldn't have found any ", count)
	if count > 0 {
		t.FailNow()
	}

	fmt.Println("")
	count = 0
	myMap = tmp.Suffix("CT", true)
	for k, v := range myMap {
		count++
		fmt.Println("Find CAPS Suffix Sensitive", k, ":", v)
	}
	log.Println("Count: should have found one ", count)
	if count < 0 {
		t.FailNow()
	}

	fmt.Println("")
	count = 0
	myMap = tmp.Suffix("ct", false)
	for k, v := range myMap {
		count++
		fmt.Println("Find Lower Suffix Non Sensitive", k, ":", v)
	}
	log.Println("Count: should have found one: ", count)
	if count != 1 {
		t.Error("should have found at least one")
		t.FailNow()
	}
}

func TestEnvironment_Includes(t *testing.T) {
	setup()
	tmp := New()
	myMap := tmp.Includes("ct", true)
	var count int
	for k, v := range myMap {
		count++
		fmt.Println("Find Includes Sensitive", k, ":", v)
	}
	log.Println("Count: shouldn't have found any ", count)
	if count > 0 {
		t.FailNow()
	}

	fmt.Println("")
	count = 0
	myMap = tmp.Includes("CT", true)
	for k, v := range myMap {
		count++
		fmt.Println("Find CAPS Includes Sensitive", k, ":", v)
	}
	log.Println("Count: should have found some ", count)
	if count <= 0 {
		t.FailNow()
	}

	fmt.Println("")
	count = 0
	myMap = tmp.Includes("ct", false)
	for k, v := range myMap {
		count++
		fmt.Println("Find Lower Includes Non Sensitive", k, ":", v)
	}
	log.Println("Count: should have found some: ", count)
	if count <= 0 {
		t.Error("should have found at least one")
		t.FailNow()
	}
}
