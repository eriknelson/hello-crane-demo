package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
	cli.RunAndExit(cli.NewCustomPlugin("SanitizeMetadata", "v1", nil, Run))
}

// Removes ExternalIPs for LoadBalancer services
func Run(u *unstructured.Unstructured, extras map[string]string) (transform.PluginResponse, error) {
	// plugin writers need to write custome code here.
	var patch jsonpatch.Patch
	var err error

	patchJSON := fmt.Sprintf(`[
{ "op": "remove", "path": "/metadata/managedFields"},
{ "op": "remove", "path": "/metadata/uid"},
{ "op": "remove", "path": "/metadata/resourceVersion"},
{ "op": "remove", "path": "/metadata/creationTimestamp"}
]`)
	patch, err = jsonpatch.DecodePatch([]byte(patchJSON))

	return transform.PluginResponse{
		Version:    "v1",
		IsWhiteOut: false,
		Patches:    patch,
	}, err
}
