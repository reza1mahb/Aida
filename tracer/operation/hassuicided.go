package operation

import (
	"encoding/binary"
	"fmt"
	"github.com/Fantom-foundation/Aida/tracer/dict"
	"github.com/Fantom-foundation/Aida/tracer/state"
	"io"
	"time"
)

// HasSuicided data structure
type HasSuicided struct {
	ContractIndex uint32 // encoded contract address
}

// GetId returns the HasSuicided operation identifier.
func (op *HasSuicided) GetId() byte {
	return HasSuicidedID
}

// NewHasSuicided creates a new HasSuicided operation.
func NewHasSuicided(cIdx uint32) *HasSuicided {
	return &HasSuicided{ContractIndex: cIdx}
}

// ReadHasSuicided reads a HasSuicided operation from a file.
func ReadHasSuicided(file io.Reader) (Operation, error) {
	data := new(HasSuicided)
	err := binary.Read(file, binary.LittleEndian, data)
	return data, err
}

// Write the HasSuicided operation to a file.
func (op *HasSuicided) Write(f io.Writer) error {
	err := binary.Write(f, binary.LittleEndian, *op)
	return err
}

// Execute the HasSuicided operation.
func (op *HasSuicided) Execute(db state.StateDB, ctx *dict.DictionaryContext) time.Duration {
	contract := ctx.DecodeContract(op.ContractIndex)
	start := time.Now()
	db.HasSuicided(contract)
	return time.Since(start)
}

// Debug prints a debug message for the HasSuicided operation.
func (op *HasSuicided) Debug(ctx *dict.DictionaryContext) {
	fmt.Printf("\t%s: %s\n", operationLabels[HasSuicidedID], ctx.DecodeContract(op.ContractIndex))
}
