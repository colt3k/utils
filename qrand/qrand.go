package qrand

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
)

const URL = "https://qrng.anu.edu.au/API/jsonI.php"

type Response struct {
	Type    string  `json:"type"`
	Length  int     `json:"length"`
	Data    []uint8 `json:"data"`
	Success bool    `json:"success"`
}

func GenerateSeedData(amount int) (uint64, error) {

	ba := make([]byte, amount)
	urlWithParms := fmt.Sprintf("%s?length=%v&type=uint8", URL, amount)
	req, err := http.NewRequest(http.MethodGet, urlWithParms, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return 0, err
	}

	var response Response
	json.NewDecoder(resp.Body).Decode(&response)

	for i:=0; i < amount; i++ {
		b := response.Data[i]
		ba[i]=b
	}

	seed := binary.BigEndian.Uint64(ba)

	return seed,nil
}
