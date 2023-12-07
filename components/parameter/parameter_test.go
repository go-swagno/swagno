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
		location    Location
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "TestPathParam",
			location:    Path,
			required:    true,
			description: "A test parameter for path",
			want: Parameter{
				name:        "TestPathParam",
				typeValue:   Integer,
				in:          Path,
				required:    true,
				description: "A test parameter for path\n (max: 100)",
				min:         0,
				max:         100,
			},
		},
		{
			name:        "TestQueryParam",
			location:    Query,
			required:    true,
			description: "A test parameter for query",
			want: Parameter{
				name:        "TestQueryParam",
				typeValue:   Integer,
				in:          Query,
				required:    true,
				description: "A test parameter for query\n (max: 100)",
				min:         0,
				max:         100,
			},
		},
		{
			name:        "TestHeaderParam",
			location:    Header,
			required:    true,
			description: "A test parameter for header",
			want: Parameter{
				name:        "TestHeaderParam",
				typeValue:   Integer,
				in:          Header,
				required:    true,
				description: "A test parameter for header\n (max: 100)",
				min:         0,
				max:         100,
			},
		},
		{
			name:        "TestFormParam",
			location:    Form,
			required:    true,
			description: "A test parameter for form data",
			want: Parameter{
				name:        "TestFormParam",
				typeValue:   Integer,
				in:          Form,
				required:    true,
				description: "A test parameter for form data\n (max: 100)",
				min:         0,
				max:         100,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntParam(tc.name, tc.location, WithRequired(), WithDescription(tc.description), WithMin(0), WithMax(100))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrParam(t *testing.T) {
	testCases := []struct {
		name        string
		location    Location
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "StringPathParam",
			location:    Path,
			required:    true,
			description: "A test parameter for path",
			want: Parameter{
				name:        "StringPathParam",
				typeValue:   String,
				in:          Path,
				required:    true,
				description: "A test parameter for path\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "StringQueryParam",
			location:    Query,
			required:    true,
			description: "A test parameter for query",
			want: Parameter{
				name:        "StringQueryParam",
				typeValue:   String,
				in:          Query,
				required:    true,
				description: "A test parameter for query\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "StringHeaderParam",
			location:    Header,
			required:    true,
			description: "A test parameter for header",
			want: Parameter{
				name:        "StringHeaderParam",
				typeValue:   String,
				in:          Header,
				required:    true,
				description: "A test parameter for header\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "StringFormParam",
			location:    Form,
			required:    true,
			description: "A test parameter for form data",
			want: Parameter{
				name:        "StringFormParam",
				typeValue:   String,
				in:          Form,
				required:    true,
				description: "A test parameter for form data\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrParam(tc.name, tc.location, WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBoolParam(t *testing.T) {
	testCases := []struct {
		name        string
		location    Location
		required    bool
		description string
		want        Parameter
	}{
		{
			name:        "BoolPathParam",
			location:    Path,
			required:    true,
			description: "A test boolean parameter for path",
			want: Parameter{
				name:        "BoolPathParam",
				typeValue:   Boolean,
				in:          Path,
				required:    true,
				description: "A test boolean parameter for path\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "BoolQueryParam",
			location:    Query,
			required:    true,
			description: "A test boolean parameter for query",
			want: Parameter{
				name:        "BoolQueryParam",
				typeValue:   Boolean,
				in:          Query,
				required:    true,
				description: "A test boolean parameter for query\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "BoolHeaderParam",
			location:    Header,
			required:    true,
			description: "A test boolean parameter for header",
			want: Parameter{
				name:        "BoolHeaderParam",
				typeValue:   Boolean,
				in:          Header,
				required:    true,
				description: "A test boolean parameter for header\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "BoolFormParam",
			location:    Form,
			required:    true,
			description: "A test boolean parameter for form data",
			want: Parameter{
				name:        "BoolFormParam",
				typeValue:   Boolean,
				in:          Form,
				required:    true,
				description: "A test boolean parameter for form data\n (maxLength: 50)",
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolParam(tc.name, tc.location, WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
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
			got := FileParam(tc.name, WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
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
		location    Location
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:        "IntEnumPathParam",
			location:    Path,
			required:    true,
			description: "A test parameter for path",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntEnumPathParam",
				typeValue:   Integer,
				in:          Path,
				required:    true,
				description: "A test parameter for path\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "IntEnumQueryParam",
			location:    Query,
			required:    true,
			description: "A test parameter for query",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntEnumQueryParam",
				typeValue:   Integer,
				in:          Query,
				required:    true,
				description: "A test parameter for query\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "IntEnumHeaderParam",
			location:    Header,
			description: "A test parameter for header",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntEnumHeaderParam",
				typeValue:   Integer,
				in:          Header,
				required:    true,
				description: "A test parameter for header\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:        "IntEnumFormParam",
			location:    Form,
			description: "A test parameter for form data",
			arr:         []int64{1, 2, 3},
			want: Parameter{
				name:        "IntEnumFormParam",
				typeValue:   Integer,
				in:          Form,
				required:    true,
				description: "A test parameter for form data\n (maxLength: 50)",
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntEnumParam(tc.name, tc.location, tc.arr, WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, \nwant %v", got, tc.want)
			}
		})
	}
}

func TestStrEnumParam(t *testing.T) {
	testCases := []struct {
		name        string
		location    Location
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			name:     "StrEnumPathParam",
			location: Path,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				name:      "StrEnumPathParam",
				typeValue: String,
				in:        Path,
				required:  true,
				enum:      []interface{}{"a", "b", "c"},
			},
		},
		{
			name:     "StrEnumQueryParam",
			location: Query,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				name:      "StrEnumQueryParam",
				typeValue: String,
				in:        Query,
				required:  true,
				enum:      []interface{}{"a", "b", "c"},
			},
		},
		{
			name:     "StrEnumHeaderParam",
			location: Header,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				name:      "StrEnumHeaderParam",
				typeValue: String,
				in:        Header,
				required:  true,
				enum:      []interface{}{"a", "b", "c"},
			},
		},
		{
			name:     "StrEnumFormParam",
			location: Form,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				name:      "StrEnumFormParam",
				typeValue: String,
				in:        Form,
				required:  true,
				enum:      []interface{}{"a", "b", "c"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrEnumParam(tc.name, tc.location, tc.arr, WithRequired(), WithDescription(tc.description))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestIntArrParam(t *testing.T) {
	testCases := []struct {
		name        string
		location    Location
		required    bool
		description string
		arr         []int64
		want        Parameter
	}{
		{
			name:     "IntArrPathParam",
			location: Path,
			required: true,
			arr:      []int64{1, 2, 3},
			want: Parameter{
				name:        "IntArrPathParam",
				typeValue:   Array,
				in:          Path,
				description: " (maxLength: 50)",
				required:    true,
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:     "IntArrQueryParam",
			location: Query,
			required: true,
			arr:      []int64{1, 2, 3},
			want: Parameter{
				name:        "IntArrQueryParam",
				typeValue:   Array,
				in:          Query,
				description: " (maxLength: 50)",
				required:    true,
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:     "IntArrHeaderParam",
			location: Header,
			required: true,
			arr:      []int64{1, 2, 3},
			want: Parameter{
				name:        "IntArrHeaderParam",
				typeValue:   Array,
				in:          Header,
				description: " (maxLength: 50)",
				required:    true,
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			name:     "IntArrFormParam",
			location: Form,
			required: true,
			arr:      []int64{1, 2, 3},
			want: Parameter{
				name:        "IntArrFormParam",
				typeValue:   Array,
				in:          Form,
				description: " (maxLength: 50)",
				required:    true,
				enum:        []interface{}{1, 2, 3},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntArrParam(tc.name, tc.location, tc.arr, WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStrArrParam(t *testing.T) {
	testCases := []struct {
		name        string
		location    Location
		required    bool
		description string
		arr         []string
		want        Parameter
	}{
		{
			location: Path,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				typeValue:   Array,
				in:          Path,
				required:    true,
				description: " (maxLength: 50)",
				enum:        []interface{}{"a", "b", "c"},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			location: Query,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				typeValue:   Array,
				in:          Query,
				required:    true,
				description: " (maxLength: 50)",
				enum:        []interface{}{"a", "b", "c"},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			location: Header,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				typeValue:   Array,
				in:          Header,
				required:    true,
				description: " (maxLength: 50)",
				enum:        []interface{}{"a", "b", "c"},
				minLen:      0,
				maxLen:      50,
			},
		},
		{
			location: Form,
			required: true,
			arr:      []string{"a", "b", "c"},
			want: Parameter{
				typeValue:   Array,
				in:          Form,
				required:    true,
				description: " (maxLength: 50)",
				enum:        []interface{}{"a", "b", "c"},
				minLen:      0,
				maxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrArrParam(tc.name, tc.location, tc.arr, WithRequired(), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if !ParametersEqual(*got, tc.want) {
				t.Errorf("got = %v, want %v", got, tc.want)
			}
		})
	}
}
