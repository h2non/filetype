package matchers

import "testing"

func TestArExcludesDebianPackages(t *testing.T) {
	deb := []byte("!<arch>\ndebian-binary")
	if !Deb(deb) {
		t.Fatal("Deb should match debian-binary member")
	}
	if Ar(deb) {
		t.Fatal("Ar should not match a Debian package")
	}

	plain := []byte("!<arch>\nfoo/")
	if Deb(plain) {
		t.Fatal("Deb should not match a generic ar archive")
	}
	if !Ar(plain) {
		t.Fatal("Ar should match a generic ar archive")
	}
}
