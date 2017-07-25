package path

import (
	"reflect"
	"strconv"
	"testing"
)

func Test_Service_All(t *testing.T) {
	testCases := []struct {
		InputBytes []byte
		Expected   []string
	}{
		// Test case 1, ensure a single unnested path can be found.
		{
			InputBytes: []byte(`{
  "k1": "v1"
}`),
			Expected: []string{
				"k1",
			},
		},

		// Test case 2, ensure multiple unnested paths can be found.
		{
			InputBytes: []byte(`{
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
			InputBytes: []byte(`{
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
			InputBytes: []byte(`{
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
			InputBytes: []byte(`{
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
			InputBytes: []byte(`{
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

		// Test case 7, ensure single paths with unnested inline JSON objects can be
		// found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": "v2"
  }`) + `
}`),
			Expected: []string{
				"k1.k2",
			},
		},

		// Test case 8, ensure single paths with nested inline JSON objects can be
		// found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": {
      "k3": "v3"
    }
  }`) + `
}`),
			Expected: []string{
				"k1.k2.k3",
			},
		},

		// Test case 9, ensure single paths with unnested inline JSON lists can be
		// found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`[
    {
      "k2": "v2"
    }
  ]`) + `
}`),
			Expected: []string{
				"k1.[0].k2",
			},
		},

		// Test case 10, ensure multiple paths with unnested inline JSON objects can
		// be found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": "v2",
    "k3": "v3"
  }`) + `
}`),
			Expected: []string{
				"k1.k2",
				"k1.k3",
			},
		},

		// Test case 11, ensure multiple paths with nested inline JSON objects can
		// be found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": {
      "k3": "v3"
    },
    "k4": {
      "k5": "v5"
    }
  }`) + `
}`),
			Expected: []string{
				"k1.k2.k3",
				"k1.k4.k5",
			},
		},

		// Test case 12, ensure multiple paths with unnested inline JSON lists can
		// be found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`[
    {
      "k2": "v2"
    },
    {
      "k3": "v3"
    }
  ]`) + `
}`),
			Expected: []string{
				"k1.[0].k2",
				"k1.[1].k3",
			},
		},

		// Test case 13, ensure single paths with unnested inline YAML objects can
		// be found.
		{
			InputBytes: []byte(`{
  "k1": "k2: v2"
}`),
			Expected: []string{
				"k1.k2",
			},
		},

		// Test case 14, ensure single paths with nested inline YAML objects can be
		// found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: v3`) + `
}`),
			Expected: []string{
				"k1.k2.k3",
			},
		},

		// Test case 15, ensure single paths with unnested inline YAML lists can be
		// found.
		{
			InputBytes: []byte(`{
  "k1": "- k2: v2"
}`),
			Expected: []string{
				"k1.[0].k2",
			},
		},

		// Test case 16, ensure multiple paths with unnested inline YAML objects can
		// be found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2: v2
k3: v3`) + `
}`),
			Expected: []string{
				"k1.k2",
				"k1.k3",
			},
		},

		// Test case 17, ensure multiple paths with nested inline YAML objects can
		// be found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: v3
k4:
  k5: v5`) + `
}`),
			Expected: []string{
				"k1.k2.k3",
				"k1.k4.k5",
			},
		},

		// Test case 18, ensure multiple paths with unnested inline YAML lists can
		// be found.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`- k2: v2
- k3: v3`) + `
}`),
			Expected: []string{
				"k1.[0].k2",
				"k1.[1].k3",
			},
		},

		// Test case 19, ensure paths with separators inside keys can be found.
		{
			InputBytes: []byte(`{
  "k1": {
    "k2.k3": {
      "k4": "v4"
    }
  }
}`),
			Expected: []string{
				`k1.k2\.k3.k4`,
			},
		},
	}

	for i, testCase := range testCases {
		config := DefaultConfig()
		config.InputBytes = testCase.InputBytes
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
		InputBytes []byte
		Path       string
		Expected   interface{}
	}{
		// Test case 1, ensure the value of an unnested path can be returned.
		{
			InputBytes: []byte(`{
  "k1": "v1"
}`),
			Path:     "k1",
			Expected: "v1",
		},

		// Test case 2, ensure the value of a nested path can be returned.
		{
			InputBytes: []byte(`{
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
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:     "k1.[0].k2",
			Expected: "v2",
		},

		// Test case 4, ensure the value of single paths with unnested inline JSON
		// objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": "v2"
  }`) + `
}`),
			Path:     "k1.k2",
			Expected: "v2",
		},

		// Test case 5, ensure the value of single paths with nested inline JSON
		// objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": {
      "k3": "v3"
    }
  }`) + `
}`),
			Path:     "k1.k2.k3",
			Expected: "v3",
		},

		// Test case 6, ensure the value of single paths with unnested inline JSON
		// lists can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`[
    {
      "k2": "v2"
    }
  ]`) + `
}`),
			Path:     "k1.[0].k2",
			Expected: "v2",
		},

		// Test case 7, ensure the value of multiple paths with unnested inline JSON
		// objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": "v2",
    "k3": "v3"
  }`) + `
}`),
			Path:     "k1.k2",
			Expected: "v2",
		},

		// Test case 8, ensure the value of multiple paths with nested inline JSON
		// objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": {
      "k3": "v3"
    },
    "k4": {
      "k5": "v5"
    }
  }`) + `
}`),
			Path:     "k1.k2.k3",
			Expected: "v3",
		},

		// Test case 9, ensure the value of multiple paths with unnested inline JSON
		// lists can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`[
    {
      "k2": "v2"
    },
    {
      "k3": "v3"
    }
  ]`) + `
}`),
			Path:     "k1.[0].k2",
			Expected: "v2",
		},

		// Test case 10, ensure the value of single paths with unnested inline YAML
		// objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": "k2: v2"
}`),
			Path:     "k1.k2",
			Expected: "v2",
		},

		// Test case 11, ensure the value of single paths with nested inline YAML
		// objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: v3`) + `
}`),
			Path:     "k1.k2.k3",
			Expected: "v3",
		},

		// Test case 12, ensure the value of single paths with unnested inline YAML
		// lists can be returned.
		{
			InputBytes: []byte(`{
  "k1": "- k2: v2"
}`),
			Path:     "k1.[0].k2",
			Expected: "v2",
		},

		// Test case 13, ensure the value of multiple paths with unnested inline
		// YAML objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2: v2
k3: v3`) + `
}`),
			Path:     "k1.k2",
			Expected: "v2",
		},

		// Test case 14, ensure the value of multiple paths with nested inline YAML
		// objects can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: v3
k4:
  k5: v5`) + `
}`),
			Path:     "k1.k2.k3",
			Expected: "v3",
		},

		// Test case 15, ensure the value of multiple paths with unnested inline
		// YAML lists can be returned.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`- k2: v2
- k3: v3`) + `
}`),
			Path:     "k1.[0].k2",
			Expected: "v2",
		},

		// Test case 16, ensure the value of paths with separators inside keys can
		// be returned.
		{
			InputBytes: []byte(`{
  "k1": {
    "k2.k3": {
      "k4": "v4"
    }
  }
}`),
			Path:     `k1.k2\.k3.k4`,
			Expected: "v4",
		},
	}

	for i, testCase := range testCases {
		config := DefaultConfig()
		config.InputBytes = testCase.InputBytes
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

func Test_Service_Get_Error(t *testing.T) {
	textCases := []struct {
		InputBytes   []byte
		Path         string
		ErrorMatcher func(error) bool
	}{
		// Test 1, when there is only 1 element in the list index [1] cannot be
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:         "k1.[1].k2",
			ErrorMatcher: IsNotFound,
		},

		// Test 2, when there is k1 at the beginning of the path key k3 cannot be
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:         "k3.[0].k2",
			ErrorMatcher: IsNotFound,
		},

		// Test 3, when there is k2 at the end of the path key k3 cannot be
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:         "k1.[0].k3",
			ErrorMatcher: IsNotFound,
		},
	}

	for i, tc := range textCases {
		config := DefaultConfig()
		config.InputBytes = tc.InputBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		_, err = newService.Get(tc.Path)
		if !tc.ErrorMatcher(err) {
			t.Fatal("test", i+1, "expected", true, "got", false)
		}
	}
}

func Test_Service_Set(t *testing.T) {
	textCases := []struct {
		InputBytes []byte
		Path       string
		Value      string
		Expected   []byte
	}{
		// Test 1,
		{
			InputBytes: []byte(`{
  "k1": "v1"
}`),
			Path:  "k1",
			Value: "modified",
			Expected: []byte(`{
  "k1": "modified"
}`),
		},

		// Test 2,
		{
			InputBytes: []byte(`{
  "k1": "v1",
  "k2": "v2"
}`),
			Path:  "k1",
			Value: "modified",
			Expected: []byte(`{
  "k1": "modified",
  "k2": "v2"
}`),
		},

		// Test 3,
		{
			InputBytes: []byte(`{
  "k1": {
    "k2": "v2"
  }
}`),
			Path:  "k1.k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": {
    "k2": "modified"
  }
}`),
		},

		// Test 4,
		{
			InputBytes: []byte(`{
  "k1": {
    "k2": "v2"
  },
  "k3": "v3"
}`),
			Path:  "k1.k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": {
    "k2": "modified"
  },
  "k3": "v3"
}`),
		},

		// Test 5,
		{
			InputBytes: []byte(`{
  "k1": {
    "k2": {
      "k3": "v3"
    }
  }
}`),
			Path:  "k1.k2.k3",
			Value: "modified",
			Expected: []byte(`{
  "k1": {
    "k2": {
      "k3": "modified"
    }
  }
}`),
		},

		// Test 6,
		{
			InputBytes: []byte(`[
  {
    "k1": "v1"
  }
]`),
			Path:  "[0].k1",
			Value: "modified",
			Expected: []byte(`[
  {
    "k1": "modified"
  }
]`),
		},

		// Test 7,
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:  "k1.[0].k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": [
    {
      "k2": "modified"
    }
  ]
}`),
		},

		// Test 8,
		{
			InputBytes: []byte(`{
  "k1": [
    [
      {
        "k2": "v2"
      }
    ],
    [
      {
        "k3": "v3"
      }
    ]
  ]
}`),
			Path:  "k1.[1].[0].k3",
			Value: "modified",
			Expected: []byte(`{
  "k1": [
    [
      {
        "k2": "v2"
      }
    ],
    [
      {
        "k3": "modified"
      }
    ]
  ]
}`),
		},

		// Test case 9, ensure the value of single paths with unnested inline JSON
		// objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": "v2"
  }`) + `
}`),
			Path:  "k1.k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`{
  "k2": "modified"
}`) + `
}`),
		},

		// Test case 10, ensure the value of single paths with nested inline JSON
		// objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": {
      "k3": "v3"
    }
  }`) + `
}`),
			Path:  "k1.k2.k3",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`{
  "k2": {
    "k3": "modified"
  }
}`) + `
}`),
		},

		// Test case 11, ensure the value of single paths with unnested inline JSON
		// lists can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`[
    {
      "k2": "v2"
    }
  ]`) + `
}`),
			Path:  "k1.[0].k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`[
  {
    "k2": "modified"
  }
]`) + `
}`),
		},

		// Test case 12, ensure the value of multiple paths with unnested inline
		// JSON objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": "v2",
    "k3": "v3"
  }`) + `
}`),
			Path:  "k1.k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`{
  "k2": "modified",
  "k3": "v3"
}`) + `
}`),
		},

		// Test case 13, ensure the value of multiple paths with nested inline JSON
		// objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`{
    "k2": {
      "k3": "v3"
    },
    "k4": {
      "k5": "v5"
    }
  }`) + `
}`),
			Path:  "k1.k2.k3",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`{
  "k2": {
    "k3": "modified"
  },
  "k4": {
    "k5": "v5"
  }
}`) + `
}`),
		},

		// Test case 14, ensure the value of multiple paths with unnested inline
		// JSON lists can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`[
    {
      "k2": "v2"
    },
    {
      "k3": "v3"
    }
  ]`) + `
}`),
			Path:  "k1.[0].k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`[
  {
    "k2": "modified"
  },
  {
    "k3": "v3"
  }
]`) + `
}`),
		},

		// Test case 15, ensure the value of single paths with unnested inline YAML
		// objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": "k2: v2"
}`),
			Path:  "k1.k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`k2: modified
`) + `
}`),
		},

		// Test case 16, ensure the value of single paths with nested inline YAML
		// objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: v3`) + `
}`),
			Path:  "k1.k2.k3",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: modified
`) + `
}`),
		},

		// Test case 17, ensure the value of single paths with unnested inline YAML
		// lists can be modified.
		{
			InputBytes: []byte(`{
  "k1": "- k2: v2"
}`),
			Path:  "k1.[0].k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`- k2: modified
`) + `
}`),
		},

		// Test case 18, ensure the value of multiple paths with unnested inline
		// YAML objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2: v2
k3: v3`) + `
}`),
			Path:  "k1.k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`k2: modified
k3: v3
`) + `
}`),
		},

		// Test case 19, ensure the value of multiple paths with nested inline YAML
		// objects can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: v3
k4:
  k5: v5`) + `
}`),
			Path:  "k1.k2.k3",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`k2:
  k3: modified
k4:
  k5: v5
`) + `
}`),
		},

		// Test case 20, ensure the value of multiple paths with unnested inline
		// YAML lists can be modified.
		{
			InputBytes: []byte(`{
  "k1": ` + strconv.Quote(`- k2: v2
- k3: v3`) + `
}`),
			Path:  "k1.[0].k2",
			Value: "modified",
			Expected: []byte(`{
  "k1": ` + strconv.Quote(`- k2: modified
- k3: v3
`) + `
}`),
		},

		// Test case 21, ensure the value of paths with separators inside keys can
		// be modified.
		{
			InputBytes: []byte(`{
  "k1": {
    "k2.k3": {
      "k4": "v4"
    }
  }
}`),
			Path:  `k1.k2\.k3.k4`,
			Value: "modified",
			Expected: []byte(`{
  "k1": {
    "k2.k3": {
      "k4": "modified"
    }
  }
}`),
		},
	}

	for i, tc := range textCases {
		config := DefaultConfig()
		config.InputBytes = tc.InputBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		err = newService.Set(tc.Path, tc.Value)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		output, err := newService.OutputBytes()
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
		if !reflect.DeepEqual(tc.Expected, output) {
			t.Fatal("test", i+1, "expected", string(tc.Expected), "got", string(output))
		}
	}
}

func Test_Service_Set_Error(t *testing.T) {
	textCases := []struct {
		InputBytes   []byte
		Path         string
		ErrorMatcher func(error) bool
	}{
		// Test 1, when there is only 1 element in the list index [1] cannot be
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:         "k1.[1].k2",
			ErrorMatcher: IsNotFound,
		},

		// Test 2, when there is k1 at the beginning of the path key k3 cannot be
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:         "k3.[0].k2",
			ErrorMatcher: IsNotFound,
		},

		// Test 3, when there is k2 at the end of the path key k3 cannot be
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Path:         "k1.[0].k3",
			ErrorMatcher: IsNotFound,
		},
	}

	for i, tc := range textCases {
		config := DefaultConfig()
		config.InputBytes = tc.InputBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		err = newService.Set(tc.Path, "value")
		if !tc.ErrorMatcher(err) {
			t.Fatal("test", i+1, "expected", true, "got", false)
		}
	}
}
