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
				"k1",
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
				"k1",
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
				"k1",
			},
		},

		// Test case 10, ensure multiple paths with unnested inline JSON lists can
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
				"k1",
			},
		},

		// Test case 11, ensure paths with separators inside keys can be found.
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

		// Test case 12, ensure a single unnested path can be found, even though its
		// value is empty.
		{
			InputBytes: []byte(`{
  "k1": ""
}`),
			Expected: []string{
				"k1",
			},
		},

		// Test case 13, ensure a single nested path can be found, even though its
		// value is empty.
		{
			InputBytes: []byte(`{
  "k1": {
    "k2": {
      "k3": ""
    }
  }
}`),
			Expected: []string{
				"k1.k2.k3",
			},
		},

		// Test case 14, ensure empty fields are handled (YAML standard allows for such fields).
		{
			InputBytes: []byte(`
k1:
  k2:
k3: "v"
k4: ["a", null "b"]
k5: null
`),
			Expected: []string{"k1.k2", "k3", "k4", "k5"},
		},

		// Test case 15, ensure array of objects is handled correctly.
		{
			InputBytes: []byte(`
tolerations:
- effect: NoSchedule
`),
			Expected: []string{"tolerations.[0].effect"},
		},
	}

	for i, tc := range testCases {
		config := DefaultConfig()
		config.InputBytes = tc.InputBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		output, err := newService.All()
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
		if !reflect.DeepEqual(tc.Expected, output) {
			t.Fatal("test", i+1, "expected", tc.Expected, "got", output)
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

		// Test case 17, ensure the value of a single unnested path can be returned,
		// even though its value is empty.
		{
			InputBytes: []byte(`{
  "k1": ""
}`),
			Path:     "k1",
			Expected: "",
		},

		// Test case 18, ensure the value of a single nested path can be returned,
		// even though its value is empty.
		{
			InputBytes: []byte(`{
  "k1": {
    "k2": {
      "k3": ""
    }
  }
}`),
			Path:     "k1.k2.k3",
			Expected: "",
		},

		// Test case 19, ensure the value of multiple paths(array index > 9) with unnested inline
		// YAML lists can be returned.
		{
			InputBytes: []byte(`{
  "k": ` + strconv.Quote(`- k0: v0
- k1: v1
- k2: v2
- k3: v3
- k4: v4
- k5: v5
- k6: v6
- k7: v7
- k8: v8
- k9: v9
- k10: v10`) + `
}`),
			Path:     "k.[10].k10",
			Expected: "v10",
		},

		// Test case 20, ensure nil values are handled properly - explicit null
		{
			InputBytes: []byte(`
k1:
  k2: null
`),
			Path:     "k1.k2",
			Expected: nil,
		},

		// Test case 21, ensure nil values are handled properly - empty value
		{
			InputBytes: []byte(`k1:`),
			Path:       "k1",
			Expected:   nil,
		},

		// Test case 22, ensure nil values are handled properly - null value in an array
		{
			InputBytes: []byte(`k1: [null]`),
			Path:       "k1",
			Expected:   []interface{}{nil},
		},
	}

	for i, tc := range testCases {
		config := DefaultConfig()
		config.InputBytes = tc.InputBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		output, err := newService.Get(tc.Path)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
		if !reflect.DeepEqual(tc.Expected, output) {
			t.Fatal("test", i+1, "expected", tc.Expected, "got", output)
		}
	}
}

func Test_Service_Get_Error(t *testing.T) {
	testCases := []struct {
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

	for i, tc := range testCases {
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
	testCases := []struct {
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

		// Test case 22, ensure the value of a single unnested path can be modified,
		// even though its value is empty.
		{
			InputBytes: []byte(`{
  "k1": ""
}`),
			Path:  "k1",
			Value: "modified",
			Expected: []byte(`{
  "k1": "modified"
}`),
		},

		// Test case 23, ensure the value of a single nested path can be modified,
		// even though its value is empty.
		{
			InputBytes: []byte(`{
  "k1": {
    "k2": {
      "k3": ""
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

		// Test case 24, ensure the value of a single unnested path can be set,
		// even though its key is missing.
		{
			InputBytes: []byte(`{
  "k1": ""
}`),
			Path:  "k2",
			Value: "added",
			Expected: []byte(`{
  "k1": "",
  "k2": "added"
}`),
		},

		// Test case 25, ensure the value of a multi level nested path can be set,
		// even though its key is missing.
		{
			InputBytes: []byte(`{
  "k1": {}
}`),
			Path:  "k1.s1.e3",
			Value: "added",
			Expected: []byte(`{
  "k1": {
    "s1": {
      "e3": "added"
    }
  }
}`),
		},

		// Test case 26, ensure the value of a multi level nested path can be set,
		// even though its key is missing.
		{
			InputBytes: []byte(`{
  "k1": ""
}`),
			Path:  "k2.s1.e3",
			Value: "added",
			Expected: []byte(`{
  "k1": "",
  "k2": {
    "s1": {
      "e3": "added"
    }
  }
}`),
		},

		// Test case 27, ensure array of objects is handled correctly.
		{
			InputBytes: []byte(`{
"tolerations": [{}],
}`),
			Path:  "tolerations.[0].effect",
			Value: "NoSchedule",
			Expected: []byte(`tolerations:
- effect: NoSchedule
`),
		},

		// Test case 28, ensure array of objects is handled correctly.
		{
			InputBytes: []byte(`{
"tolerations": [],
}`),
			Path:  "tolerations.[0].effect",
			Value: "NoSchedule",
			Expected: []byte(`tolerations:
- effect: NoSchedule
`),
		},
	}

	for i, tc := range testCases {
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
	testCases := []struct {
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
	}

	for i, tc := range testCases {
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

func Test_Service_Validate(t *testing.T) {
	testCases := []struct {
		InputBytes []byte
		Paths      []string
	}{
		// Test 1, when there is only 1 field to valiate with one level nesting
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": "v2"
    }
  ]
}`),
			Paths: []string{"k2"},
		},

		// Test 2, when there is only 1 field to valiate with two level nesting
		// found.
		{
			InputBytes: []byte(`{
  "k1": [
    {
      "k2": {
				  "k3": "v3"
			}
    }
  ]
}`),
			Paths: []string{"k3"},
		},

		// Test 3, when there are only 2 fields to valiate with different level nesting
		// found.
		{
			InputBytes: []byte(`{
	"k1": [
		{
			"k2": "v2",
			"k3": {
				"k4": "v4"
		  }
		}
	]
}`),
			Paths: []string{"k2", "k4"},
		},
	}

	for i, tc := range testCases {
		config := DefaultConfig()
		config.InputBytes = tc.InputBytes
		newService, err := New(config)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}

		err = newService.Validate(tc.Paths)
		if err != nil {
			t.Fatal("test", i+1, "expected", nil, "got", err)
		}
	}
}
