package assets

import "testing"

func TestGetAsset(t *testing.T) {
	s, _ := GetAsset("ias/cloudformation/ecs.yml")

	if s == "" {
		t.Log("error retrieving assets")
		t.Fail()
	}

}
func TestGetAssetFail(t *testing.T) {
	_, err := GetAsset("ias/cloudformation/ecsmis.yml")

	if err == nil {
		t.Log("error not returned")
		t.Fail()
	}

}
