package apireplay

import (
	"encoding/json"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/Fantom-foundation/Aida/iterator"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/status-im/keycard-go/hexutils"
)

const (
	hexOne  = "0x1"
	hexZero = "0x0"

	// 32 bytes returned by EVM as result one
	longHexOne = "0x0000000000000000000000000000000000000000000000000000000000000001"

	// 32 bytes returned by EVM as a zero result
	longHexZero = "0x0000000000000000000000000000000000000000000000000000000000000000"
)

// Test_compareBalanceOK tests compare func for getBalance method
// It expects no error since results are same
func Test_compareBalanceOK(t *testing.T) {
	bigRes, _ := new(big.Int).SetString("1", 16)
	data := &OutData{
		Method: "eth_getBalance",
		Recorded: &RecordedData{
			Result: []byte(hexOne),
		},
		StateDB: &StateDBData{
			Result: bigRes,
		},
	}

	err := compareBalance(data)
	if err != nil {
		t.Errorf("error must be nil; err: %v", err)
	}

}

// Test_compareBalanceErrorNoMatchingResult tests compare func for getBalance method
// It expects an error of no matching results since results are different
func Test_compareBalanceErrorNoMatchingResult(t *testing.T) {
	bigRes, _ := new(big.Int).SetString("1", 16)
	data := &OutData{
		Method: "eth_getBalance",
		Recorded: &RecordedData{
			Result: []byte(hexZero),
		},
		StateDB: &StateDBData{
			Result: bigRes,
		},
	}

	err := compareBalance(data)
	if err == nil {
		t.Errorf("error must not be nil; err: %v", err)
		return
	}

	if err.typ != noMatchingResult {
		t.Errorf("error must be type 'noMatchingResult'; err: %v", err)
	}

}

// Test_compareTransactionCountOK tests compare func for getTransactionCount method
// It expects no error since results are same
func Test_compareTransactionCountOK(t *testing.T) {

	data := &OutData{
		Method: "eth_getTransactionCount",
		Recorded: &RecordedData{
			Result: []byte(hexOne),
		},
		StateDB: &StateDBData{
			Result: uint64(1),
		},
	}

	err := compareTransactionCount(data)

	if err != nil {
		t.Errorf("error must be nil; err: %v", err)
	}

}

// Test_compareTransactionCountErrorNoMatchingResult tests compare func for getTransactionCount method
// It expects an error of no matching results since results are different
func Test_compareTransactionCountErrorNoMatchingResult(t *testing.T) {
	data := &OutData{
		Method: "eth_getTransactionCount",
		Recorded: &RecordedData{
			Result: []byte(hexZero),
		},
		StateDB: &StateDBData{
			Result: uint64(1),
		},
	}

	err := compareTransactionCount(data)
	if err == nil {
		t.Errorf("error must not be nil; err: %v", err)
		return
	}

	if err.typ != noMatchingResult {
		t.Errorf("error must be type 'noMatchingResult'; err: %v", err)
	}

}

// Test_compareCallOK tests compare func for call method
// It expects no error since results are same
func Test_compareCallOK(t *testing.T) {

	rec, _ := json.Marshal(strings.TrimPrefix(longHexOne, "0x"))
	data := &OutData{
		Method: "eth_call",
		Recorded: &RecordedData{
			Result: rec,
		},
		StateDB: &StateDBData{
			Result: hexutils.HexToBytes(strings.TrimPrefix(longHexOne, "0x")),
		},
	}

	err := compareCall(data)
	if err != nil {
		t.Errorf("error must be nil; err: %v", err)
	}
}

// Test_compareCallErrorNoMatchingResult tests compare func for call method
// It expects an error of no matching results since results are different
func Test_compareCallErrorNoMatchingResult(t *testing.T) {
	rec, _ := json.Marshal(longHexOne)
	data := &OutData{
		Method: "eth_call",
		Recorded: &RecordedData{
			Result: rec,
		},
		StateDB: &StateDBData{
			Result: hexutils.HexToBytes(strings.TrimPrefix(longHexZero, "0x")),
		},
	}

	err := compareCall(data)
	if err == nil {
		t.Errorf("error must not be nil; err: %v", err)
		return
	}

	if err.typ != noMatchingResult {
		t.Errorf("error must be type 'noMatchingResult'; err: %v", err)
	}

}

// Test_compareCallErrorExpectedResultGotErr tests compare func for call method
// It expects an error of "expected valid result, got error" since recorded data is a valid result but EVM returns error
func Test_compareCallErrorExpectedResultGotErr(t *testing.T) {
	data := &OutData{
		Method: "eth_call",
		Recorded: &RecordedData{
			Result: []byte(hexOne),
		},
		StateDB: &StateDBData{
			Error: errors.New("error"),
		},
	}

	err := compareCall(data)
	if err == nil {
		t.Errorf("error must not be nil; err: %v", err)
		return
	}

	if err.typ != expectedResultGotError {
		t.Errorf("error must be type 'expectedResultGotError'; err: %v", err)
	}

}

// Test_compareCallErrorExpectedErrGotResult tests compare func for call method
// It expects an error of "expected error, got valid result" since recorded data is an error but EVM returns valid result
func Test_compareCallErrorExpectedErrGotResult(t *testing.T) {
	data := &OutData{
		Method: "eth_call",
		Recorded: &RecordedData{
			Error: &iterator.ErrorMessage{
				Code:    -32000,
				Message: "error",
			},
		},
		StateDB: &StateDBData{
			Result: hexutils.HexToBytes(strings.TrimPrefix(longHexZero, "0x")),
		},
	}

	err := compareCall(data)
	if err == nil {
		t.Errorf("error must not be null")
		return
	}

	if err.typ != expectedErrorGotResult {
		t.Errorf("error must be type 'expectedErrorGotResult'; err: %v", err)
	}

}

// Test_compareEstimateGasOK tests compare func for estimateGas method
// It expects no error since results are same
func Test_compareEstimateGasOK(t *testing.T) {
	data := &OutData{
		Method: "eth_getTransactionCount",
		Recorded: &RecordedData{
			Result: []byte(hexOne),
		},
		StateDB: &StateDBData{
			Result: hexutil.Uint64(1),
		},
	}

	err := compareEstimateGas(data)
	if err != nil {
		t.Errorf("error must be nil; err: %v", err)
	}
}

// Test_compareEstimateGasErrorNoMatchingResult tests compare func for estimateGas method
// It expects an error of no matching results since results are different
func Test_compareEstimateGasErrorNoMatchingResult(t *testing.T) {
	data := &OutData{
		Method: "eth_estimateGas",
		Recorded: &RecordedData{
			Result: []byte(hexOne),
		},
		StateDB: &StateDBData{
			Result: hexutil.Uint64(0),
		},
	}

	err := compareEstimateGas(data)
	if err == nil {
		t.Errorf("error must not be null")
		return
	}

	if err.typ != noMatchingResult {
		t.Errorf("error must be type 'noMatchingResult'; err: %v", err)
	}

}

// Test_compareEstimateGasErrorExpectedResultGotErr tests compare func for estimateGas method
// It expects an error of "expected valid result, got error" since recorded data is a valid result but EVM returns error
func Test_compareEstimateGasErrorExpectedResultGotErr(t *testing.T) {
	data := &OutData{
		Method: "eth_estimateGas",
		Recorded: &RecordedData{
			Result: []byte(hexOne),
		},
		StateDB: &StateDBData{
			Error: errors.New("error"),
		},
	}

	err := compareEstimateGas(data)
	if err == nil {
		t.Errorf("error must be nil; err: %v", err)
		return
	}

	if err.typ != expectedResultGotError {
		t.Errorf("error must be type 'expectedResultGotError'; err: %v", err)
	}
}

// Test_compareEstimateGasErrorExpectedErrGotResult tests compare func for estimateGas method
// It expects an error of "expected error, got valid result" since recorded data is an error but EVM returns valid result
func Test_compareEstimateGasErrorExpectedErrGotResult(t *testing.T) {
	data := &OutData{
		Method: "eth_estimateGas",
		Recorded: &RecordedData{
			Error: &iterator.ErrorMessage{
				Code:    1000,
				Message: "error",
			},
		},
		StateDB: &StateDBData{
			Result: hexutil.Uint64(0),
		},
	}

	err := compareEstimateGas(data)
	if err == nil {
		t.Errorf("error must not be null")
		return
	}

	if err.typ != expectedErrorGotResult {
		t.Errorf("error must be type 'expectedErrorGotResult'; err: %v", err)
	}

}