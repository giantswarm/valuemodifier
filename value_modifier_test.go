package valuemodifier

import (
	"fmt"
	"testing"
)

type testModifier struct{}

func (m testModifier) Modify(value []byte) ([]byte, error) {
	return []byte(string(value) + "-modified"), nil
}

func Test_ValueModifier_TraverseJSON(t *testing.T) {
	var err error
	var newService *Service
	{
		config := DefaultConfig()
		config.ValueModifiers = []ValueModifier{
			testModifier{},
		}
		config.IgnoreFields = []string{
			"noSecret1",
			"noSecret2",
		}
		newService, err = New(config)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	var expectedJSON string
	var testJSON string
	{
		testJSON = `{
  "block1": {
    "block11": {
      "pass1": "pass1"
    },
    "pass2": "pass2"
  },
  "block2": {
    "block21": {
      "pass3": "pass3"
    },
    "pass4": "pass4"
  },
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass5": "pass5",
  "pass6": 123456
}`
		expectedJSON = `{
  "block1": {
    "block11": {
      "pass1": "pass1-modified"
    },
    "pass2": "pass2-modified"
  },
  "block2": {
    "block21": {
      "pass3": "pass3-modified"
    },
    "pass4": "pass4-modified"
  },
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass5": "pass5-modified",
  "pass6": "123456-modified"
}`
	}

	{
		output, err := newService.TraverseJSON([]byte(testJSON))
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		if string(output) != expectedJSON {
			t.Fatal("expected", fmt.Sprintf("%q", expectedJSON), "got", fmt.Sprintf("%q", output))
		}
	}
}

func Test_ValueModifier_TraverseYAML(t *testing.T) {
	var err error
	var newService *Service
	{
		config := DefaultConfig()
		config.ValueModifiers = []ValueModifier{
			testModifier{},
		}
		config.IgnoreFields = []string{
			"noSecret1",
			"noSecret2",
		}
		newService, err = New(config)
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
	}

	var expectedYAML string
	var testYAML string
	{
		testYAML = `block1:
  block11:
    pass1: pass1
  pass2: pass2
block2:
  block21:
    pass3: pass3
  pass4: pass4
noSecret1: foo,
noSecret2: bar,
pass5: pass5
pass6: 1234565
`
		expectedYAML = `block1:
  block11:
    pass1: pass1-modified
  pass2: pass2-modified
block2:
  block21:
    pass3: pass3-modified
  pass4: pass4-modified
noSecret1: foo,
noSecret2: bar,
pass5: pass5-modified
pass6: 1234565-modified
`
	}

	{
		output, err := newService.TraverseYAML([]byte(testYAML))
		if err != nil {
			t.Fatal("expected", nil, "got", err)
		}
		if string(output) != expectedYAML {
			t.Fatal("expected", fmt.Sprintf("%q", expectedYAML), "got", fmt.Sprintf("%q", output))
		}
	}
}
