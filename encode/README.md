# Encoding

### Example
    
    string := encode.Encode([]byte("some text"), encodeenum.B64STD)
    []byte := encode.Decode([]byte(string), encodeenum.B64STD)
    
    string := encode.Encode([]byte("some text"), encodeenum.B64URL)
    []byte := encode.Decode([]byte(string), encodeenum.B64URL)
    
    string := encode.Encode([]byte("some text"), encodeenum.Hex)
    []byte := encode.Decode([]byte(string), encodeenum.Hex)