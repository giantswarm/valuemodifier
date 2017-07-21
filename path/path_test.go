package path

import (
	"reflect"
	"testing"
)

func Test_Service_All(t *testing.T) {
	testCases := []struct {
		JSONBytes []byte
		Expected  []string
	}{
		// Test case 1, ensure a single unnested path can be found.
		{
			JSONBytes: []byte(`{
  "k1": "v1"
}`),
			Expected: []string{
				"k1",
			},
		},
		// Test case 2, ensure multiple unnested paths can be found.
		{
			JSONBytes: []byte(`{
  "k1": "v1",
  "k2": "v2",
  "k3": "v3"
}`),
			Expected: []string{
				"k1",
				"k2",
				"k3",
			},
		},
		// Test case 3, ensure a single nested path can be found.
		{
			JSONBytes: []byte(`{
  "k1": {
    "k2": {
      "k3": "v3"
    }
  }
}`),
			Expected: []string{
				"k1.k2.k3",
			},
		},
		// Test case 4, ensure multiple nested paths can be found.
		{
			JSONBytes: []byte(`{
  "k1": {
    "k2": {
      "k3": "v3"
    }
  },
  "k11": {
    "k22": {
      "k33": "v33"
    }
  }
}`),
			Expected: []string{
				"k1.k2.k3",
				"k11.k22.k33",
			},
		},
		// Test case 5, multiple unnested paths and multiple nested paths can be
		// found at the same time.
		{
			JSONBytes: []byte(`{
  "k1": "v1",
  "k2": "v2",
  "k3": "v3",
  "k4": {
    "k5": {
      "k6": "v6"
    }
  },
  "k7": {
    "k8": {
      "k9": "v9"
    }
  }
}`),
			Expected: []string{
				"k1",
				"k2",
				"k3",
				"k4.k5.k6",
				"k7.k8.k9",
			},
		},
		// Test case 6, ensure paths across lists can be found.
		{
			JSONBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    },
    {
      "k3": "v3"
    }
  ]
}`),
			Expected: []string{
				"k1.[0].k2",
				"k1.[1].k3",
			},
		},
	}

	for i, testCase := range testCases {
		config := DefaultConfig()
		config.JSONBytes = testCase.JSONBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		output, err := newService.All()
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
		if !reflect.DeepEqual(testCase.Expected, output) {
			t.Fatal("test", i+1, "expected", testCase.Expected, "got", output)
		}
	}
}

func Test_Service_Get(t *testing.T) {
	testCases := []struct {
		JSONBytes []byte
		Path      string
		Expected  interface{}
	}{
		// Test case 1, ensure the value of an unnested path can be returned.
		{
			JSONBytes: []byte(`{
  "k1": "v1"
}`),
			Path:     "k1",
			Expected: "v1",
		},
		// Test case 2, ensure the value of a nested path can be returned.
		{
			JSONBytes: []byte(`{
  "k1": {
    "k2": {
      "k3": "v3"
    }
  }
}`),
			Path:     "k1.k2.k3",
			Expected: "v3",
		},
		// Test case 3, ensure the value of a list path can be returned.
		{
			JSONBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:     "k1.[0].k2",
			Expected: "v2",
		},
	}

	for i, testCase := range testCases {
		config := DefaultConfig()
		config.JSONBytes = testCase.JSONBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		output, err := newService.Get(testCase.Path)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
		if !reflect.DeepEqual(testCase.Expected, output) {
			t.Fatal("test", i+1, "expected", testCase.Expected, "got", output)
		}
	}
}

func Test_Service_Set(t *testing.T) {
	tests := []struct {
		json     []byte
		key      string
		value    string
		expected []byte
	}{
		{
			json: []byte(`{
	"k1": "v1"
}`),
			key:   "k1",
			value: "modified",
			expected: []byte(`{
	"k1": "modified"
}`),
		},

		{
			json: []byte(`{
	"k1": "v1",
	"k2": "v2"
}`),
			key:   "k1",
			value: "modified",
			expected: []byte(`{
	"k1": "modified",
	"k2": "v2"
}`),
		},

		{
			json: []byte(`{
	"k1": {
		"k2": "v2"
	}
}`),
			key:   "k1.k2",
			value: "modified",
			expected: []byte(`{
	"k1": {
		"k2": "modified"
	}
}`),
		},

		{
			json: []byte(`{
	"k1": {
		"k2": "v2"
	},
	"k3": "v3"
}`),
			key:   "k1.k2",
			value: "modified",
			expected: []byte(`{
	"k1": {
		"k2": "modified"
	},
	"k3": "v3"
}`),
		},

		{
			json: []byte(`{
	"k1": {
		"k2": {
			"k3": "v3"
		}
	}
}`),
			key:   "k1.k2.k3",
			value: "modified",
			expected: []byte(`{
	"k1": {
		"k2": {
			"k3": "modified"
		}
	}
}`),
		},
	}

	for index, test := range tests {
		config := DefaultConfig()
		config.JSONBytes = test.json
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", index, "expected", nil, "got", err)
		}

		if err := newService.Set(test.key, test.value); err != nil {
			t.Fatal("test", index, "expected", nil, "got", err)
		}

		output := newService.JSONBytes()
		if !reflect.DeepEqual(test.expected, output) {
			t.Fatal("test", index, "expected", string(test.expected), "got", string(output))
		}
	}
}
