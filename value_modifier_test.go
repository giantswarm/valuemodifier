package valuemodifier

import (
	"fmt"
	"strconv"
	"testing"
)

type testModifier1 struct{}

func (m testModifier1) Modify(value []byte) ([]byte, error) {
	return []byte(string(value) + "-modified1"), nil
}

type testModifier2 struct{}

func (m testModifier2) Modify(value []byte) ([]byte, error) {
	return []byte(string(value) + "-modified2"), nil
}

func Test_ValueModifier_Traverse_JSON(t *testing.T) {
	testCases := []struct {
		ValueModifiers []ValueModifier
		IgnoreFields   []string
		Input          string
		Expected       string
	}{
		// Test case 0, a single modifier modifies all secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `{
  "noSecret1": "noSecret1",
  "pass1": "pass1"
}`,
			Expected: `{
  "noSecret1": "noSecret1-modified1",
  "pass1": "pass1-modified1"
}`,
		},
		// Test case 1, a single modifier modifies all numeric secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `{
  "noSecret1": "noSecret1",
  "pass1": 12345
}`,
			Expected: `{
  "noSecret1": "noSecret1-modified1",
  "pass1": "12345-modified1"
}`,
		},
		// Test case 2, a single modifier modifies all secrets inside lists.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			Input: `{
  "list1": [
    {
      "pass1": "pass1"
    }
  ]
}`,
			Expected: `{
  "list1": [
    {
      "pass1": "pass1-modified1"
    }
  ]
}`,
		},
		// Test case 3, a single modifier modifies all secrets, but ignores the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{
				"noSecret1",
			},
			Input: `{
  "noSecret1": "foo",
  "pass1": "pass1"
}`,
			Expected: `{
  "noSecret1": "foo",
  "pass1": "pass1-modified1"
}`,
		},
		// Test case 4, multiple modifiers modify all secrets, but ignore the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			Input: `{
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass1": "pass1",
  "pass2": "pass2"
}`,
			Expected: `{
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass1": "pass1-modified1-modified2",
  "pass2": "pass2-modified1-modified2"
}`,
		},
		// Test case 5, nested blocks, multiple modifiers modify all secrets, but
		// ignore the ones configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			Input: `{
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
}`,
			Expected: `{
  "block1": {
    "block11": {
      "pass1": "pass1-modified1-modified2"
    },
    "pass2": "pass2-modified1-modified2"
  },
  "block2": {
    "block21": {
      "pass3": "pass3-modified1-modified2"
    },
    "pass4": "pass4-modified1-modified2"
  },
  "noSecret1": "foo",
  "noSecret2": "bar",
  "pass5": "pass5-modified1-modified2",
  "pass6": "123456-modified1-modified2"
}`,
		},
	}

	for i, testCase := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			config := DefaultConfig()
			config.ValueModifiers = testCase.ValueModifiers
			config.IgnoreFields = testCase.IgnoreFields
			newService, err := New(config)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			output, err := newService.Traverse([]byte(testCase.Input))
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
			if string(output) != testCase.Expected {
				t.Fatal("expected", fmt.Sprintf("%q", testCase.Expected), "got", fmt.Sprintf("%q", output))
			}
		})
	}
}

func Test_ValueModifier_Traverse_YAML(t *testing.T) {
	testCases := []struct {
		ValueModifiers []ValueModifier
		IgnoreFields   []string
		SelectFields   []string
		Input          string
		Expected       string
	}{
		// Test case 0, a single modifier modifies all secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `noSecret1: noSecret1
pass1: pass1
`,
			Expected: `noSecret1: noSecret1-modified1
pass1: pass1-modified1
`,
		},

		// Test case 1, a single modifier modifies all numeric secrets.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `noSecret1: noSecret1
pass1: 12345
`,
			Expected: `noSecret1: noSecret1-modified1
pass1: 12345-modified1
`,
		},

		// Test case 2, a single modifier modifies all secrets inside lists.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `list1:
- pass1: pass1
`,
			Expected: `list1:
- pass1: pass1-modified1
`,
		},

		// Test case 3, a single modifier modifies all secrets, but ignores the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{
				"noSecret1",
			},
			SelectFields: []string{},
			Input: `noSecret1: noSecret1
pass1: pass1
`,
			Expected: `noSecret1: noSecret1
pass1: pass1-modified1
`,
		},

		// Test case 4, multiple modifiers modify all secrets, but ignore the ones
		// configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			SelectFields: []string{},
			Input: `noSecret1: noSecret1
noSecret2: noSecret2
pass1: pass1
pass2: pass2
`,
			Expected: `noSecret1: noSecret1
noSecret2: noSecret2
pass1: pass1-modified1-modified2
pass2: pass2-modified1-modified2
`,
		},

		// Test case 5, nested blocks, multiple modifiers modify all secrets, but
		// ignore the ones configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
				testModifier2{},
			},
			IgnoreFields: []string{
				"noSecret1",
				"noSecret2",
			},
			SelectFields: []string{},
			Input: `block1:
  block11:
    pass1: pass1
  pass2: pass2
block2:
  block21:
    pass3: pass3
  pass4: pass4
noSecret1: foo
noSecret2: bar
pass5: pass5
pass6: 123456
`,
			Expected: `block1:
  block11:
    pass1: pass1-modified1-modified2
  pass2: pass2-modified1-modified2
block2:
  block21:
    pass3: pass3-modified1-modified2
  pass4: pass4-modified1-modified2
noSecret1: foo
noSecret2: bar
pass5: pass5-modified1-modified2
pass6: 123456-modified1-modified2
`,
		},

		// Test case 6, modifiers modify string blocks.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `pass1: |
  foo
  bar
`,
			Expected: `pass1: |-
  foo
  bar
  -modified1
`,
		},

		// Test case 7, modifiers modify secrets of string blocks representing YAML.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `pass1: |
  bar:
    baz: pass2
  foo: pass3
`,
			Expected: `pass1: |-
  bar:
    baz: pass2
  foo: pass3
  -modified1
`,
		},

		// Test case 8, modifiers modify secrets of string blocks representing JSON.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `pass1: |
  {
    "block1": {
      "block11": {
        "pass2": "pass2"
      }
    }
  }
`,
			Expected: `pass1: |-
  {
    "block1": {
      "block11": {
        "pass2": "pass2"
      }
    }
  }
  -modified1
`,
		},

		// Test case 9, a single modifier modifies all secrets, but ignores the
		// ones configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{
				"block1.block11.pass1",
				"pass1",
			},
			Input: `block1:
  block11:
    pass1: pass1
  pass2: pass2
pass1: pass1
pass2: pass2
`,
			Expected: `block1:
  block11:
    pass1: pass1-modified1
  pass2: pass2
pass1: pass1-modified1
pass2: pass2
`,
		},

		// Test case 10, ensure a real config example for modifier modifies all secrets,
		// but ignores the ones configured using IgnoreFields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{
				"ipmiIp",
				"privateIp",
				"providerId",
			},
			SelectFields: []string{},
			Input: `credentials:
- ipmiIp: 192.168.0.2
  password: password
  privateIp: 10.0.0.1
  providerId: providerId
  username: username
`,
			Expected: `credentials:
- ipmiIp: 192.168.0.2
  password: password-modified1
  privateIp: 10.0.0.1
  providerId: providerId
  username: username-modified1
`,
		},

		// Test case 11, ensure a real world example works with all fields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `Installation:
  V1:
    Secret:
      Alertmanager:
        Nginx:
          Auth: magic
      Prometheus:
        Nginx:
          Auth: magic
      Registry:
        PullSecret:
          DockerConfigJSON: |-
            {
              "auths": {
                "quay.io": {
                  "auth": "magic"
                }
              }
            }
`,
			Expected: `Installation:
  V1:
    Secret:
      Alertmanager:
        Nginx:
          Auth: magic-modified1
      Prometheus:
        Nginx:
          Auth: magic-modified1
      Registry:
        PullSecret:
          DockerConfigJSON: |-
            {
              "auths": {
                "quay.io": {
                  "auth": "magic"
                }
              }
            }-modified1
`,
		},

		// Test case 11, ensure a real world example works with selected fields.
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{
				"Installation.V1.Secret.Alertmanager.Nginx.Auth",
			},
			Input: `Installation:
  V1:
    Secret:
      Alertmanager:
        Nginx:
          Auth: magic
      Prometheus:
        Nginx:
          Auth: magic
      Registry:
        PullSecret:
          DockerConfigJSON: |-
            {
              "auths": {
                "quay.io": {
                  "auth": "magic"
                }
              }
            }
`,
			Expected: `Installation:
  V1:
    Secret:
      Alertmanager:
        Nginx:
          Auth: magic-modified1
      Prometheus:
        Nginx:
          Auth: magic
      Registry:
        PullSecret:
          DockerConfigJSON: |-
            {
              "auths": {
                "quay.io": {
                  "auth": "magic"
                }
              }
            }
`,
		},
		// Test case 12, modifier operating on a slice
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `k1:
- k2
- k3`,
			Expected: `k1:
- k2-modified1
- k3-modified1
`,
		},
		// Test case 13, modifier operating on a slice, mixed values
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `k1:
- k2
- k3
- 8080`,
			Expected: `k1:
- k2-modified1
- k3-modified1
- 8080-modified1
`,
		},
		// Test case 14, modifier operating on a slice, null
		{
			ValueModifiers: []ValueModifier{
				testModifier1{},
			},
			IgnoreFields: []string{},
			SelectFields: []string{},
			Input: `k1:
- k2
- k3
- null
- 8080`,
			Expected: `k1:
- k2-modified1
- k3-modified1
- null
- 8080-modified1
`,
		},
	}

	for i, testCase := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			config := DefaultConfig()
			config.ValueModifiers = testCase.ValueModifiers
			config.IgnoreFields = testCase.IgnoreFields
			config.SelectFields = testCase.SelectFields
			newService, err := New(config)
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}

			output, err := newService.Traverse([]byte(testCase.Input))
			if err != nil {
				t.Fatal("expected", nil, "got", err)
			}
			if string(output) != testCase.Expected {
				t.Fatal("expected", fmt.Sprintf("%q", testCase.Expected), "got", fmt.Sprintf("%q", output))
			}
		})
	}
}
