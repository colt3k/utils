# encryption



## Abilities

### Key ALGORITHMS

1. PBKDF2 build key off password
2. Scrypt build key off password (better alternative to #1)


### ENCRYPTION ALGORITHM/TRANSFORMATION

1. AES/CBC/PKCS5Padding Encrypt/Decrypt using prior keys (with AES PKCS7Padding is actually used)


### MAC

1. Poly1305


#### TERMS

- Password is passed by user

- Salt is a unique salt for the system this will be running on and doesn't change

- IV or nonce is generated on every encryption and is part of the combination of a key return

##### NOTES

[golang-crypto](https://github.com/chain/golang-crypto)

[exploring-1password-crypto](http://sosedoff.com/2015/05/30/exploring-1password-crypto.html)

[data-encryption-in-go-using-openssl](http://sosedoff.com/2015/05/22/data-encryption-in-go-using-openssl.html)

[golang cipher](https://golang.org/pkg/crypto/cipher/)

