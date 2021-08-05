package main

import (
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/konveyor/crane-lib/transform"
	"github.com/konveyor/crane-lib/transform/cli"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sort"
)

func main() {
	cli.RunAndExit(cli.NewCustomPlugin("HelloCraneWhiteout", "v1", nil, Run))
}

func contains(s []string, searchterm string) bool {
	i := sort.SearchStrings(s, searchterm)
	return i < len(s) && s[i] == searchterm
}

func Run(u *unstructured.Unstructured, extras map[string]string) (transform.PluginResponse, error) {
	// plugin writers need to write custome code here.
	var patch jsonpatch.Patch
	var whiteout bool

	// NOTE: Don't repeat this! I happen to know that none of the kinds in this
	// specific demo are relevant and should not be migrated. If you implement a
	// plugin and follow this, it's likely you'll end up erasing important instances.
	whiteoutKinds := []string{
		"Endpoints",
		"ClusterServiceVersion",
		"EndpointSlice",
		"ConfigMap",
		"RoleBinding",
		"ServiceAccount",
		"Secret",
	}

	sort.Strings(whiteoutKinds)

	kind := u.GetKind()
	if contains(whiteoutKinds, kind) {
		whiteout = true
	}

	return transform.PluginResponse{
		Version:    "v1",
		IsWhiteOut: whiteout,
		Patches:    patch,
	}, nil
}
