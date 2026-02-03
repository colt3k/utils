package ss

/**
This is an example setup using BasicAuth
	In this example you set a single user(you can add more by using a DB)
	In the test there is a pre-set token, meant to be replaced at build time
*/
import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	sessions = map[string]session{}
)

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

// DynamicAuth this should be used to return some rotating piece of data i.e. totp
func DynamicAuth() string {
	t := time.Now()
	minutes := ReverseNumber(t.Minute())
	return minutes
}

// ReverseNumber take a number and reverse it keeping zero and return as string
func ReverseNumber(num int) string {
	tmp := strconv.Itoa(num)
	if num < 10 {
		tmp = "0" + tmp
	}
	runes := []rune(tmp)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
func validateCookie(r *http.Request, cookieName string) bool {
	if len(strings.TrimSpace(cookieName)) == 0 {
		cookieName = "session_token"
	}
	c, err := r.Cookie(cookieName)
	if err != nil {
		// For any other type of error, return a bad request status
		return false
	}
	sessionToken := c.Value
	// We then get the session from our session map
	_, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		return false
	}
	return true
}

// BasicAuth handler func that performs basic auth
func BasicAuth(next http.HandlerFunc, authToken, cookieName, user string, sessionExpirationSeconds int, dynamicAuth interface{}) http.HandlerFunc {

	if len(strings.TrimSpace(cookieName)) == 0 {
		cookieName = "session_token"
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract the username and password from the request
		// Authorization header. If no Authentication header is present
		// or the header value is invalid, then the 'ok' return value
		// will be false.

		dynVal := dynamicAuth.(func() string)()
		validCookie := validateCookie(r, cookieName)
		if validCookie {
			next.ServeHTTP(w, r)
			return
		}
		username, password, ok := r.BasicAuth()
		if ok {
			// Calculate SHA-256 hashes for the provided and expected
			// usernames and passwords.
			//fmt.Printf("current time %v\n", t)
			//fmt.Printf("current minutes %v\n", t.Minute())

			//fmt.Printf("received user: %v - pass: %v\n", username, password)
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(user))
			expectedPasswordHash := sha256.Sum256([]byte(authToken + dynVal))
			//fmt.Printf("expected user: %v - pass: %v\n", util.USER, util.PASS+minutes)
			// Use the subtle.ConstantTimeCompare() function to check if
			// the provided username and password hashes equal the
			// expected username and password hashes. ConstantTimeCompare
			// will return 1 if the values are equal, or 0 otherwise.
			// Importantly, we should to do the work to evaluate both the
			// username and password before checking the return values to
			// avoid leaking information.
			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			// If the username and password are correct, then call
			// the next handler in the chain. Make sure to return
			// afterwards, so that none of the code below is run.
			if usernameMatch && passwordMatch {
				sessionToken := uuid.New().String()
				// expire after an hour
				expiresAt := time.Now().Add(time.Duration(sessionExpirationSeconds) * time.Second)
				sessions[sessionToken] = session{
					username: sessionToken,
					expiry:   expiresAt,
				}
				http.SetCookie(w, &http.Cookie{
					Name:    cookieName,
					Value:   sessionToken,
					Expires: expiresAt,
				})
				next.ServeHTTP(w, r)
				return
			}
		}

		// If the Authentication header is not present, is invalid, or the
		// username or password is wrong, then set a WWW-Authenticate
		// header to inform the client that we expect them to use basic
		// authentication and send a 401 Unauthorized response.
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
