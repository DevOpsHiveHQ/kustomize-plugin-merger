package merger

import (
	"reflect"
	"runtime"
	"testing"

	"dario.cat/mergo"
	"github.com/stretchr/testify/assert"
)

func TestSetMergeStrategy(t *testing.T) {
	tests := []struct {
		actual   resourceMerge
		expected func(*mergo.Config)
	}{
		{
			actual: resourceMerge{
				Strategy: "replace",
			},
			expected: mergo.WithOverride,
		},
		{
			actual: resourceMerge{
				Strategy: "append",
			},
			expected: mergo.WithAppendSlice,
		},
		{
			actual: resourceMerge{
				Strategy: "combine",
			},
			expected: mergo.WithSliceDeepCopy,
		},
	}

	for _, tt := range tests {
		tt.actual.setStrategy()
		actual := runtime.FuncForPC(reflect.ValueOf(tt.actual.config).Pointer()).Name()
		expected := runtime.FuncForPC(reflect.ValueOf(tt.expected).Pointer()).Name()

		assert.Equal(t, expected, actual,
			"Expected config to be %v, but got %v", tt.expected, tt.actual.config)
	}
}

func TestSetRoot(t *testing.T) {
	tests := []struct {
		actual   resourceInputFiles
		expected string
	}{
		{
			actual: resourceInputFiles{
				Root: "/mnt",
			},
			expected: "/mnt/",
		},
		{
			actual: resourceInputFiles{
				Root: "/mnt/",
			},
			expected: "/mnt/",
		},
		{
			actual: resourceInputFiles{
				Root: "my-dir",
			},
			expected: "my-dir/",
		},
	}

	for _, tt := range tests {
		tt.actual.setRoot()
		assert.Equal(t, tt.expected, tt.actual.Root,
			"Expected root to be %v, but got %v", tt.expected, tt.actual)
	}
}
