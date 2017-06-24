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
	testJSON := `{
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
  "pass5": "pass5",
  "pass6": 123456
}`

	valueModifiers := []ValueModifier{
		testModifier{},
	}

	output, err := TraverseJSON([]byte(testJSON), valueModifiers...)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expectedJSON := `{
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
  "pass5": "pass5-modified",
  "pass6": "123456-modified"
}`

	if string(output) != expectedJSON {
		t.Fatal("expected", fmt.Sprintf("%q", expectedJSON), "got", fmt.Sprintf("%q", output))
	}
}

func Test_ValueModifier_TraverseYAML(t *testing.T) {
	testYAML := `block1:
  block11:
    pass1: pass1
  pass2: pass2
block2:
  block21:
    pass3: pass3
  pass4: pass4
pass5: pass5
pass6: 1234565
`

	valueModifiers := []ValueModifier{
		testModifier{},
	}

	output, err := TraverseYAML([]byte(testYAML), valueModifiers...)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expectedYAML := `block1:
  block11:
    pass1: pass1-modified
  pass2: pass2-modified
block2:
  block21:
    pass3: pass3-modified
  pass4: pass4-modified
pass5: pass5-modified
pass6: 1234565-modified
`

	if string(output) != expectedYAML {
		t.Fatal("expected", fmt.Sprintf("%q", expectedYAML), "got", fmt.Sprintf("%q", output))
	}
}
