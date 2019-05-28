package regx

import (
	"regexp"
)

// https://golang.org/pkg/regexp/syntax/

const (
	ALPHANUMERIC = "[a-zA-Z0-9]+"				// Alphanumeric Patterns
	ALPHA = "([a-zA-Z])+"						// Alpha Patterns

	BTC = "[13][a-km-zA-HJ-NP-Z1-9]{25,34}"		// Bitcoin Address

	// DIGITS
	DIGITS = "[0-9]+"
	// DOMAIN
	DOMAIN = "([a-z][a-z0-9-]+(\\.|-*\\.))+[a-z]{2,6}"
	EMAIL = "[a-zA-Z0-9._%-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,4}"
	GITREPO = "((git|ssh|http(s)?)|(git@[\\w\\.]+))(:(//)?)([\\w\\.@\\:/\\-~]+)(\\.git)(/)?"
	HEXCOLOR = "#([0-9a-f]{3,6})"
	HTMLTAG = "<([a-z]+)([^<]+)*(?:>(.*)<\\/\\1>|\\s+\\/>)"

	IBAN = "[a-zA-Z]{2}[0-9]{2}[a-zA-Z0-9]{4}[0-9]{7}([a-zA-Z0-9]?){0,16}" // International Bank account numbers

	//IP = "(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
	IP = "\\A(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\z"
	//IPV4 = "((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])"
	IPV4 = "\\A(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\z"
	IPV6 = "(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))"
	MACADDRESS = "([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})"

	MD5 = "\\b[A-Fa-f0-9]{32}\\b"
	// One lower, one uppercase letter one number and 6 long
	//PASSWORD = "(?=^.{6,}$)((?=.*[A-Za-z0-9])(?=.*[A-Z])(?=.*[a-z]))^.*"
	//PASSWORD = "(^.{6,}$)((^.*[A-Za-z0-9]))^.*"
	//PASSWORD = "(^.{6,}$)((?=.*[A-Za-z0-9])(?=.*[A-Z])(?=.*[a-z]))^.*"
	PHONENUMBER = "\\b(([0-9]{1})*[- .(]*([0-9]{3})[- .)]*[0-9]{3}[- .]*[0-9]{4})+\\b"
	POBOX = "\\b[PO.|Post\\sOffice]*\\s?B(ox)?.*\\d+\\b"

	SHA1 = "\\b[A-Fa-f0-9]{7,40}\\b"
	SHA256 = "\\b[A-Fa-f0-9]{64}\\b"
	SSN = "\\d{3}-\\d{2}-\\d{4}"
	STREETADDRESS = "\\d+[ ](?:[A-Za-z0-9.-]+[ ]?)+(?:Avenue|Lane|Road|Boulevard|Drive|Street|Ave|Dr|Rd|Blvd|Ln|St)\\.?"
	TIME = "([0[0-9]|1[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?(\\s?[PA]M)?"

	URLSLUG = "[a-z0-9-]+"
	//URL = "(https?:\\/\\/(?:www\\.|(?!www))[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\\.[^\\s]{2,}|www\\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\\.[^\\s]{2,}|https?:\\/\\/(?:www\\.|(?!www))[a-zA-Z0-9]\\.[^\\s]{2,}|www\\.[a-zA-Z0-9]\\.[^\\s]{2,})"
	USSTATEABBRV = "AL|AK|AS|AZ|AR|CA|CO|CT|DE|DC|FM|FL|GA|GU|HI|ID|IL|IN|IA|KS|KY|LA|ME|MH|MD|MA|MI|MN|MS|MO|MT|NE|NV|NH|NJ|NM|NY|NC|ND|MP|OH|OK|OR|PW|PA|PR|RI|SC|SD|TN|TX|UT|VT|VI|VA|WA|WV|WI|WY"
	USSTATE = "Alabama|Alaska|Arizona|Arkansas|California|Colorado|Connecticut|Delaware|Florida|Georgia|Hawaii|Idaho|Illinois|Indiana|Iowa|Kansas|Kentucky|Louisiana|Maine|Maryland|Massachusetts|Michigan|Minnesota|Mississippi|Missouri|Montana|Nebraska|Nevada|New[ ]Hampshire|New[ ]Jersey|New[ ]Mexico|New[ ]York|North[ ]Carolina|North[ ]Dakota|Ohio|Oklahoma|Oregon|Pennsylvania|Rhode[ ]Island|South[ ]Carolina|South[ ]Dakota|Tennessee|Texas|Utah|Vermont|Virginia|Washington|West[ ]Virginia|Wisconsin|Wyoming"
	USERNAME = "[a-z0-9_]"
	UUID = "[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}"
	ZIPCODE = "\\b\\d{5}(?:-\\d{4})?\\b"
)

func Find(regex, data string) string {
	r, _ := regexp.Compile(regex)
	return r.FindString(data)
}
func Match(regex, data string) bool {
	r := regexp.MustCompile(regex)
	return r.MatchString(data)
}
func CC(cardType, data string) string {
	pattern := "(?:4[0-9]{12}(?:[0-9]{3})?)"

	switch cardType {
	case "mastercard":
		pattern = "(?:5[1-5][0-9]{14})"
	case "americanexpress":
		fallthrough
	case "amex":
		pattern = "(?:3[47][0-9]{13})"
	case "discover":
		pattern = "(?:6(?:011|5[0-9][0-9])[0-9]{12})"
	}
	r, _ := regexp.Compile(pattern)
	return r.FindString(data)
}

func DATE(format, data string) string {

	DD := "(0?[1-9]|[12][0-9]|3[01])"
	MM := "(0?[1-9]|1[012])"
	YYYY := "\\d{4}"
	pattern := "[DD, MM, YYYY]"	// default


	switch format {
	case "MM/DD/YYYY":
		fallthrough
	case "MM-DD-YYYY":
		fallthrough
	case "MM.DD.YYYY":
		pattern = "["+MM+", "+DD+", "+YYYY+"]"
	case "YYYY/MM/DD":
		fallthrough
	case "YYYY-MM-DD":
		fallthrough
	case "YYYY.MM.DD":
		pattern = "["+YYYY+", "+MM+", "+DD+"]"
	}

	r, _ := regexp.Compile(pattern)
	return r.FindString(data)
}

func ISBN(format, data string) string {
	pattern := "(?:ISBN(?:-10)?:?\\s)?(?=[0-9X]{10}$|(?=(?:[0-9]+[-\\s]){3})[-\\s0-9X]{13}$)[0-9]{1,5}[-\\s]?[0-9]+[-\\s]?[0-9]+[-\\s]?[0-9X]"

	if format != "ISBN-10" {
		pattern = "(?:ISBN(?:-13)?:?\\s)?(?=[0-9]{13}$|(?=(?:[0-9]+[-\\s]){4})[-\\s0-9]{17}$)97[89][-\\s]?[0-9]{1,5}[-\\s]?[0-9]+[-\\s]?[0-9]+[-\\s]?[0-9]"
	}

	r, _ := regexp.Compile(pattern)
	return r.FindString(data)
}

func POSTALCODE(format, data string) string {
	pattern := "\b\\d{5}(?:-\\d{4})?\b"

	switch format {
	case "AD":
		pattern = "(?:AD)*(\\d{3})"
	case "AM":
		fallthrough
	case "BY":
		fallthrough
	case "CN":
		fallthrough
	case "IN":
		fallthrough
	case "KG":
		fallthrough
	case "KP":
		fallthrough
	case "KZ":
		fallthrough
	case "MN":
		fallthrough
	case "NG":
		fallthrough
	case "RO":
		fallthrough
	case "RS":
		fallthrough
	case "RU":
		fallthrough
	case "SG":
		fallthrough
	case "TJ":
		fallthrough
	case "TM":
		fallthrough
	case "UZ":
		fallthrough
	case "VN":
		pattern = "(\\d{6})"
	case "AR":
		pattern = "([A-Z]\\d{4}[A-Z]{3})"
	case "AT":
		fallthrough
	case "AU":
		fallthrough
	case "BD":
		fallthrough
	case "BE":
		fallthrough
	case "BG":
		fallthrough
	case "CH":
		fallthrough
	case "CR":
		fallthrough
	case "CV":
		fallthrough
	case "CX":
		fallthrough
	case "CY":
		fallthrough
	case "DK":
		fallthrough
	case "ET":
		fallthrough
	case "GE":
		fallthrough
	case "GL":
		fallthrough
	case "GW":
		fallthrough
	case "HU":
		fallthrough
	case "LI":
		fallthrough
	case "LR":
		fallthrough
	case "LU":
		fallthrough
	case "MK":
		fallthrough
	case "MZ":
		fallthrough
	case "NE":
		fallthrough
	case "NF":
		fallthrough
	case "NO":
		fallthrough
	case "NZ":
		fallthrough
	case "PH":
		fallthrough
	case "PY":
		fallthrough
	case "TN":
		fallthrough
	case "VE":
		fallthrough
	case "ZA":
		pattern = "(\\d{4})"
	case "AX":
	case "FI":
		pattern = "(?:FI)*(\\d{5})"
	case "AZ":
		pattern = "(?:AZ)*(\\d{4})"
	case "BA":
		fallthrough
	case "CZ":
		fallthrough
	case "DE":
		fallthrough
	case "DO":
		fallthrough
	case "DZ":
		fallthrough
	case "EE":
		fallthrough
	case "EG":
		fallthrough
	case "ES":
		fallthrough
	case "FM":
		fallthrough
	case "FR":
		fallthrough
	case "GR":
		fallthrough
	case "GT":
		fallthrough
	case "ID":
		fallthrough
	case "IL":
		fallthrough
	case "IQ":
		fallthrough
	case "IT":
		fallthrough
	case "JO":
		fallthrough
	case "KE":
		fallthrough
	case "KH":
		fallthrough
	case "KW":
		fallthrough
	case "LA":
		fallthrough
	case "LK":
		fallthrough
	case "MA":
		fallthrough
	case "MC":
		fallthrough
	case "ME":
		fallthrough
	case "MM":
		fallthrough
	case "MQ":
		fallthrough
	case "MV":
		fallthrough
	case "MX":
		fallthrough
	case "MY":
		fallthrough
	case "NC":
		fallthrough
	case "NP":
		fallthrough
	case "PK":
		fallthrough
	case "PL":
		fallthrough
	case "SA":
		fallthrough
	case "SD":
		fallthrough
	case "SK":
		fallthrough
	case "SN":
		fallthrough
	case "TH":
		fallthrough
	case "TR":
		fallthrough
	case "TW":
		fallthrough
	case "UA":
		fallthrough
	case "UY":
		fallthrough
	case "VA":
		fallthrough
	case "YT":
		fallthrough
	case "ZM":
		fallthrough
	case "CS":
		pattern = "(\\d{5})"
	case "BB":
		pattern = "(?:BB)*(\\d{5})"
	case "BH":
		pattern = "(\\d{3}\\d?)"
	case "BM":
		pattern = "([A-Z]{2}\\d{2})"
	case "BN":
		fallthrough
	case "HN":
		pattern = "([A-Z]{2}\\d{4})"
	case "BR":
		pattern = "(\\d{8})"
	case "CA":
		pattern = "([ABCEGHJKLMNPRSTVXY]\\d[ABCEGHJKLMNPRSTVWXYZ]) ?(\\d[ABCEGHJKLMNPRSTVWXYZ]\\d)"
	case "CL":
		fallthrough
	case "JP":
		fallthrough
	case "NI":
		fallthrough
	case "PT":
		pattern = "(\\d{7})"
	case "CU":
		pattern = "(?:CP)*(\\d{5})"
	case "EC":
		pattern = "([a-zA-Z]\\d{4}[a-zA-Z])"
	case "FO":
		pattern = "(?:FO)*(\\d{3})"
	case "GB":
		fallthrough
	case "GG":
		fallthrough
	case "IM":
		fallthrough
	case "JE":
		pattern = "([Gg][Ii][Rr] 0[Aa]{2})|((([A-Za-z][0-9]{1,2})|(([A-Za-z][A-Ha-hJ-Yj-y][0-9]{1,2})|(([A-Za-z][0-9][A-Za-z])|([A-Za-z][A-Ha-hJ-Yj-y][0-9]?[A-Za-z]))))\\s?[0-9][A-Za-z]{2})"
	case "GF":
		pattern = "((97|98)3\\d{2})"
	case "GP":
		pattern = "((97|98)\\d{3})"
	case "GU":
		pattern = "(969\\d{2})"
	case "HR":
		pattern = "(?:HR)*(\\d{5})"
	case "HT":
		pattern = "(?:HT)*(\\d{4})"
	case "IR":
		pattern = "(\\d{10})"
	case "IS":
		fallthrough
	case "LS":
		fallthrough
	case "MG":
		fallthrough
	case "OM":
		fallthrough
	case "PG":
		pattern = "(\\d{3})"
	case "KR":
		pattern = "(?:SEOUL)*(\\d{6})"
	case "LB":
		pattern = "(\\d{4}(\\d{4})?)"
	case "LT":
		pattern = "(?:LT)*(\\d{5})"
	case "LV":
		pattern = "(?:LV)*(\\d{4})"
	case "MD":
		pattern = "(?:MD)*(\\d{4})"
	case "MT":
		pattern = "([A-Z]{3}\\s\\d{2}\\d?)"
	case "NL":
		pattern = "(\\d{4}[A-Z]{2})"
	case "PF":
		pattern = "((97|98)7\\d{2})"
	case "PM":
		pattern = "(97500)"
	case "PR":
		pattern = "(\\d{9})"
	case "PW":
		pattern = "(96940)"
	case "RE":
		pattern = "((97|98)(4|7|8)\\d{2})"
	case "SE":
		pattern = "(?:SE)*(\\d{5})"
	case "SH":
		pattern = "(STHL1ZZ)"
	case "SI":
		pattern = "(?:SI)*(\\d{4})"
	case "SM":
		pattern = "(4789\\d)"
	case "SO":
		pattern = "([A-Z]{2}\\d{5})"
	case "SV":
		pattern = "(?:CP)*(\\d{4})"
	case "SZ":
		pattern = "([A-Z]\\d{3})"
	case "TC":
		pattern = "(TKCA 1ZZ)"
	case "US":
		fallthrough
	case "VI":
		pattern = "\\d{5}(-\\d{4})?"
	case "WF":
		pattern = "(986\\d{2})"
	}
	r, _ := regexp.Compile(pattern)
	return r.FindString(data)
}

func PRICE(data string, matchCurrency bool) string {
	pattern := "(\\d+[,\\.\\s]+\\d+)\\b"
	if matchCurrency {
		pattern = "([\\$\\xA2-\\xA5\\u058F\\u060B\\u09F2\\u09F3\\u09FB\\u0AF1\\u0BF9\\u0E3F\\u17DB\\u20A0-\\u20BD\\uA838\\uFDFC\\uFE69\\uFF04\\uFFE0\\uFFE1\\uFFE5\\uFFE6]?\\s?\\d+[,\\.\\s]+\\d+)\\b"
	}
	r, _ := regexp.Compile(pattern)
	return r.FindString(data)
}