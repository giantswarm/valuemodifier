package overwrite

import (
	"testing"
)

func Test_Overwrite_Service_Modify(t *testing.T) {
	config := DefaultConfig()
	config.NewValue = "foo"
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("foo")
	value := []byte("bar")

	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if string(modified) != string(expected) {
		t.Fatal("expected", expected, "got", modified)
	}
}
