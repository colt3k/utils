package stringut

import "testing"

func TestReverse(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}
	for _, c := range cases {
		got := Reverse(c.in)
		if got != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestContainsOnlyNumeric(t *testing.T) {
	cases := []struct {
		in     string
		length int
		want   bool
	}{
		{"1", 1, true},
		{"13", 2, true},
		{"13", 3, false},
		{"098765", 5, false},
		{"099765", 6, true},
		{"test0987", 8, false},
		{"0t8e9a2b", 6, false},
	}
	for _, c := range cases {
		got := ContainsOnlyNumeric(c.in, c.length)
		if got != c.want {
			t.Errorf("ContainsOnlyNumeric(%q, %d) == %v, want %v", c.in, c.length, got, c.want)
		}
	}
}
