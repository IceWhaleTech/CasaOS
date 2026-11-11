package samba

import "testing"

func TestMountOptions(t *testing.T) {
	cases := []struct {
		name, username, password, port, want string
	}{
		{"default port", "casa", "secret", "445", "username=casa,password=secret,iocharset=utf8,port=445"},
		{"custom port", "casa", "secret", "4455", "username=casa,password=secret,iocharset=utf8,port=4455"},
		{"no port", "casa", "secret", "", "username=casa,password=secret,iocharset=utf8"},
		{"comma in password", "casa", "a,b,,c", "445", "username=casa,password=a,,b,,,,c,iocharset=utf8,port=445"},
	}
	for _, c := range cases {
		if got := MountOptions(c.username, c.password, c.port); got != c.want {
			t.Errorf("%s: MountOptions(%q, %q, %q) = %q, want %q", c.name, c.username, c.password, c.port, got, c.want)
		}
	}
}
