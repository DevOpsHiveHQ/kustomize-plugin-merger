package main

import (
	"log"

	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"

	"github.com/DevOpsHiveHQ/kustomize-plugin-merger/pkg/merger"
)

func main() {
	rlSource := &kio.ByteReadWriter{}
	processor := &framework.VersionedAPIProcessor{FilterProvider: framework.GVKFilterMap{
		merger.ResourceKind: {
			merger.ResourceGroup + "/" + merger.ResourceVersion: &merger.Merger{},
		},
	}}
	if err := framework.Execute(processor, rlSource); err != nil {
		log.Fatal(err)
	}
}
