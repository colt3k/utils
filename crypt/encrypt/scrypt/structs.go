package scrypt

import (
	"bytes"
	"strconv"
)

//Params store values for scrypt encryption
type Params struct {
	N       int    `json:"iterations"`               // CPU/memory cost parameter (logN)
	R       int    `json:"block_size"`               // block size parameter (octets)
	P       int    `json:"parallelism"`              // parallelisation parameter (positive int)
	SaltLen int    `json:"salt_bytes"`               // bytes to use as salt (octets)
	DKLen   int    `json:"derived_key_length"`       // length of the derived key (octets)
	Salt    []byte `json:"generated_or_stored_salt"` // generated or stored salt
	Dk      []byte `json:"derived_key"`              // derived key
}

func (p *Params) String() string {
	var byt bytes.Buffer
	byt.WriteString("\n\tIterations: ")
	byt.WriteString(strconv.Itoa(p.N))
	byt.WriteString("\n\tBlockSize: ")
	byt.WriteString(strconv.Itoa(p.R))
	byt.WriteString("\n\tParallelism: ")
	byt.WriteString(strconv.Itoa(p.P))
	byt.WriteString("\n\tSaltLen: ")
	byt.WriteString(strconv.Itoa(p.SaltLen))
	byt.WriteString("\n\tDerivedKeyLength: ")
	byt.WriteString(strconv.Itoa(p.DKLen))
	byt.WriteString("\n\tSalt: ")
	byt.WriteString(string(p.Salt))
	byt.WriteString("\n\tDerivedKey: ")
	byt.WriteString(string(p.Dk))

	return byt.String()
}