package hc

import "testing"

func TestProcessorMake(t *testing.T) {
	mapData, status, err := MakeCall("GET", "http://www.google.com", nil, nil, nil, NewClientSettings(true, true, 120, 120))
	if err != nil {
		t.Logf("error %v", err)
		t.Fail()
	}
	t.Logf("success received status %v, %v", status, mapData)
}

func TestProcessorDirect(t *testing.T) {
	c := NewHTTPClient("GET", "http://www.google.com", nil, nil, NewClientSettings(true, true, 120, 120))
	mapData, status, err := c.Process(nil)
	if err != nil {
		t.Logf("error %v", err)
		t.Fail()
	}
	t.Logf("success received status %v, %v", status, mapData)
}
