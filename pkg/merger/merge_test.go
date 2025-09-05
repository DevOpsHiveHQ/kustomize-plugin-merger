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

type testMergeCase struct {
	name        string
	key         string
	contains    string
	notContains string
}

func TestMergeReplace(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		actual      mergerResource
		expected    []testMergeCase
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
			// In the overlay the output key are "resource.input.files.sources[*]".
			expected: []testMergeCase{
				{"src01", "dev.yaml", "image: dev:1.0.0", "image: sidecar"},
				{"src02", "stage.yaml", "image: stage:1.0.0", "image: sidecar"},
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
			// In the patch the output key are "resource.name".
			expected: []testMergeCase{
				{"src", "test-replace-with-patch", "image: stage:1.0.0", "image: sidecar"},
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
		for _, tmc := range tt.expected {
			output := tt.actual.Output.items[tmc.key]
			assert.Contains(t, output, tmc.contains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.contains, output)
			assert.NotContains(t, output, tmc.notContains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.notContains, output)
		}
	}
}

func TestMergeAppend(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		actual      mergerResource
		expected    []testMergeCase
		expectedLen int
	}{
		{
			name: "Test append with overlay",
			desc: "The number of the outputs in 'append+overlay' should be the same as input sources",
			actual: mergerResource{
				Name: "test-append-with-overlay",
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
				Merge:  resourceMerge{Strategy: "append"},
				Output: resourceOutput{Format: "raw"},
			},
			// In the overlay the output key are "resource.input.files.sources[*]".
			expected: []testMergeCase{
				{"src01", "dev.yaml", "image: dev:1.0.0", "image: stage:1.0.0"},
				{"src01", "dev.yaml", "image: sidecar", "image: stage:1.0.0"},
				{"src02", "stage.yaml", "image: stage:1.0.0", "image: dev:1.0.0"},
				{"src02", "stage.yaml", "image: sidecar", "image: dev:1.0.0"},
			},
			expectedLen: 2,
		},
		{
			name: "Test append with patch",
			desc: "The number of the outputs in 'append+patch' should always be one",
			actual: mergerResource{
				Name: "test-append-with-patch",
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
				Merge:  resourceMerge{Strategy: "append"},
				Output: resourceOutput{Format: "raw"},
			},
			// In the patch the output key are "resource.name".
			expected: []testMergeCase{
				{"src", "test-append-with-patch", "image: sidecar", "dummy"},
				{"src", "test-append-with-patch", "image: dev:1.0.0", "dummy"},
				{"src", "test-append-with-patch", "image: stage:1.0.0", "dummy"},
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
		for _, tmc := range tt.expected {
			output := tt.actual.Output.items[tmc.key]
			assert.Contains(t, output, tmc.contains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.contains, output)
			assert.NotContains(t, output, tmc.notContains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.notContains, output)
		}
	}
}

func TestMergeCombine(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		actual      mergerResource
		expected    []testMergeCase
		expectedLen int
	}{
		{
			name: "Test combine with overlay",
			desc: "The number of the outputs in 'combine+overlay' should be the same as input sources",
			actual: mergerResource{
				Name: "test-combine-with-overlay",
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
				Merge:  resourceMerge{Strategy: "combine"},
				Output: resourceOutput{Format: "raw"},
			},
			// In the overlay the output key are "resource.input.files.sources[*]".
			expected: []testMergeCase{
				{"src01.1", "dev.yaml", "image: dev:1.0.0", "image: stage:1.0.0"},
				{"src01.2", "dev.yaml", "image: dev:1.0.0", "image: sidecar"},
				{"src02.1", "stage.yaml", "image: stage:1.0.0", "image: dev:1.0.0"},
				{"src02.2", "stage.yaml", "image: stage:1.0.0", "image: sidecar"},
			},
			expectedLen: 2,
		},
		{
			name: "Test combine with patch",
			desc: "The number of the outputs in 'combine+patch' should always be one",
			actual: mergerResource{
				Name: "test-combine-with-patch",
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
				Merge:  resourceMerge{Strategy: "combine"},
				Output: resourceOutput{Format: "raw"},
			},
			// In the patch the output key are "resource.name".
			expected: []testMergeCase{
				{"src", "test-combine-with-patch", "image: stage:1.0.0", "image: sidecar"},
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
		for _, tmc := range tt.expected {
			output := tt.actual.Output.items[tmc.key]
			assert.Contains(t, output, tmc.contains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.contains, output)
			assert.NotContains(t, output, tmc.notContains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.notContains, output)
		}
	}
}

func TestOutputRaw(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		actual      mergerResource
		expected    []testMergeCase
		expectedLen int
	}{
		{
			name: "Test Raw output with overlay",
			desc: "Raw output is used when the input is a manifest.",
			actual: mergerResource{
				Name: "test-configmap",
				Output: resourceOutput{
					Format: "raw",
					items: map[string]string{
						"iampolicy.yaml": `---
              apiVersion: iam.cnrm.cloud.google.com/v1beta1
              kind: IAMPolicy
              metadata:
                name: storage-admin-policy
            `,
					},
				},
			},
			expected: []testMergeCase{
				{"src01.1", "iampolicy.yaml", "kind: IAMPolicy", ""},
				{"src01.2", "iampolicy.yaml", "storage-admin-policy", ""},
			},
		},
	}

	for _, tt := range tests {
		// Export.
		output := tt.actual.export()[0]
		outputManifest, _ := output.String()

		// Assert.
		for _, tmc := range tt.expected {
			assert.Contains(t, outputManifest, tmc.contains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.contains, outputManifest)
		}
	}
}

func TestOutputConfigMap(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		actual      mergerResource
		expected    []testMergeCase
		expectedLen int
	}{
		{
			name: "Test ConfigMap output with overlay",
			desc: "The number of the outputs in 'ConfigMap+overlay' should be the same as input sources",
			actual: mergerResource{
				Name: "test-output-configmap",
				Output: resourceOutput{
					Format: "configmap",
					items: map[string]string{
						"prom.yaml": `
							global:
								scrape_interval: 15s
								evaluation_interval: 15s
						`,
					},
				},
			},
			expected: []testMergeCase{
				{"src01.1", "prom.yaml", "kind: ConfigMap", ""},
				{"src01.2", "prom.yaml", "name: test-output-configmap", ""},
				{"src01.3", "prom.yaml", "evaluation_interval", ""},
			},
		},
		{
			name: "Test ConfigMap output with overlay+annotations",
			desc: "The number of the outputs in 'ConfigMap+overlay' should be the same as input sources",
			actual: mergerResource{
				Name: "test-output-configmap",
				Output: resourceOutput{
					Format: "configmap",
					Annotations: map[string]string{
						"kustomize.config.k8s.io/needs-hash": "true",
					},
					items: map[string]string{
						"prom.yaml": `
							global:
								scrape_interval: 15s
								evaluation_interval: 15s
						`,
					},
				},
			},
			expected: []testMergeCase{
				{"src01.1", "prom.yaml", "kind: ConfigMap", ""},
				{"src01.2", "prom.yaml", "name: test-output-configmap", ""},
				{"src01.3", "prom.yaml", "evaluation_interval", ""},
				{"src01.4", "prom.yaml", "kustomize.config.k8s.io/needs-hash: \"true\"", ""},
			},
		},
	}

	for _, tt := range tests {
		// Export.
		output := tt.actual.export()[0]
		outputData := output.GetDataMap()
		outputManifest, _ := output.String()

		// Assert.
		for _, tmc := range tt.expected {
			assert.Contains(t, outputData, tmc.key)
			assert.Contains(t, outputManifest, tmc.contains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.contains, outputManifest)
		}
	}
}

func TestOutputSecret(t *testing.T) {
	tests := []struct {
		name        string
		desc        string
		actual      mergerResource
		expected    []testMergeCase
		expectedLen int
	}{
		{
			name: "Test Secret output with overlay",
			desc: "The number of the outputs in 'Secret+overlay' should be the same as input sources",
			actual: mergerResource{
				Name: "test-output-secret",
				Output: resourceOutput{
					Format: "secret",
					items: map[string]string{
						"prom.yaml": `
							global:
								scrape_interval: 15s
								evaluation_interval: 15s
						`,
					},
				},
			},
			expected: []testMergeCase{
				{"src01.2", "prom.yaml", "kind: Secret", ""},
				{"src01.2", "prom.yaml", "name: test-output-secret", ""},
				{"src01.3", "prom.yaml", "CgkJCQkJCQlnbG9iYWw6CgkJCQkJCQkJc2NyYXBlX2ludGVydmFsOiAxNXMKCQkJCQkJCQ", ""},
			},
		},
		{
			name: "Test Secret output with overlay+annotations",
			desc: "The number of the outputs in 'Secret+overlay' should be the same as input sources",
			actual: mergerResource{
				Name: "test-output-secret",
				Output: resourceOutput{
					Format: "secret",
					Annotations: map[string]string{
						"kustomize.config.k8s.io/needs-hash": "false",
					},
					items: map[string]string{
						"prom.yaml": `
							global:
								scrape_interval: 15s
								evaluation_interval: 15s
						`,
					},
				},
			},
			expected: []testMergeCase{
				{"src01.2", "prom.yaml", "kind: Secret", ""},
				{"src01.2", "prom.yaml", "name: test-output-secret", ""},
				{"src01.3", "prom.yaml", "CgkJCQkJCQlnbG9iYWw6CgkJCQkJCQkJc2NyYXBlX2ludGVydmFsOiAxNXMKCQkJCQkJCQ", ""},
				{"src01.4", "prom.yaml", "kustomize.config.k8s.io/needs-hash: \"false\"", ""},
			},
		},
	}

	for _, tt := range tests {
		// Export.
		output := tt.actual.export()[0]
		outputData := output.GetDataMap()
		outputManifest, _ := output.String()

		// Assert.
		for _, tmc := range tt.expected {
			assert.Contains(t, outputData, tmc.key)
			assert.Contains(t, outputManifest, tmc.contains,
				"Expected to be %v %v=%v, but got %v", tmc.name, tmc.key, tmc.contains, outputManifest)
		}
	}
}
