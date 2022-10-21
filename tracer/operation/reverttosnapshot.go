package operation

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/Fantom-foundation/Aida/tracer/dict"
	"github.com/Fantom-foundation/Aida/tracer/state"
)

// Revert-to-snapshot operation's data structure with returned snapshot id
type RevertToSnapshot struct {
	SnapshotID uint16
}

// Return the revert-to-snapshot operation identifier.
func (op *RevertToSnapshot) GetOpId() byte {
	return RevertToSnapshotID
}

// Create a new revert-to-snapshot operation.
func NewRevertToSnapshot(SnapshotID int) *RevertToSnapshot {
	return &RevertToSnapshot{SnapshotID: uint16(SnapshotID)}
}

// Read a revert-to-snapshot operation from file.
func ReadRevertToSnapshot(file *os.File) (Operation, error) {
	data := new(RevertToSnapshot)
	err := binary.Read(file, binary.LittleEndian, data)
	return data, err
}

// Write the revert-to-snapshot operation to file.
func (op *RevertToSnapshot) Write(f *os.File) error {
	err := binary.Write(f, binary.LittleEndian, *op)
	return err
}

// Execute the revert-to-snapshot operation.
func (op *RevertToSnapshot) Execute(db state.StateDB, ctx *dict.DictionaryContext) {
	id := ctx.GetSnapshot(op.SnapshotID)
	db.RevertToSnapshot(int(id))
}

// Print a debug message for revert-to-snapshot operation.
func (op *RevertToSnapshot) Debug(ctx *dict.DictionaryContext) {
	fmt.Printf("\tsnapshot id: %v\n", op.SnapshotID)
}
