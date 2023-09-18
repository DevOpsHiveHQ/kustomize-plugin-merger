package merger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

func TestMerger(t *testing.T) {
	// Create a Merger resource.
	merger := &Merger{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-merger",
		},
		Spec: mergerSpec{
			Resources: []mergerResource{
				{
					Name: "my-resource",
					Input: resourceInput{
						Method: Overlay,
						Files: resourceInputFiles{
							Root: "../../examples/testdata",
							Sources: []string{
								"input/dev.yaml",
								"input/stage.yaml",
							},
							Destination: "input/base.yaml",
						},
					},
					Merge: resourceMerge{
						Strategy: Append,
					},
					Output: resourceOutput{
						Format: Raw,
					},
				},
			},
		},
	}

	//
	// Schema.
	_, err := merger.Schema()
	if err != nil {
		t.Errorf("Failed to load schema: %v", err)
	}

	//
	// Default.
	err = merger.Default()
	if err != nil {
		t.Errorf("Failed to set defaults: %v", err)
	}
	// The trailing slash should be added to the root.
	assert.Equal(t, "../../examples/testdata/", merger.Spec.Resources[0].Input.Files.Root)
	// Create the output items as an empty map.
	assert.Equal(t, map[string]string{}, merger.Spec.Resources[0].Output.items)

	//
	// Validate.
	err = merger.Validate()
	if err != nil {
		t.Errorf("Failed to validate: %v", err)
	}

	//
	// Filter.
	rlItemsIn := []*yaml.RNode{}
	var rlItemsOut []*yaml.RNode

	rlItemsOut, err = merger.Filter(rlItemsIn)
	if err != nil {
		t.Errorf("Failed to filter resource list: %v", err)
	}
	// Output items should be two since the input method is "overlay".
	assert.Len(t, rlItemsOut, 2)

	// Verify that the filtered resource list contains the merged resource.
	// NOTE: Since the items are generated from a Map; hence the output order is no preserved.
	assert.Equal(t, "Pod", rlItemsOut[0].GetKind())
}
