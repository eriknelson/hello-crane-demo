package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	transformtypes "github.com/konveyor/crane-lib/transform/types"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
	cli.RunAndExit(cli.NewCustomPlugin("RemoveStatus", "v1", nil, Run))
}

func Run(u *unstructured.Unstructured, extras map[string]string) (transform.PluginResponse, error) {
	// plugin writers need to write custome code here.
	patch, err := RemoveStatus(*u)

	if err != nil {
		return transform.PluginResponse{}, err
	}
	return transform.PluginResponse{
		Version:    "v1",
		IsWhiteOut: false,
		Patches:    patch,
	}, nil
}

func RemoveStatus(u unstructured.Unstructured) (jsonpatch.Patch, error) {
	jsonPatch := jsonpatch.Patch{}
	hasStatus, _ := transformtypes.HasStatusObject(u)
	if hasStatus {
		patchJSON := fmt.Sprintf(`[
				{ "op": "remove", "path": "/status"}
				]`)
		patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
		if err != nil {
			return nil, err
		}
		jsonPatch = append(jsonPatch, patch...)
	}
	return jsonPatch, nil
}
