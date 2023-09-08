package merger

import (
	_ "embed"

	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/kustomize/kyaml/errors"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/resid"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

//go:embed "schema/generators.kustomize.aabouzaid.com_mergers.yaml"
var mergerSchemaDefinition string

func (m *Merger) Schema() (*spec.Schema, error) {
	schema, err := framework.SchemaFromFunctionDefinition(
		resid.NewGvk(ResourceGroup, ResourceVersion, ResourceKind),
		mergerSchemaDefinition)
	return schema, errors.WrapPrefixf(err, "failed to parse Merger schema")
}

func (m *Merger) Default() error {
	for index := range m.Spec.Resources {
		// Defaults input files.
		if err := m.Spec.Resources[index].setInputFiles(); err != nil {
			return err
		}
		// Defaults merge strategy.
		m.Spec.Resources[index].setMergeStrategy()
	}
	return nil
}

func (mg *Merger) Validate() error {
	// At the moment this only validates Merger manifest against its openAPIV3Schema.
	return nil
}

func (m *Merger) Filter(rlItems []*yaml.RNode) ([]*yaml.RNode, error) {
	for _, resource := range m.Spec.Resources {
		resource.merge(resource.Input.items)
		rlItems = append(rlItems, resource.Output.rlItems...)
	}

	return rlItems, nil
}
