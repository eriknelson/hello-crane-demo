package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
	cli.RunAndExit(cli.NewCustomPlugin("PVCProcessor", "v1", nil, Run))
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
	patchJSON := fmt.Sprintf(`[
		{ "op": "remove", "path": "/spec/volumeName"}
	]`)

	patch, err := jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		return nil, err
	}

	jsonPatch = append(jsonPatch, patch...)

	patchJSON = fmt.Sprintf(`[
		{ "op": "remove", "path": "/spec/volumeMode"}
	]`)

	patch, err = jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		return nil, err
	}

	jsonPatch = append(jsonPatch, patch...)

	patchJSON = fmt.Sprintf(`[
		{ "op": "remove", "path": "/metadata/annotations"}
	]`)

	patch, err = jsonpatch.DecodePatch([]byte(patchJSON))
	if err != nil {
		return nil, err
	}

	jsonPatch = append(jsonPatch, patch...)
	return jsonPatch, nil
}
