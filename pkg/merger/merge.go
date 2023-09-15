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

// setMergeStrategy sets "mergo" merging strategy, mainly for merging arrays.
func (r *mergerResource) setMergeStrategy() {
	switch r.Merge.Strategy {
	case Replace:
		r.Merge.config = mergo.WithOverride
	case Append:
		r.Merge.config = mergo.WithAppendSlice
	case Combine:
		r.Merge.config = mergo.WithSliceDeepCopy
	default:
		r.Merge.config = func(*mergo.Config) {}
	}
}

func (r *mergerResource) setInputFilesRoot() {
	if r.Input.Files.Root != "" {
		r.Input.Files.Root = strings.TrimSuffix(r.Input.Files.Root, "/") + "/"
	}
}

func (r *mergerResource) loadDestinationFile() *koanf.Koanf {
	k := koanf.New(".")
	dstFile := r.Input.Files.Root + r.Input.Files.Destination
	if err := k.Load(koanfFile.Provider(dstFile), koanfYaml.Parser()); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return k
}

// merge performs the actual merging of configuration files from resourceInputFiles sources.
func (r *mergerResource) merge() {
	// TODO: Simplify/split the logic in merge method.
	k := r.loadDestinationFile()
	fileKey := r.Name

	for _, srcFile := range r.Input.Files.Sources {
		srcFile = r.Input.Files.Root + srcFile

		if r.Input.Method == Overlay {
			k = r.loadDestinationFile()
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
		rNode, _ := yaml.FromMap(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]string{
				"name": r.Name,
			},
		})
		if err := rNode.LoadMapIntoConfigMapData(r.Output.items); err != nil {
			log.Fatalf("Error creating ConfigMap data: %v", err)
		}
		rlItems = append(rlItems, rNode)

	case Secret:
		rNode, _ := yaml.FromMap(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]string{
				"name": r.Name,
			},
		})
		if err := rNode.LoadMapIntoSecretData(r.Output.items); err != nil {
			log.Fatalf("Error creating Secret data: %v", err)
		}
		rlItems = append(rlItems, rNode)
	}

	return rlItems
}
