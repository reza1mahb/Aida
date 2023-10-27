package main

import (
	"math/big"
	"testing"

	"github.com/Fantom-foundation/Aida/executor"
	"github.com/Fantom-foundation/Aida/utils"
	substate "github.com/Fantom-foundation/Substate"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/mock/gomock"
)

var testingAddress = common.Address{1}

func TestSdbRecord_AllDbEventsAreIssuedInOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	provider := executor.NewMockProvider[*substate.Substate](ctrl)
	processor := executor.NewMockProcessor[*substate.Substate](ctrl)
	ext := executor.NewMockExtension[*substate.Substate](ctrl)
	path := t.TempDir() + "test_trace"
	cfg := &utils.Config{
		First:            10,
		Last:             11,
		ChainID:          utils.MainnetChainID,
		SkipPriming:      true,
		TraceFile:        path,
		SyncPeriodLength: 1,
	}

	provider.EXPECT().
		Run(10, 12, gomock.Any()).
		DoAndReturn(func(from int, to int, consumer executor.Consumer[*substate.Substate]) error {
			for i := from; i < to; i++ {
				consumer(executor.TransactionInfo[*substate.Substate]{Block: i, Transaction: 3, Data: emptyTx})
				consumer(executor.TransactionInfo[*substate.Substate]{Block: i, Transaction: utils.PseudoTx, Data: emptyTx})
			}
			return nil
		})

	// All transactions are processed in order
	gomock.InOrder(
		ext.EXPECT().PreRun(executor.AtBlock[*substate.Substate](10), gomock.Any()),

		// block 10
		ext.EXPECT().PreBlock(executor.AtBlock[*substate.Substate](10), gomock.Any()),

		ext.EXPECT().PreTransaction(executor.AtTransaction[*substate.Substate](10, 3), gomock.Any()),
		processor.EXPECT().Process(executor.AtTransaction[*substate.Substate](10, 3), gomock.Any()),
		ext.EXPECT().PostTransaction(executor.AtTransaction[*substate.Substate](10, 3), gomock.Any()),

		ext.EXPECT().PreTransaction(executor.AtTransaction[*substate.Substate](10, utils.PseudoTx), gomock.Any()),
		processor.EXPECT().Process(executor.AtTransaction[*substate.Substate](10, utils.PseudoTx), gomock.Any()),
		ext.EXPECT().PostTransaction(executor.AtTransaction[*substate.Substate](10, utils.PseudoTx), gomock.Any()),

		ext.EXPECT().PostBlock(executor.AtBlock[*substate.Substate](10), gomock.Any()),

		// block 11
		ext.EXPECT().PreBlock(executor.AtBlock[*substate.Substate](11), gomock.Any()),

		ext.EXPECT().PreTransaction(executor.AtTransaction[*substate.Substate](11, 3), gomock.Any()),
		processor.EXPECT().Process(executor.AtTransaction[*substate.Substate](11, 3), gomock.Any()),
		ext.EXPECT().PostTransaction(executor.AtTransaction[*substate.Substate](11, 3), gomock.Any()),

		ext.EXPECT().PreTransaction(executor.AtTransaction[*substate.Substate](11, utils.PseudoTx), gomock.Any()),
		processor.EXPECT().Process(executor.AtTransaction[*substate.Substate](11, utils.PseudoTx), gomock.Any()),
		ext.EXPECT().PostTransaction(executor.AtTransaction[*substate.Substate](11, utils.PseudoTx), gomock.Any()),

		ext.EXPECT().PostBlock(executor.AtBlock[*substate.Substate](11), gomock.Any()),

		ext.EXPECT().PostRun(executor.AtBlock[*substate.Substate](12), gomock.Any(), nil),
	)

	if err := record(cfg, provider, processor, []executor.Extension[*substate.Substate]{ext}); err != nil {
		t.Errorf("record failed: %v", err)
	}
}

// emptyTx is a dummy substate that will be processed without crashing.
var emptyTx = &substate.Substate{
	Env: &substate.SubstateEnv{},
	Message: &substate.SubstateMessage{
		GasPrice: big.NewInt(12),
	},
	Result: &substate.SubstateResult{
		GasUsed: 1,
	},
}

// testTx is a dummy substate used for testing validation.
var testTx = &substate.Substate{
	InputAlloc: substate.SubstateAlloc{testingAddress: substate.NewSubstateAccount(1, new(big.Int).SetUint64(1), []byte{})},
	Env:        &substate.SubstateEnv{},
	Message: &substate.SubstateMessage{
		GasPrice: big.NewInt(12),
	},
	Result: &substate.SubstateResult{
		GasUsed: 1,
	},
}