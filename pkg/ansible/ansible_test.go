package ansible

import (
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testDataDir   = "../../testdata"
	vaultPassword = "pa$$word"
)

func TestFind(t *testing.T) {
	want := []Result{
		{Path: testDataDir + "/defaults/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/group_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/inventories/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/vars/vars.yml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
	}
	results, err := Find(testDataDir, vaultPassword, "test_var")

	assertValuesMatch(t, want, results, err)
}

func TestFind_error(t *testing.T) {
	tests := map[string]struct {
		path     string
		password string
	}{
		"non-existing path":    {path: "does/not/exist", password: vaultPassword},
		"wrong vault password": {path: testDataDir + "/defaults/vault.yaml", password: "wrong password"},
		"non-existing file":    {path: testDataDir + "/defaults/does_not_exist.yaml", password: vaultPassword},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Find(tc.path, tc.password, "test_var")
			assert.Error(t, err)
		})
	}
}

func TestFindRegex(t *testing.T) {
	want := []Result{
		{Path: testDataDir + "/defaults/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/defaults/vault.yaml", Variable: "test_var2", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/group_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/inventories/host_vars/vault.yaml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
		{Path: testDataDir + "/vars/vars.yml", Variable: "test_var", Value: yaml.Node{Value: "value"}},
	}
	results, err := FindRegex(testDataDir, vaultPassword, "^test_.*")

	assertValuesMatch(t, want, results, err)
}

func TestFindRegex_badRegex(t *testing.T) {
	_, err := FindRegex(testDataDir, vaultPassword, "*")
	assert.Error(t, err)
}

func assertValuesMatch(t *testing.T, want []Result, results []Result, err error) {
	assert.NoError(t, err)
	assert.Len(t, results, len(want))

	for _, result := range results {
		matched := false
		for i, wanted := range want {
			if result.Path == wanted.Path && result.Variable == wanted.Variable && result.Value.Value == wanted.Value.Value {
				matched = true
				want = slices.Delete(want, i, i+1)
				break
			}
		}
		assert.Truef(t, matched, "result %v not found in want", result)
	}
	assert.Emptyf(t, want, "not all results found in want: %v", want)
}
