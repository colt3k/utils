// Package qrand pulls quantum random bytes from the ANU QRNG API and converts
// them into seeds.
package qrand

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
)

// URL is the public QRNG endpoint used by this package.
const URL = "https://qrng.anu.edu.au/API/jsonI.php"

// Response matches the QRNG JSON payload.
type Response struct {
	Type    string  `json:"type"`
	Length  int     `json:"length"`
	Data    []uint8 `json:"data"`
	Success bool    `json:"success"`
}

// GenerateSeedData requests amount bytes and returns the first 8 bytes as a
// big-endian uint64 seed.
func GenerateSeedData(amount int) (uint64, error) {

	ba := make([]byte, amount)

	response, err := pullData(amount)
	if err != nil {
		return 0, err
	}

	for i := 0; i < amount; i++ {
		b := response.Data[i]
		ba[i] = b
	}

	seed := binary.BigEndian.Uint64(ba)

	return seed, nil
}

func pullData(amount int) (Response, error) {
	var response Response

	urlWithParms := fmt.Sprintf("%s?length=%v&type=uint8", URL, amount)
	req, err := http.NewRequest(http.MethodGet, urlWithParms, nil)
	if err != nil {
		return response, err
	}

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return response, err
	}

	json.NewDecoder(resp.Body).Decode(&response)

	return response, nil
}
