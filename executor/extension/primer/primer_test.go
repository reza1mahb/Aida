package primer

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Fantom-foundation/Aida/executor"
	"github.com/Fantom-foundation/Aida/executor/extension"
	"github.com/Fantom-foundation/Aida/logger"
	"github.com/Fantom-foundation/Aida/state"
	"github.com/Fantom-foundation/Aida/utils"
	"go.uber.org/mock/gomock"
)

func TestStateDbPrimerExtension_NoPrimerIsCreatedIfDisabled(t *testing.T) {
	cfg := &utils.Config{}
	cfg.SkipPriming = true

	ext := MakeStateDbPrimer[any](cfg)
	if _, ok := ext.(extension.NilExtension[any]); !ok {
		t.Errorf("Primer is enabled although not set in configuration")
	}

}

func TestStateDbPrimerExtension_PrimingDoesNotTriggerForExistingStateDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := logger.NewMockLogger(ctrl)

	cfg := &utils.Config{}
	cfg.SkipPriming = false
	cfg.IsExistingStateDb = true

	log.EXPECT().Warning("Skipping priming due to usage of pre-existing StateDb")

	ext := makeStateDbPrimer[any](cfg, log)

	ext.PreRun(executor.State[any]{}, nil)

}

func TestStateDbPrimerExtension_PrimingDoesTriggerForNonExistingStateDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := logger.NewMockLogger(ctrl)

	cfg := &utils.Config{}
	cfg.SkipPriming = false
	cfg.StateDbSrc = ""
	cfg.First = 2

	log.EXPECT().Noticef("Priming to block %v", cfg.First-1)

	ext := makeStateDbPrimer[any](cfg, log)

	ext.PreRun(executor.State[any]{}, &executor.Context{})
}

func TestStateDbPrimerExtension_AttemptToPrimeBlockZeroDoesNotFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := logger.NewMockLogger(ctrl)

	cfg := &utils.Config{}
	cfg.SkipPriming = false
	cfg.StateDbSrc = ""
	cfg.First = 0

	ext := makeStateDbPrimer[any](cfg, log)

	err := ext.PreRun(executor.State[any]{}, &executor.Context{})
	if err != nil {
		t.Errorf("priming should not happen hence should not fail")
	}
}

// TestStatedb_PrimeStateDB tests priming fresh state DB with randomized world state data
func TestPrime_PrimeStateDB(t *testing.T) {
	log := logger.NewLogger("Warning", "TestPrimeStateDB")
	for _, tc := range utils.GetStateDbTestCases() {
		t.Run(fmt.Sprintf("DB variant: %s; shadowImpl: %s; archive variant: %s", tc.Variant, tc.ShadowImpl, tc.ArchiveVariant), func(t *testing.T) {
			cfg := utils.MakeTestConfig(tc)

			// Initialization of state DB
			sDB, sDbDir, err := utils.PrepareStateDB(cfg)
			defer os.RemoveAll(sDbDir)

			if err != nil {
				t.Fatalf("failed to create state DB: %v", err)
			}

			// Closing of state DB
			defer func(sDB state.StateDB) {
				err = sDB.Close()
				if err != nil {
					t.Fatalf("failed to close state DB: %v", err)
				}
			}(sDB)

			// Generating randomized world state
			ws, _ := utils.MakeWorldState(t)

			pc := utils.NewPrimeContext(cfg, sDB, log)
			// Priming state DB
			pc.PrimeStateDB(ws, sDB)

			// Checks if state DB was primed correctly
			for key, account := range ws {
				if sDB.GetBalance(key).Cmp(account.Balance) != 0 {
					t.Fatalf("failed to prime account balance; Is: %v; Should be: %v", sDB.GetBalance(key), account.Balance)
				}

				if sDB.GetNonce(key) != account.Nonce {
					t.Fatalf("failed to prime account nonce; Is: %v; Should be: %v", sDB.GetNonce(key), account.Nonce)
				}

				if bytes.Compare(sDB.GetCode(key), account.Code) != 0 {
					t.Fatalf("failed to prime account code; Is: %v; Should be: %v", sDB.GetCode(key), account.Code)
				}

				for sKey, sValue := range account.Storage {
					if sDB.GetState(key, sKey) != sValue {
						t.Fatalf("failed to prime account storage; Is: %v; Should be: %v", sDB.GetState(key, sKey), sValue)
					}
				}
			}
		})
	}
}