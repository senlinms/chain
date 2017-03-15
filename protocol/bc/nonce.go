package bc

import (
	"chain/errors"
	"chain/protocol/vm"
)

// Nonce contains data used, among other things, for distinguishing
// otherwise-identical issuances (when used as those issuances'
// "anchors"). It satisfies the Entry interface.
type Nonce struct {
	body struct {
		Program   Program
		TimeRange Hash
		ExtHash   Hash
	}

	witness struct {
		Arguments [][]byte
		Anchored  Hash
	}

	// TimeRange contains (a pointer to) the manifested entry
	// corresponding to body.TimeRange.
	TimeRange *TimeRange

	// Anchored contains a pointer to the manifested entry corresponding
	// to witness.Anchored.
	Anchored Entry
}

func (Nonce) Type() string         { return "nonce1" }
func (n *Nonce) Body() interface{} { return n.body }

func (Nonce) Ordinal() int { return -1 }

// NewNonce creates a new Nonce.
func NewNonce(p Program, tr *TimeRange) *Nonce {
	n := new(Nonce)
	n.body.Program = p
	n.body.TimeRange = EntryID(tr)
	n.TimeRange = tr
	return n
}

func (n *Nonce) CheckValid(state *validationState) error {
	err := vm.Verify(newTxVMContext(state.currentTx, n, n.body.Program, n.witness.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking nonce program")
	}

	// xxx recursively validate the timerange?

	if n.TimeRange.body.MinTimeMS == 0 || n.TimeRange.body.MaxTimeMS == 0 {
		return errZeroTime
	}

	if state.currentTx.body.Version == 1 && (n.body.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
