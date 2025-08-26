// The main logic for Merger plugin.
package merger

import (
	"log"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	koanfYaml "github.com/knadh/koanf/parsers/yaml"
	koanfFile "github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// setStrategy sets "mergo" merging strategy, mainly for merging arrays.
func (rm *resourceMerge) setStrategy() {
	switch rm.Strategy {
	case Replace:
		rm.config = mergo.WithOverride
	case Append:
		rm.config = mergo.WithAppendSlice
	case Combine:
		rm.config = mergo.WithSliceDeepCopy
	}
}

func (rif *resourceInputFiles) setRoot() {
	if rif.Root != "" {
		rif.Root = strings.TrimSuffix(rif.Root, "/") + "/"
	}
}

func (rif *resourceInputFiles) loadDestinationFile() *koanf.Koanf {
	k := koanf.New(".")
	dstFile := rif.Root + rif.Destination
	if err := k.Load(koanfFile.Provider(dstFile), koanfYaml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return k
}

// merge performs the actual merging of configuration files from resourceInputFiles sources.
func (r *mergerResource) merge() {
	// TODO: Simplify/split the logic in merge method.
	k := r.Input.Files.loadDestinationFile()
	fileKey := r.Name

	for _, srcFile := range r.Input.Files.Sources {
		srcFile = r.Input.Files.Root + srcFile

		if r.Input.Method == Overlay {
			k = r.Input.Files.loadDestinationFile()
			fileKey = filepath.Base(srcFile)
		}

		err := k.Load(koanfFile.Provider(srcFile), koanfYaml.Parser(),
			koanf.WithMergeFunc(func(src, dst map[string]interface{}) error {
				if err := mergo.Merge(&dst, src, mergo.WithOverride, r.Merge.config); err != nil {
					log.Fatalf("Error merging config: %v", err)
				}
				return nil
			}))
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}
		//
		mergedData, err := k.Marshal(koanfYaml.Parser())
		if err != nil {
			log.Fatalf("Error marshaling yaml: %v", err)
		}
		r.Output.items[fileKey] = string(mergedData)
	}
}

func (r *mergerResource) export() []*yaml.RNode {
	rlItems := []*yaml.RNode{}
	switch r.Output.Format {
	case Raw:
		for _, value := range r.Output.items {
			rlItems = append(rlItems, yaml.MustParse(value))
		}

	case ConfigMap:
		rNode, err := yaml.FromMap(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":        r.Name,
				"annotations": r.Output.Annotations,
			},
		})
		if err != nil {
			log.Fatalf("Error creating ConfigMap: %v", err)
		}
		if err := rNode.LoadMapIntoConfigMapData(r.Output.items); err != nil {
			log.Fatalf("Error creating ConfigMap data: %v", err)
		}
		rlItems = append(rlItems, rNode)

	case Secret:
		rNode, err := yaml.FromMap(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":        r.Name,
				"annotations": r.Output.Annotations,
			},
		})
		if err != nil {
			log.Fatalf("Error creating Secret: %v", err)
		}
		if err := rNode.LoadMapIntoSecretData(r.Output.items); err != nil {
			log.Fatalf("Error creating Secret data: %v", err)
		}
		rlItems = append(rlItems, rNode)
	}

	return rlItems
}
