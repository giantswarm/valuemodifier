package inspect

import (
	"testing"
)

func Test_Inspect_Service_Inspect(t *testing.T) {
	config := DefaultConfig()
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("bar")
	value := []byte("bar")

	inspectedValue := newService.InspectedValue()
	if string(inspectedValue) != string("") {
		t.Fatal("expected", "empty string", "got", inspectedValue)
	}

	_, err = newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	inspectedValue = newService.InspectedValue()
	if string(inspectedValue) != string(expected) {
		t.Fatal("expected", expected, "got", inspectedValue)
	}
}

func Test_Inspect_Service_Modify(t *testing.T) {
	config := DefaultConfig()
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("bar")
	value := []byte("bar")

	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if string(modified) != string(expected) {
		t.Fatal("expected", expected, "got", modified)
	}
}
