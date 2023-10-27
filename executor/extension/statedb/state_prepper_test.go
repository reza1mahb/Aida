package statedb

import (
	"testing"

	"github.com/Fantom-foundation/Aida/executor"
	"github.com/Fantom-foundation/Aida/state"
	substate "github.com/Fantom-foundation/Substate"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"
)

func TestStatePrepper_PreparesStateBeforeEachTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := state.NewMockStateDB(ctrl)

	allocA := substate.SubstateAlloc{common.Address{1}: nil}
	allocB := substate.SubstateAlloc{common.Address{2}: nil}
	ctx := &executor.Context{State: db}

	gomock.InOrder(
		db.EXPECT().PrepareSubstate(&allocA, uint64(5)),
		db.EXPECT().PrepareSubstate(&allocB, uint64(7)),
	)

	prepper := MakeStateDbPrepper()

	prepper.PreTransaction(executor.State[*substate.Substate]{
		Block: 5,
		Data:  &substate.Substate{InputAlloc: allocA},
	}, ctx)

	prepper.PreTransaction(executor.State[*substate.Substate]{
		Block: 7,
		Data:  &substate.Substate{InputAlloc: allocB},
	}, ctx)
}

func TestStatePrepper_DoesNotCrashOnMissingStateOrSubstate(t *testing.T) {
	ctrl := gomock.NewController(t)
	db := state.NewMockStateDB(ctrl)
	ctx := &executor.Context{State: db}

	prepper := MakeStateDbPrepper()
	prepper.PreTransaction(executor.State[*substate.Substate]{Block: 5}, nil)                             // misses both
	prepper.PreTransaction(executor.State[*substate.Substate]{Block: 5}, ctx)                             // misses the substate
	prepper.PreTransaction(executor.State[*substate.Substate]{Block: 5, Data: &substate.Substate{}}, nil) // misses the state
}