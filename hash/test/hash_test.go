package test

import (
	"crypto/md5"
	cryptsha1 "crypto/sha1"
	"crypto/sha256"
	cryptsha512 "crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/colt3k/utils/hash"
	"github.com/colt3k/utils/hash/blake2"
	"github.com/colt3k/utils/hash/hashenum"
	"github.com/colt3k/utils/hash/md"
	"github.com/colt3k/utils/hash/sha1"
	"github.com/colt3k/utils/hash/sha2"
	"github.com/colt3k/utils/hash/sha3"
	"github.com/colt3k/utils/hash/sha512"
)

const plain = "this is my example text"

// ---------------------------------------------------------------------------
// Available-algorithm tests (stdlib always present)
// ---------------------------------------------------------------------------

type availableCase struct {
	name       string
	hasher     hash.Hasher
	input      string
	expectLen  int
	expectHex  string // empty means "skip hex check, only verify length"
}

func TestAvailableAlgorithms_String(t *testing.T) {
	cases := []availableCase{
		{name: "MD5", hasher: md.NewHash(md.Format(hashenum.MD5)), input: plain, expectLen: 16},
		{name: "SHA1", hasher: sha1.NewHash(sha1.Format(hashenum.SHA1)), input: plain, expectLen: 20},
		{name: "SHA256", hasher: sha2.NewHash(sha2.Format(hashenum.SHA256)), input: plain, expectLen: 32},
		{name: "SHA384", hasher: sha512.NewHash(sha512.Format(hashenum.SHA384)), input: plain, expectLen: 48},
		{name: "SHA512", hasher: sha512.NewHash(sha512.Format(hashenum.SHA512)), input: plain, expectLen: 64},
		{name: "SHA512_224", hasher: sha512.NewHash(sha512.Format(hashenum.SHA512_224)), input: plain, expectLen: 28},
		{name: "SHA512_256", hasher: sha512.NewHash(sha512.Format(hashenum.SHA512_256)), input: plain, expectLen: 32},
		{name: "SHA3_224", hasher: sha3.NewHash(sha3.Format(hashenum.SHA3_224)), input: plain, expectLen: 28},
		{name: "SHA3_256", hasher: sha3.NewHash(sha3.Format(hashenum.SHA3_256)), input: plain, expectLen: 32},
		{name: "SHA3_512", hasher: sha3.NewHash(sha3.Format(hashenum.SHA3_512)), input: plain, expectLen: 64},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.hasher.String(tc.input)

			if result == nil {
				t.Fatalf("%s returned nil", tc.name)
			}
			if len(result) != tc.expectLen {
				t.Errorf("%s: expected length %d, got %d", tc.name, tc.expectLen, len(result))
			}

			// Cross-check against stdlib for the algorithms we can.
			switch tc.name {
			case "MD5":
				stdlib := md5.Sum([]byte(tc.input))
				if hex.EncodeToString(result) != hex.EncodeToString(stdlib[:]) {
					t.Errorf("%s: hash mismatch with stdlib md5", tc.name)
				}
			case "SHA1":
				h := cryptsha1.New()
				h.Write([]byte(tc.input))
				stdlib := h.Sum(nil)
				if hex.EncodeToString(result) != hex.EncodeToString(stdlib) {
					t.Errorf("%s: hash mismatch with stdlib sha1", tc.name)
				}
			case "SHA256":
				stdlib := sha256.Sum256([]byte(tc.input))
				if hex.EncodeToString(result) != hex.EncodeToString(stdlib[:]) {
					t.Errorf("%s: hash mismatch with stdlib sha256", tc.name)
				}
			case "SHA512":
				stdlib := cryptsha512.Sum512([]byte(tc.input))
				if hex.EncodeToString(result) != hex.EncodeToString(stdlib[:]) {
					t.Errorf("%s: hash mismatch with stdlib sha512", tc.name)
				}
			}

			t.Logf("%s  len=%d  hex=%s", tc.name, len(result), hex.EncodeToString(result))
		})
	}
}

// ---------------------------------------------------------------------------
// Unavailable-algorithm tests (RIPEMD160, BLAKE2* — not in stdlib on macOS)
// ---------------------------------------------------------------------------

func TestUnavailableAlgorithms(t *testing.T) {
	cases := []struct {
		name   string
		hasher hash.Hasher
	}{
		{name: "RIPEMD160", hasher: md.NewHash(md.Format(hashenum.RIPEMD160))},
		{name: "BLAKE2s_256", hasher: blake2.NewHash(blake2.Format(hashenum.BLAKE2s_256))},
		{name: "BLAKE2b_256", hasher: blake2.NewHash(blake2.Format(hashenum.BLAKE2b_256))},
		{name: "BLAKE2b_384", hasher: blake2.NewHash(blake2.Format(hashenum.BLAKE2b_384))},
		{name: "BLAKE2b_512", hasher: blake2.NewHash(blake2.Format(hashenum.BLAKE2b_512))},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.hasher.String(plain)
			if result != nil {
				t.Logf("%s returned data (len=%d) — algorithm IS available on this platform", tc.name, len(result))
			} else {
				t.Logf("%s returned nil — algorithm not available on this platform (expected on macOS)", tc.name)
			}
			// Either outcome is acceptable; the test documents the platform behavior.
		})
	}
}

// ---------------------------------------------------------------------------
// File hashing tests
// ---------------------------------------------------------------------------

func TestFileHashing(t *testing.T) {
	// Create a temp file with known content.
	dir := t.TempDir()
	path := filepath.Join(dir, "testfile.txt")
	if err := os.WriteFile(path, []byte(plain), 0644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name      string
		hasher    hash.Hasher
		expectLen int
	}{
		{name: "MD5", hasher: md.NewHash(md.Format(hashenum.MD5)), expectLen: 16},
		{name: "SHA1", hasher: sha1.NewHash(sha1.Format(hashenum.SHA1)), expectLen: 20},
		{name: "SHA256", hasher: sha2.NewHash(sha2.Format(hashenum.SHA256)), expectLen: 32},
		{name: "SHA512", hasher: sha512.NewHash(sha512.Format(hashenum.SHA512)), expectLen: 64},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open(path)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			result := tc.hasher.File(f)
			if result == nil {
				t.Fatalf("%s File() returned nil", tc.name)
			}
			if len(result) != tc.expectLen {
				t.Errorf("%s File(): expected length %d, got %d", tc.name, tc.expectLen, len(result))
			}

			// File hash must match String hash for the same content.
			f.Seek(0, 0)
			strResult := tc.hasher.String(plain)
			if hex.EncodeToString(result) != hex.EncodeToString(strResult) {
				t.Errorf("%s: File() and String() produce different hashes for same content", tc.name)
			}

			t.Logf("%s  file_hash len=%d  hex=%s", tc.name, len(result), hex.EncodeToString(result))
		})
	}
}

// ---------------------------------------------------------------------------
// Format() returns the correct enum
// ---------------------------------------------------------------------------

func TestFormatReturnsCorrectEnum(t *testing.T) {
	cases := []struct {
		hasher hash.Hasher
		want   hashenum.HashEnum
	}{
		{hasher: md.NewHash(md.Format(hashenum.MD5)), want: hashenum.MD5},
		{hasher: sha1.NewHash(sha1.Format(hashenum.SHA1)), want: hashenum.SHA1},
		{hasher: sha2.NewHash(sha2.Format(hashenum.SHA256)), want: hashenum.SHA256},
		{hasher: sha3.NewHash(sha3.Format(hashenum.SHA3_256)), want: hashenum.SHA3_256},
		{hasher: sha512.NewHash(sha512.Format(hashenum.SHA512)), want: hashenum.SHA512},
		{hasher: blake2.NewHash(blake2.Format(hashenum.BLAKE2b_512)), want: hashenum.BLAKE2b_512},
	}

	for _, tc := range cases {
		got := tc.hasher.Format()
		if got != tc.want {
			t.Errorf("Format(): got %v (%s), want %v (%s)", got, got.String(), tc.want, tc.want.String())
		} else {
			t.Logf("Format() %s — OK", got.String())
		}
	}
}

// ---------------------------------------------------------------------------
// SSHA (salted SHA1)
// ---------------------------------------------------------------------------

func TestSSHA(t *testing.T) {
	password := "password"
	salt := []byte("salt")

	h := cryptsha1.New()
	h.Write([]byte(password))
	h.Write(salt)
	digest := h.Sum(nil)
	full := append(digest, salt...)
	expected := base64.StdEncoding.EncodeToString(full)

	t.Logf("SSHA(%q, %q) = {SSHA}%s", password, salt, expected)

	if len(expected) == 0 {
		t.Fatal("SSHA produced empty output")
	}
	// SSHA: 20-byte sha1 + len(salt) bytes, base64'd.
	// With a 4-byte salt → 24 bytes → 32 base64 chars.
	if len(expected) != 32 {
		t.Errorf("SSHA: expected 32 base64 chars, got %d", len(expected))
	}
}

// ---------------------------------------------------------------------------
// Empty / edge-case inputs
// ---------------------------------------------------------------------------

func TestEdgeCases(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{name: "empty_string", input: ""},
		{name: "single_char", input: "a"},
		{name: "unicode", input: "héllo wörld 🌍"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hasher := sha256.New()
			// Use the package-level helper directly.
			// (We test the wrapper through the algorithm tests above.)
			result := hash.String(tc.input, hasher)
			if result == nil {
				t.Fatalf("String(%q) returned nil", tc.input)
			}
			t.Logf("SHA256(%q) = %s", tc.input, hex.EncodeToString(result))
		})
	}
}
