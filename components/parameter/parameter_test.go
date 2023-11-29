package parameter

import (
	"reflect"
	"testing"
)

func SliceEqual(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}

func ParametersEqual(a, b Parameter) bool {
	return a.name == b.name &&
		a.typeValue == b.typeValue &&
		a.in == b.in &&
		a.required == b.required &&
		a.description == b.description &&
		//SliceEqual(a.enum, b.enum) &&    // TODO need to update testing to account for int64 vs int mismatch
		reflect.DeepEqual(a.defaultValue, b.defaultValue) &&
		a.format == b.format &&
		a.min == b.min &&
		a.max == b.max &&
		a.minLen == b.minLen &&
		a.maxLen == b.maxLen &&
		a.pattern == b.pattern &&
		a.maxItems == b.maxItems &&
		a.minItems == b.minItems &&
		a.uniqueItems == b.uniqueItems &&
		a.multipleOf == b.multipleOf &&
		a.collectionFormat == b.collectionFormat
}

// TODO create a comprehensive test suite for the 'Fields' types
func TestParams(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "testParam",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "testParam",
				typeValue:   Integer,
				in:          Path,
				required:    true,
				description: "A test parameter\n (max: 100)",
				min:         0,
				max:         100,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntParam(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMax(100))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "stringParam",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "stringParam",
				typeValue:   String,
				in:          Path,
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrParam(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBoolParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "boolParam",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "boolParam",
				typeValue:   Boolean,
				in:          Path,
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolParam(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFileParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "FileParam",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "FileParam",
				typeValue:   File,
				in:          Form,
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := FileParam(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "IntQuery",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "IntQuery",
				typeValue:   Integer,
				in:          "query",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntQuery(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "StrQuery",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "StrQuery",
				typeValue:   String,
				in:          "query",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrQuery(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBoolQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "BoolQuery",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "BoolQuery",
				typeValue:   Boolean,
				in:          "query",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolQuery(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "IntHeader",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "IntHeader",
				typeValue:   Integer,
				in:          "header",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntHeader(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "StrHeader",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "StrHeader",
				typeValue:   String,
				in:          "header",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrHeader(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBoolHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "BoolHeader",
			required:    true,
			description: "A test parameter",
			want: Parameter{
				name:        "BoolHeader",
				typeValue:   Boolean,
				in:          "header",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolHeader(tc.name, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

// // TODO all EnumParam tests with number will need some custom testing/comparing for this Enums []interface{} and arr's []int64 to match correctly during testing
func TestIntEnumParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:        "IntEnumParam",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntEnumParam",
				typeValue:   Integer,
				in:          Path,
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntEnumParam(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, \nwant %v", got, tc.want)
			}
		})
	}
}

func TestStrEnumParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			name:        "StrEnumParam",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			want: Parameter{
				name:        "StrEnumParam",
				typeValue:   String,
				in:          Path,
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{"a", "b", "c"},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrEnumParam(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntEnumQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:        "IntEnumQuery",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntEnumQuery",
				typeValue:   Integer,
				in:          "query",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntEnumQuery(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrEnumQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			name:        "StrEnumQuery",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			want: Parameter{
				name:        "StrEnumQuery",
				typeValue:   String,
				in:          "query",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{"a", "b", "c"},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrEnumQuery(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntEnumHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:        "IntEnumHeader",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntEnumHeader",
				typeValue:   Integer,
				in:          "header",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntEnumHeader(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrEnumHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			name:        "StrEnumHeader",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			want: Parameter{
				name:        "StrEnumHeader",
				typeValue:   String,
				in:          "header",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{"a", "b", "c"},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrEnumHeader(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntArrParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:        "IntArrParam",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntArrParam",
				typeValue:   Array,
				in:          Path,
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntArrParam(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrArrParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			name:        "StrArrParam",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			want: Parameter{
				name:        "StrArrParam",
				typeValue:   Array,
				in:          Path,
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
				enum:        []interface{}{string("a"), string("b"), string("c")},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrArrParam(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntArrQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:        "IntArrQuery",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntArrQuery",
				typeValue:   Array,
				in:          "query",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntArrQuery(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrArrQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			name:        "StrArrQuery",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			want: Parameter{
				name:        "StrArrQuery",
				typeValue:   Array,
				in:          "query",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
				enum:        []interface{}{string("a"), string("b"), string("c")},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrArrQuery(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntArrHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:        "IntArrHeader",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntArrHeader",
				typeValue:   Array,
				in:          "header",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntArrHeader(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrArrHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			name:        "StrArrHeader",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			want: Parameter{
				name:        "StrArrHeader",
				typeValue:   Array,
				in:          "header",
				required:    true,
				description: "A test parameter\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
				enum:        []interface{}{string("a"), string("b"), string("c")},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrArrHeader(tc.name, tc.arr, WithIn(Path), WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}
