package main

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
	cli.RunAndExit(cli.NewCustomPlugin("ServicePlugin", "v1", nil, Run))
}

// Removes ExternalIPs for LoadBalancer services
func Run(u *unstructured.Unstructured, extras map[string]string) (transform.PluginResponse, error) {
	// plugin writers need to write custome code here.
	var patch jsonpatch.Patch
	var err error
	if u.GetKind() == "Service" {
		var intPatch jsonpatch.Patch
		if IsLoadBalancerService(*u) {
			patchJSON := fmt.Sprintf(`[
{ "op": "remove", "path": "/spec/externalIPs"}
]`)
			patch, err = jsonpatch.DecodePatch([]byte(patchJSON))
		}
		if !IsServiceClusterIPNone(*u) {
			patchJSON := fmt.Sprintf(`[
{ "op": "remove", "path": "/spec/clusterIP"}
]`)
			intPatch, err = jsonpatch.DecodePatch([]byte(patchJSON))
			patch = append(patch, intPatch...)
		}

		if !IsServiceClusterIPNone(*u) {
			patchJSON := fmt.Sprintf(`[
{ "op": "remove", "path": "/spec/clusterIPs"}
]`)
			intPatch, err = jsonpatch.DecodePatch([]byte(patchJSON))
			patch = append(patch, intPatch...)
		}

		// HACK: Just need to make this go away, smarter ways to do this like making
		// sure its even present. Also probably needs some amount of processing
		// to read through each of the port structs and to check if nodePort is
		// present and then strip. This is just hardcoded for demo purposes.
		patchJSON := fmt.Sprintf(`[
			{ "op": "remove", "path": "/spec/ports/0/nodePort"}
		]`)
		intPatch, err = jsonpatch.DecodePatch([]byte(patchJSON))
		patch = append(patch, intPatch...)
	}
	return transform.PluginResponse{
		Version:    "v1",
		IsWhiteOut: false,
		Patches:    patch,
	}, err
}

func IsLoadBalancerService(u unstructured.Unstructured) bool {
	if u.GetKind() != "Service" {
		return false
	}
	// Get Spec
	spec, ok := u.UnstructuredContent()["spec"]
	if !ok {
		return false
	}

	specMap, ok := spec.(map[string]interface{})
	if !ok {
		return false
	}
	// Get type
	serviceType, ok := specMap["type"]
	if !ok {
		return false
	}
	return serviceType == "LoadBalancer"
}

func IsServiceClusterIPNone(u unstructured.Unstructured) bool {
	if u.GetKind() != "Service" {
		return false
	}
	// Get Spec
	spec, ok := u.UnstructuredContent()["spec"]
	if !ok {
		return false
	}

	specMap, ok := spec.(map[string]interface{})
	if !ok {
		return false
	}
	// Get type
	clusterIP, ok := specMap["clusterIP"]
	if !ok {
		return false
	}
	return clusterIP == "None"
}

func IsServiceClusterIPsNone(u unstructured.Unstructured) bool {
	if u.GetKind() != "Service" {
		return false
	}
	// Get Spec
	spec, ok := u.UnstructuredContent()["spec"]
	if !ok {
		return false
	}

	specMap, ok := spec.(map[string]interface{})
	if !ok {
		return false
	}
	// Get type
	clusterIPs, ok := specMap["clusterIPs"]
	if !ok {
		return false
	}
	return clusterIPs == "None"
}
