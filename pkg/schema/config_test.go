package schema_test

import (
	"testing"

	"github.com/algolia/grue/pkg/schema"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	for _, tc := range []struct {
		Path string
		Err  string
	}{
		{Path: "testdata/valid.yaml", Err: ""},
		{Path: "testdata/missing_context.yaml", Err: "Build.Artifacts.0.Context: non zero value required"},
		{Path: "testdata/missing_image.yaml", Err: "Build.Artifacts.0.Image: non zero value required"},
		{Path: "testdata/missing_manifests.yaml", Err: "Deploy.Clusters.0.Manifests: non zero value required"},
	} {
		t.Run(tc.Path, func(t *testing.T) {
			_, err := schema.New(tc.Path)
			if tc.Err == "" {
				require.Nil(t, err)
			} else {
				require.Equal(t, tc.Err, err.Error())
			}
		})
	}
}
