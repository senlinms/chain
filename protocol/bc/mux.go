package bc

import (
	"chain/errors"
	"chain/math/checked"
	"chain/protocol/vm"
)

// Mux splits and combines value from one or more source entries,
// making it available to one or more destination entries. It
// satisfies the Entry interface.
type Mux struct {
	body struct {
		Sources []valueSource // issuances, spends, and muxes
		Program Program
		ExtHash Hash
	}

	witness struct {
		Destinations []ValueDestination // outputs, retirements, and muxes
		Arguments    [][]byte
	}
}

func (Mux) Type() string         { return "mux1" }
func (m *Mux) Body() interface{} { return m.body }

func (Mux) Ordinal() int { return -1 }

func (m *Mux) Destinations() []ValueDestination {
	return m.witness.Destinations
}

// NewMux creates a new Mux.
func NewMux(sources []valueSource, program Program) *Mux {
	m := new(Mux)
	m.body.Sources = sources
	m.body.Program = program
	return m
}

func (mux *Mux) CheckValid(state *validationState) error {
	err := vm.Verify(newTxVMContext(state.currentTx, mux, mux.body.Program, mux.witness.Arguments))
	if err != nil {
		return errors.Wrap(err, "checking mux program")
	}

	for i, src := range mux.body.Sources {
		srcState := *state
		srcState.sourcePosition = uint64(i)
		err := src.CheckValid(&srcState)
		if err != nil {
			return errors.Wrapf(err, "checking mux source %d", i)
		}
	}

	for i, dest := range mux.witness.Destinations {
		destState := *state
		destState.destPosition = uint64(i)
		err := dest.CheckValid(&destState)
		if err != nil {
			return errors.Wrapf(err, "checking mux destination %d", i)
		}
	}

	parity := make(map[AssetID]int64)
	for i, src := range mux.body.Sources {
		sum, ok := checked.AddInt64(parity[src.Value.AssetID], int64(src.Value.Amount))
		if !ok {
			return errors.WithDetailf(errOverflow, "adding %d units of asset %x from mux source %d to total %d overflows int64", src.Value.Amount, src.Value.AssetID[:], i, parity[src.Value.AssetID])
		}
		parity[src.Value.AssetID] = sum
	}

	for i, dest := range mux.witness.Destinations {
		sum, ok := parity[dest.Value.AssetID]
		if !ok {
			return errors.WithDetailf(errNoSource, "mux destination %d, asset %x, has no corresponding source", i, dest.Value.AssetID[:])
		}

		diff, ok := checked.SubInt64(sum, int64(dest.Value.Amount))
		if !ok {
			return errors.WithDetailf(errOverflow, "subtracting %d units of asset %x from mux destination %d from total %d underflows int64", dest.Value.Amount, dest.Value.AssetID[:], i, sum)
		}
		parity[dest.Value.AssetID] = diff
	}

	for assetID, amount := range parity {
		if amount != 0 {
			return errors.WithDetailf(errUnbalanced, "asset %x sources - destinations = %d (should be 0)", assetID[:], amount)
		}
	}

	if state.currentTx.body.Version == 1 && (mux.body.ExtHash != Hash{}) {
		return errNonemptyExtHash
	}

	return nil
}
