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
			"Expected to be %v, but got %v", tt.expected, tt.actual.config)
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
			"Expected to be %v, but got %v", tt.expected, tt.actual)
	}
}

func TestLoadDestinationFile(t *testing.T) {
	tests := []struct {
		actual   resourceInputFiles
		expected string
	}{
		{
			actual: resourceInputFiles{
				Root:        "../../examples/testdata/",
				Destination: "input/base.yaml",
			},
			expected: "metadata.name",
		},
	}

	for _, tt := range tests {
		k := tt.actual.loadDestinationFile()
		assert.Contains(t, k.Keys(), tt.expected,
			"Expected to be %v, but got %v", tt.expected, tt.actual)
	}
}

func TestMergeReplace(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		actual      mergerResource
		expected    map[string]string
		expectedLen int
	}{
		{
			name: "Test replace with overlay",
			desc: "The number of the outputs in 'replace+overlay' should be the same as input sources",
			actual: mergerResource{
				Name: "test-replace-with-overlay",
				Input: resourceInput{
					Method: "overlay",
					Files: resourceInputFiles{
						Root: "../../examples/testdata/",
						Sources: []string{
							"input/dev.yaml",
							"input/stage.yaml",
						},
						Destination: "input/base.yaml",
					},
				},
				Merge:  resourceMerge{Strategy: "replace"},
				Output: resourceOutput{Format: "raw"},
			},
			expected: map[string]string{
				// In the overlay the output key are "resource.input.files.sources[*]".
				"dev.yaml":   "image: dev:1.0.0",
				"stage.yaml": "image: stage:1.0.0",
			},
			expectedLen: 2,
		},
		{
			name: "Test replace with patch",
			desc: "The number of the outputs in 'replace+patch' should always be one",
			actual: mergerResource{
				Name: "test-replace-with-patch",
				Input: resourceInput{
					Method: "patch",
					Files: resourceInputFiles{
						Root: "../../examples/testdata/",
						Sources: []string{
							"input/dev.yaml",
							"input/stage.yaml",
						},
						Destination: "input/base.yaml",
					},
				},
				Merge:  resourceMerge{Strategy: "replace"},
				Output: resourceOutput{Format: "raw"},
			},
			expected: map[string]string{
				// In the patch the output key are "resource.name".
				"test-replace-with-patch": "image: stage:1.0.0",
			},
			expectedLen: 1,
		},
	}

	for _, tt := range tests {
		// Those methods are called in the Default method.
		tt.actual.Input.Files.setRoot()
		tt.actual.Merge.setStrategy()
		tt.actual.Output.items = make(map[string]string)

		// Merge.
		tt.actual.merge()

		// Assert.
		assert.Equal(t, len(tt.actual.Output.items), tt.expectedLen)
		for k, v := range tt.expected {
			output := tt.actual.Output.items[k]
			assert.Contains(t, output, v,
				"Expected to be %v=%v, but got %v", k, v, output)
		}
	}
}
