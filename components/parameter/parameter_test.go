package parameter

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// TODO create a comprehensive test suite for the 'Fields' types
func TestParams(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		// IntParam
		{
			name:        "testParam",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{Min: 0, Max: 100}},
			want: Parameter{
				Name:        "testParam",
				Type:        Integer,
				In:          Path,
				Required:    true,
				Description: "A test parameter\n (max: 100)",
				Min:         0,
				Max:         100,
			},
		},
		{
			name:        "testParam",
			required:    false,
			description: "A test parameter",
			args:        []Fields{{Min: 0, Max: 100}},
			want: Parameter{
				Name:        "testParam",
				Type:        Integer,
				In:          Path,
				Required:    false,
				Description: "A test parameter\n (max: 100)",
				Min:         0,
				Max:         100,
			},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntParam(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMax(100))
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("StrParam() = %v, want %v", got, tc.want)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("JsonSwagger() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestStrParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		// StrParam
		{
			name:        "stringParam",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "stringParam",
				Type:        String,
				In:          Path,
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrParam(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StrParam() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestBoolParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "boolParam",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "boolParam",
				Type:        Boolean,
				In:          Path,
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolParam(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("BoolParam() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestFileParam(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "FileParam",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "FileParam",
				Type:        File,
				In:          Form,
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := FileParam(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FileParam() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestIntQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntQuery",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntQuery",
				Type:        Integer,
				In:          "query",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntQuery(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("IntQuery() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestStrQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrQuery",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrQuery",
				Type:        String,
				In:          "query",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrQuery(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StrQuery() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestBoolQuery(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "BoolQuery",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "BoolQuery",
				Type:        Boolean,
				In:          "query",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolQuery(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("BoolQuery() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestIntHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntHeader",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntHeader",
				Type:        Integer,
				In:          "header",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntHeader(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("IntHeader() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestStrHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrHeader",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrHeader",
				Type:        String,
				In:          "header",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrHeader(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StrHeader() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}

func TestBoolHeader(t *testing.T) {
	testCases := []struct {
		name        string
		required    bool
		description string
		args        []Fields
		want        Parameter
	}{
		{
			name:        "BoolHeader",
			required:    true,
			description: "A test parameter",
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "BoolHeader",
				Type:        Boolean,
				In:          "header",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BoolHeader(tc.name, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("BoolHeader() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntEnumParam",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntEnumParam",
				Type:        Integer,
				In:          Path,
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{1, 2, 3},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntEnumParam(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("IntEnumParam() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrEnumParam",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrEnumParam",
				Type:        String,
				In:          Path,
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{"a", "b", "c"},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrEnumParam(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("StrEnumParam() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntEnumQuery",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntEnumQuery",
				Type:        Integer,
				In:          "query",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{1, 2, 3},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntEnumQuery(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("StrEnumParam() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrEnumQuery",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrEnumQuery",
				Type:        String,
				In:          "query",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{"a", "b", "c"},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrEnumQuery(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("StrEnumQuery() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntEnumHeader",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntEnumHeader",
				Type:        Integer,
				In:          "header",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{1, 2, 3},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntEnumHeader(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("IntEnumHeader() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrEnumHeader",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrEnumHeader",
				Type:        String,
				In:          "header",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{"a", "b", "c"},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrEnumHeader(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StrEnumHeader() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntArrParam",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntArrParam",
				Type:        Array,
				In:          Path,
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{1, 2, 3},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntArrParam(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("IntArrParam() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrArrParam",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrArrParam",
				Type:        Array,
				In:          Path,
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
				Enum:        []interface{}{string("a"), string("b"), string("c")},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrArrParam(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StrArrParam() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntArrQuery",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntArrQuery",
				Type:        Array,
				In:          "query",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{1, 2, 3},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntArrQuery(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("IntArrQuery() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrArrQuery",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrArrQuery",
				Type:        Array,
				In:          "query",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
				Enum:        []interface{}{string("a"), string("b"), string("c")},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrArrQuery(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StrArrQuery() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "IntArrHeader",
			required:    true,
			description: "A test parameter",
			arr:         []int64{1, 2, 3},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "IntArrHeader",
				Type:        Array,
				In:          "header",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				Enum:        []interface{}{1, 2, 3},
				MinLen:      0,
				MaxLen:      50,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntArrHeader(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got, cmpopts.IgnoreFields(Parameter{}, "Enum")); diff != "" {
				t.Errorf("IntArrHeader() mismatch (-expected +got):\n%s", diff)
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
		args        []Fields
		want        Parameter
	}{
		{
			name:        "StrArrHeader",
			required:    true,
			description: "A test parameter",
			arr:         []string{"a", "b", "c"},
			args:        []Fields{{MinLen: 0, MaxLen: 50}},
			want: Parameter{
				Name:        "StrArrHeader",
				Type:        Array,
				In:          "header",
				Required:    true,
				Description: "A test parameter\n (maxLength: 50)",
				MinLen:      0,
				MaxLen:      50,
				Enum:        []interface{}{string("a"), string("b"), string("c")},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := StrArrHeader(tc.name, tc.arr, WithRequired(tc.required), WithDescription(tc.description), WithMin(0), WithMaxLen(50))
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StrArrHeader() mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
