package kubernetes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const buildOutputTemplate = `RUN apt-get update && \
apt-get install -y apt-transport-https curl && \
curl -o kubectl https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/kubectl && \
mv kubectl /usr/local/bin && \
chmod a+x /usr/local/bin/kubectl
`

func TestMixin_Build(t *testing.T) {
	t.Run("build with the default Kubernetes version", func(t *testing.T) {
		m := NewTestMixin(t)
		m.Debug = false
		err := m.Build()
		require.NoError(t, err)

		wantOutput := fmt.Sprintf(buildOutputTemplate, "v1.15.5")

		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with custom Kubernetes version", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-version.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)
		err = m.Build()
		require.NoError(t, err)

		wantOutput := fmt.Sprintf(buildOutputTemplate, "v1.18.3")

		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

	t.Run("build with custom Kubernetes version", func(t *testing.T) {
		b, err := ioutil.ReadFile("testdata/build-input-with-namespaces.yaml")
		require.NoError(t, err)

		m := NewTestMixin(t)
		m.Debug = false
		m.In = bytes.NewReader(b)
		err = m.Build()
		require.NoError(t, err)

		wantOutput := fmt.Sprintf(buildOutputTemplate, "v1.15.5") +
			"\nRUN kubectl create namespace my-namespace-1 || true; kubectl create namespace my-namespace-2 || true; " +
			"kubectl create namespace my-namespace-3 || true; kubectl create namespace my-namespace-4 || true;"
		gotOutput := m.TestContext.GetOutput()
		assert.Equal(t, wantOutput, gotOutput)
	})

}
