package stringut

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	log "github.com/colt3k/nglog/ng"
)

//ToChar change from int to char
func ToChar(i int) rune {
	return rune('A' - 1 + i)
}

//ToCharStr change from int to string
func ToCharStr(i int) string {
	return string('A' - 1 + i)
}

//Has dummy first character so you don't have to subtract one
var arr = [...]string{".", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
	"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

//ToCharStrArr change from int to array
func ToCharStrArr(i int) string { return arr[i] }

const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

//ToCharStrConst change from int to string constant
func ToCharStrConst(i int) string {
	return abc[i-1 : i]
}

// ValidString confirms that the passed in string is UTF-8
func ValidString(s string) bool {
	return utf8.ValidString(s)
}

// Reverse returns its argument string reversed rune-wise left to right.
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func HRByteCount(bytes int64, si bool) string {

	var unit int64

	if si {
		unit = 1000
	} else {
		unit = 1024
	}

	if bytes < unit {
		out := strconv.FormatInt(bytes, 10) + " B"
		return out
	}
	exp := int(math.Log(float64(bytes)) / math.Log(float64(unit)))

	var pre string

	if si {
		pre = string("kMGTPE"[exp-1]) + ""
	} else {
		pre = string("KMGTPE"[exp-1]) + "i"
	}

	//log.Printf("%.1f %sB", float64(f.size) / math.Pow(float64(unit), float64(exp)), pre)
	out := fmt.Sprintf("%.1f %sB", float64(bytes)/math.Pow(float64(unit), float64(exp)), pre)
	return out
}

func ToBool(str string) (bool, error) {

	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "y", "Y":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "n", "N":
		return false, nil
	}
	return false, &strconv.NumError{Func: "ParseBool", Num: str, Err: strconv.ErrSyntax}
}

func SubString(str string, start, end int) string {
	rns := []rune(str)
	tmp := string(rns[start:end])
	return tmp
}

func SubAfterIndexString(str string, start int) string {
	rns := []rune(str)
	tmp := string(rns[start:])
	return tmp
}
func SubBeforeIndexString(str string, end int) string {
	rns := []rune(str)
	tmp := string(rns[:end])
	return tmp
}

func ToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Logf(log.WARN, "issue converting to int\n%+v", err)
	}
	return i
}

// REGEX https://yourbasic.org/golang/regexp-cheat-sheet/

// ExtractDomainFromURL
func ExtractDomainFromURL(url string) string {
	re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)
	re.MatchString(url)
	submatchall := re.FindAllString(url, -1)
	for _, element := range submatchall {
		return element
	}
	return ""
}

// ExtractTextBetweenBrackets
func ExtractTextBetweenBrackets(str1 string) []string {
	found := make([]string, 0)

	re := regexp.MustCompile(`\[([^\[\]]*)\]`)

	submatchall := re.FindAllString(str1, -1)
	for _, element := range submatchall {
		element = strings.Trim(element, "[")
		element = strings.Trim(element, "]")
		//fmt.Println(element)
		found = append(found, element)
	}
	return found
}

//ExtractAllNonAlphaNumeric
func ExtractAllNonAlphaNumeric(str1 string) []string {
	found := make([]string, 0)
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	submatchall := re.FindAllString(str1, -1)
	for _, element := range submatchall {
		found = append(found, element)
	}
	return found
}

// ExtractDateFromString
func ExtractDateFromString(str1 string) string {
	var found string
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	submatchall := re.FindAllString(str1, -1)
	for _, element := range submatchall {
		found = element
	}
	return found
}

// ExtractDNSorIPAddress
func ExtractDNSorIPAddress(str1 string) []string {

	found := make([]string, 0)
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)

	submatchall := re.FindAllString(str1, -1)
	for _, element := range submatchall {
		found = append(found, element)
	}
	return found
}

// ValidateEmail
func ValidateEmail(str1 string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(str1)
}

/* ValidatePhoneNumber
Numerous formats are supported
*/
func ValidatePhoneNumber(str1 string) bool {

	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return re.MatchString(str1)
}

// ValidateDateFormat dd/mm/yyyy
func ValidateDateFormat(str1 string) bool {

	re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])/(0?[1-9]|1[012])/((19|20)\\d\\d)")
	return re.MatchString(str1)
}

/* ValidateCCNumbers - no dashes allowed
Regular expression validate Visa, MasterCard, American Express, Diners Club, Discover, and JCB cards
*/
func ValidateCCNumbers(str1 string) bool {
	re := regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`)
	return re.MatchString(str1)
}

// ReplaceAnyNonAlphaNumericWithX
func ReplaceAnyNonAlphaNumericWithX(str1, repl string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		log.Logf(log.FATAL, "issue regex\n%+v", err)
	}
	return reg.ReplaceAllString(str1, repl)
}

/* ReplaceFirstOccurenceOfX
Need to figure a way to reuse this
*/
func ReplaceFirstOccurenceOfX(str1, repl string) string {
	strEx := "Php-Golang-Php-Python-Php-Kotlin"
	reStr := regexp.MustCompile("^(.*?)Php(.*)$")
	repStr := "${1}Java$2"
	return reStr.ReplaceAllString(strEx, repStr)
}

func ContainsOnlyNumeric(str string, length int) bool {

	reStr := regexp.MustCompile("^\\d{" + strconv.Itoa(length) + "}$")
	return reStr.MatchString(str)
}

func Eq(one, two string) bool {
	return strings.EqualFold(one, two)
}
