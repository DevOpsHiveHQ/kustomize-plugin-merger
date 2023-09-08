package main

import (
	"log"

	"github.com/aabouzaid/kustomize-plugin-merger/pkg/merger"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/kio"
)

func main() {
	rlSource := &kio.ByteReadWriter{}
	processor := &framework.VersionedAPIProcessor{FilterProvider: framework.GVKFilterMap{
		merger.ResourceKind: {
			merger.ResourceGroup + "/" + merger.ResourceVersion: &merger.Merger{},
		}}}
	if err := framework.Execute(processor, rlSource); err != nil {
		log.Fatal(err)
	}
}
