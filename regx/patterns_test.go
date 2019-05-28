package regx

import (
	"fmt"
	"reflect"
	"testing"
)


var data = []struct{
	pattern	string
	dataGlob string
	valid bool
}{
	{ALPHANUMERIC, "abc5678def", true},
	{ALPHA, "abcdef", true},
	{BTC, "__1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2^&", true},
	{DIGITS, "abc5678f", true},
	{DOMAIN, "www.test.com", true},
	{EMAIL, "none@none.com", true},
	{GITREPO, "ssh://user@host.xz:port/path/to/repo.git/", true},
	{GITREPO, "ssh://user@host.xz/path/to/repo.git/", true},
	{GITREPO, "user@host.xz:/path/to/repo.git/", false},
	{GITREPO, "host.xz:/path/to/repo.git/", false},
	{GITREPO, "rsync://host.xz/path/to/repo.git/", false},
	{GITREPO, "git://host.xz/~user/path/to/repo.git/", true},
	{GITREPO, "/path/to/repo.git/", false},
	{GITREPO, "~/path/to/repo.git", false},
	{GITREPO, "file:///path/to/repo.git/", false},
	{GITREPO, "file://~/path/to/repo.git/", false},
	{HEXCOLOR, "#999999", true},
	//{HTMLTAG, "&lt;div&gt;", true},
	{IBAN, "BE71096123456769", true},
	{IBAN, "FR7630006000011234567890189", true},
	{IBAN, "DE91100000000123456789", true},
	{IBAN, "GR9608100010000001234567890", true},
	{IBAN, "RO09BCYP0000001234567890", true},
	{IBAN, "SA4420000001234567891234", true},
	{IBAN, "ES7921000813610123456789", true},
	{IBAN, "CH5604835012345678009", true},
	{IBAN, "GB98MIDL07009312345678", true},
	{IP, "255.255.255.0", true},
	{IP, "192.168.2.999", false},
	{IPV4, "255.255.255.0", true},
	{IPV4, "192.168.2.999", false},
	{IPV6, "fe80::8b0:8fb0:c02f:3aee%en0", true},
	{MACADDRESS, "66:00:dc:50:04:01", true},
	{MD5, "7f1c3bbd142f7da9ff5535f0311d728c", true},
	//{PASSWORD, "Abc9Ta*", true},
	{PHONENUMBER, "6305555555", true},
	{POBOX, "PO Box 098", true},
	{SHA1, "FDF49D0CB812BEA37AD500DEA2863128A2FB715D", true},
	{SHA256, "6235BB10721B13A4AC93881E840CE1A086DCB63F19934E3596D97A8CAFE240AE", true},
	{SSN, "111-11-1111", true},
	{STREETADDRESS, "123 Some Lane", true},
	{TIME, "1:30:20 PM", true},
	{URLSLUG, "www.home.com", true},
	//{URL, "https://www.home.com", true},
	{USSTATEABBRV, "WY", true},
	{USSTATE, "Wyoming", true},
	{USERNAME, "joe", true},
	{UUID, "ef38562d-e420-4e2a-a02e-0d4e0c9b51cd", true},
	{ZIPCODE, "11111", true},
}

func Test(t *testing.T) {
	for _,k := range data {
		fmt.Println( k.dataGlob)
		if Match(k.pattern, k.dataGlob) && k.valid {
			//fmt.Println("valid", k.dataGlob)
		} else if !Match(k.pattern, k.dataGlob) && !k.valid {
			//fmt.Println("valid", k.dataGlob)
		} else if Match(k.pattern, k.dataGlob) && !k.valid {

			fmt.Println("INvalid", k.dataGlob, reflect.TypeOf(k.pattern).Name())
			t.Fail()
		} else if !Match(k.pattern, k.dataGlob) && k.valid {
			fmt.Println("INvalid", k.dataGlob, reflect.TypeOf(k.pattern).Name())
			t.Fail()
		}
	}
}
