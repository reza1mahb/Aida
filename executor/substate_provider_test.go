// Copyright 2024 Fantom Foundation
// This file is part of Aida Testing Infrastructure for Sonic
//
// Aida is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Aida is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Aida. If not, see <http://www.gnu.org/licenses/>.

package executor

//go:generate mockgen -source substate_provider_test.go -destination substate_provider_test_mocks.go -package executor

import (
	"errors"
	"math/big"
	"testing"

	"github.com/Fantom-foundation/Aida/txcontext"
	"github.com/Fantom-foundation/Aida/utils"
	substate "github.com/Fantom-foundation/Substate"
	"go.uber.org/mock/gomock"
)

func TestSubstateProvider_OpeningANonExistingDbResultsInAnError(t *testing.T) {
	cfg := utils.Config{}
	cfg.AidaDb = t.TempDir()
	// Important: the following code does not panic.
	_, err := OpenSubstateDb(&cfg, nil)
	if err == nil {
		t.Errorf("attempting to open a non-existing substate DB should fail")
	}
}

func TestSubstateProvider_IterateOverExistingDb(t *testing.T) {
	ctrl := gomock.NewController(t)
	consumer := NewMockTxConsumer(ctrl)

	// Prepare a directory containing some substate data.
	path := t.TempDir()
	if err := createSubstateDb(path); err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	// Open the substate data for reading.
	provider, err := openSubstateDb(path)
	if err != nil {
		t.Fatalf("failed to open substate DB: %v", err)
	}
	defer provider.Close()

	gomock.InOrder(
		consumer.EXPECT().Consume(10, 7, gomock.Any()),
		consumer.EXPECT().Consume(10, 9, gomock.Any()),
		consumer.EXPECT().Consume(12, 5, gomock.Any()),
	)

	if err := provider.Run(0, 20, toSubstateConsumer(consumer)); err != nil {
		t.Fatalf("failed to iterate through states: %v", err)
	}
}

func TestSubstateProvider_LowerBoundIsInclusive(t *testing.T) {
	ctrl := gomock.NewController(t)
	consumer := NewMockTxConsumer(ctrl)

	// Prepare a directory containing some substate data.
	path := t.TempDir()
	if err := createSubstateDb(path); err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	// Open the substate data for reading.
	provider, err := openSubstateDb(path)
	if err != nil {
		t.Fatalf("failed to open substate DB: %v", err)
	}
	defer provider.Close()

	gomock.InOrder(
		consumer.EXPECT().Consume(10, 7, gomock.Any()),
		consumer.EXPECT().Consume(10, 9, gomock.Any()),
		consumer.EXPECT().Consume(12, 5, gomock.Any()),
	)

	if err := provider.Run(10, 20, toSubstateConsumer(consumer)); err != nil {
		t.Fatalf("failed to iterate through states: %v", err)
	}
}

func TestSubstateProvider_UpperBoundIsExclusive(t *testing.T) {
	ctrl := gomock.NewController(t)
	consumer := NewMockTxConsumer(ctrl)

	// Prepare a directory containing some substate data.
	path := t.TempDir()
	if err := createSubstateDb(path); err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	// Open the substate data for reading.
	provider, err := openSubstateDb(path)
	if err != nil {
		t.Fatalf("failed to open substate DB: %v", err)
	}
	defer provider.Close()

	gomock.InOrder(
		consumer.EXPECT().Consume(10, 7, gomock.Any()),
		consumer.EXPECT().Consume(10, 9, gomock.Any()),
	)

	if err := provider.Run(10, 12, toSubstateConsumer(consumer)); err != nil {
		t.Fatalf("failed to iterate through states: %v", err)
	}
}

func TestSubstateProvider_RangeCanBeEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	consumer := NewMockTxConsumer(ctrl)

	// Prepare a directory containing some substate data.
	path := t.TempDir()
	if err := createSubstateDb(path); err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	// Open the substate data for reading.
	provider, err := openSubstateDb(path)
	if err != nil {
		t.Fatalf("failed to open substate DB: %v", err)
	}
	defer provider.Close()

	if err := provider.Run(5, 10, toSubstateConsumer(consumer)); err != nil {
		t.Fatalf("failed to iterate through states: %v", err)
	}
}

func TestSubstateProvider_IterationCanBeAbortedByConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	consumer := NewMockTxConsumer(ctrl)

	// Prepare a directory containing some substate data.
	path := t.TempDir()
	if err := createSubstateDb(path); err != nil {
		t.Fatalf("failed to setup test DB: %v", err)
	}

	// Open the substate data for reading.
	provider, err := openSubstateDb(path)
	if err != nil {
		t.Fatalf("failed to open substate DB: %v", err)
	}
	defer provider.Close()

	stop := errors.New("stop!")
	gomock.InOrder(
		consumer.EXPECT().Consume(10, 7, gomock.Any()),
		consumer.EXPECT().Consume(10, 9, gomock.Any()).Return(stop),
	)

	if got, want := provider.Run(10, 20, toSubstateConsumer(consumer)), stop; !errors.Is(got, want) {
		t.Errorf("provider run did not finish with expected exception, wanted %d, got %d", want, got)
	}
}

func openSubstateDb(path string) (Provider[txcontext.TxContext], error) {
	cfg := utils.Config{}
	cfg.AidaDb = path
	cfg.Workers = 1
	return OpenSubstateDb(&cfg, nil)
}

func createSubstateDb(path string) error {
	substate.SetSubstateDb(path)
	substate.OpenSubstateDB()

	state := substate.Substate{
		Env: &substate.SubstateEnv{},
		Message: &substate.SubstateMessage{
			Value: big.NewInt(12),
		},
		InputAlloc:  substate.SubstateAlloc{},
		OutputAlloc: substate.SubstateAlloc{},
		Result:      &substate.SubstateResult{},
	}

	substate.PutSubstate(10, 7, &state)
	substate.PutSubstate(10, 9, &state)
	substate.PutSubstate(12, 5, &state)

	substate.CloseSubstateDB()
	return nil
}
