package main

import (
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
	cli.RunAndExit(cli.NewCustomPlugin("SpecifyNamespace", "v1", nil, Run))
}

func Run(u *unstructured.Unstructured, extras map[string]string) (transform.PluginResponse, error) {
	// plugin writers need to write custome code here.
	var patch jsonpatch.Patch
	val, ok := extras["dest-namespace"]
	if !ok {
		// Passthrough if no argument specified
		return transform.PluginResponse{
			Version:    "v1",
			IsWhiteOut: false,
			Patches:    patch,
		}, nil
	}

	// TODO: Check for presence? What if cluster scoped?
	// TODO: Validate val; ex: DNS compliant
	patchJSON := fmt.Sprintf(`[
		{ "op": "replace", "path": "/metadata/namespace", "value": "%v" }
	]`, val)

	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))

	return transform.PluginResponse{
		Version:    "v1",
		IsWhiteOut: false,
		Patches:    patch,
	}, err
}
