package decrypt

import (
	"testing"
)

func Test_GPG_Decrypt_Service_Modify(t *testing.T) {
	config := DefaultConfig()
	config.Pass = "foo"
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("hello world")
	value := []byte(`-----BEGIN PGP SIGNATURE-----

wx4EBwMItflyy+CkVHfgNJ9CwJz0SXR8DVmT+GrIQpbSPAFfOlMN/2J8XF2/hRCm
oHm+HyYpiGLqnC/rncq3SRJ9z0xSEbhS5l+Dp3xMTGniaNEU2xtt72M35kS+HA==
=R5GK
-----END PGP SIGNATURE-----`) // "hello world"
	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if string(modified) != string(expected) {
		t.Fatal("expected", expected, "got", modified)
	}
}

func Test_GPG_Decrypt_Service_Modify_Empty(t *testing.T) {
	config := DefaultConfig()
	config.Pass = "foo"
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("")
	value := []byte("")
	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if string(modified) != string(expected) {
		t.Fatal("expected", expected, "got", modified)
	}
}
