package main

import (
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
	cli.RunAndExit(cli.NewCustomPlugin("SpecifyIngress", "v1", nil, Run))
}

func Run(u *unstructured.Unstructured, extras map[string]string) (transform.PluginResponse, error) {
	// plugin writers need to write custome code here.
	var patch jsonpatch.Patch
	var err error
	val, ok := extras["ingress-host"]
	if !ok {
		// Passthrough if no argument specified
		return transform.PluginResponse{
			Version:    "v1",
			IsWhiteOut: false,
			Patches:    patch,
		}, nil
	}

	if u.GetKind() == "Ingress" {
		// TODO: This is not gonna fly with more than one rule
		patchJSON := fmt.Sprintf(`[
			{ "op": "replace", "path": "/spec/rules/0/host", "value": "%v" }
		]`, val)

		patch, err = jsonpatch.DecodePatch([]byte(patchJSON))
	}

	return transform.PluginResponse{
		Version:    "v1",
		IsWhiteOut: false,
		Patches:    patch,
	}, err
}
