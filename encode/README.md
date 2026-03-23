# Encoding

Simple encoding helpers for standard base64, URL-safe base64, and hex.

## Supported Encodings

- `encodeenum.B64STD`
- `encodeenum.B64URL`
- `encodeenum.Hex`

## Example

```go
string := encode.Encode([]byte("some text"), encodeenum.B64STD)
[]byte := encode.Decode([]byte(string), encodeenum.B64STD)

string := encode.Encode([]byte("some text"), encodeenum.B64URL)
[]byte := encode.Decode([]byte(string), encodeenum.B64URL)

string := encode.Encode([]byte("some text"), encodeenum.Hex)
[]byte := encode.Decode([]byte(string), encodeenum.Hex)
```

`B64DecodeStdSanitized` is useful when the input may contain surrounding quotes or backticks.
