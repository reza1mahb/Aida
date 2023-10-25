package extension

import (
	"testing"

	"github.com/Fantom-foundation/Aida/utils"
)

func TestAidaDbManager_NoManagerIsCreatedIfPathIsNotProvided(t *testing.T) {
	cfg := &utils.Config{}
	ext := MakeAidaDbManager[any](cfg)

	if _, ok := ext.(NilExtension[any]); !ok {
		t.Errorf("manager is enabled although not set in configuration")
	}
}
