package stream

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
)

/*
Usage example:
	file := //open file object to write into
	reader := bytes.NewReader([]byte("some text to encrypt or file data"))
	se, err := NewStreamEncrypter(key, reader)
	io.Copy(file, se) OR encrypted, err := ioutil.ReadAll(se)


	sd, err := NewStreamDecrypter(key, se.Meta(), bytes.NewReader(encrypted))
	decrypted, err := ioutil.ReadAll(sd)
 */


// NewStreamEncrypter creates a new stream encrypter
func NewStreamEncrypter(key []byte, plainText io.Reader) (*StreamEncrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%+v", err)
	}
	iv := make([]byte, block.BlockSize())
	_, err = rand.Read(iv)
	if err != nil {
		return nil, fmt.Errorf("%+v", err)
	}
	stream := cipher.NewCTR(block, iv)
	mac := hmac.New(sha256.New, key)
	return &StreamEncrypter{
		Source: plainText,
		Block:  block,
		Stream: stream,
		Mac:    mac,
		IV:     iv,
	}, nil
}

// NewStreamDecrypter creates a new stream decrypter
func NewStreamDecrypter(key []byte, meta StreamMeta, cipherText io.Reader) (*StreamDecrypter, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("%+v", err)
	}
	stream := cipher.NewCTR(block, meta.IV)
	mac := hmac.New(sha256.New, key)
	return &StreamDecrypter{
		Source: cipherText,
		Block:  block,
		Stream: stream,
		Mac:    mac,
		Meta:   meta,
	}, nil
}

// StreamEncrypter is an encrypter for a stream of data with authentication
type StreamEncrypter struct {
	Source io.Reader
	Block  cipher.Block
	Stream cipher.Stream
	Mac    hash.Hash
	IV     []byte
}

// StreamDecrypter is a decrypter for a stream of data with authentication
type StreamDecrypter struct {
	Source io.Reader
	Block  cipher.Block
	Stream cipher.Stream
	Mac    hash.Hash
	Meta   StreamMeta
}

// Read encrypts the bytes of the inner reader and places them into p
func (s *StreamEncrypter) Read(p []byte) (int, error) {
	n, readErr := s.Source.Read(p)
	if n > 0 {
		s.Stream.XORKeyStream(p[:n], p[:n])
		err := writeHash(s.Mac, p[:n])
		if err != nil {
			return n, fmt.Errorf("%+v", err)
		}
		return n, readErr
	}
	return 0, io.EOF
}

// Meta returns the encrypted stream metadata for use in decrypting. This should only be called after the stream is finished
func (s *StreamEncrypter) Meta() StreamMeta {
	return StreamMeta{IV: s.IV, Hash: s.Mac.Sum(nil)}
}

// Read reads bytes from the underlying reader and then decrypts them
func (s *StreamDecrypter) Read(p []byte) (int, error) {
	n, readErr := s.Source.Read(p)
	if n > 0 {
		err := writeHash(s.Mac, p[:n])
		if err != nil {
			return n, fmt.Errorf("%+v", err)
		}
		s.Stream.XORKeyStream(p[:n], p[:n])
		return n, readErr
	}
	return 0, io.EOF
}

// Authenticate verifys that the hash of the stream is correct. This should only be called after processing is finished
func (s *StreamDecrypter) Authenticate() error {
	if !hmac.Equal(s.Meta.Hash, s.Mac.Sum(nil)) {
		return fmt.Errorf("authentication failed")
	}
	return nil
}

func writeHash(mac hash.Hash, p []byte) error {
	m, err := mac.Write(p)
	if err != nil {
		return fmt.Errorf("%+v", err)
	}
	if m != len(p) {
		return fmt.Errorf("could not write all bytes to hmac")
	}
	return nil
}

func checkedWrite(dst io.Writer, p []byte) (int, error) {
	n, err := dst.Write(p)
	if err != nil {
		return n, fmt.Errorf("%+v", err)
	}
	if n != len(p) {
		return n, fmt.Errorf("unable to write all bytes")
	}
	return len(p), nil
}

// StreamMeta is metadata about an encrypted stream
type StreamMeta struct {
	// IV is the initial value for the crypto function
	IV []byte
	// Hash is the sha256 hmac of the stream
	Hash []byte
}
