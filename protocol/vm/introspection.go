package vm

import (
	"bytes"
	"fmt"
	"math"
)

func opCheckOutput(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(16)
	if err != nil {
		return err
	}

	code, err := vm.pop(true)
	if err != nil {
		return err
	}
	vmVersion, err := vm.popInt64(true)
	if err != nil {
		return err
	}
	if vmVersion < 0 {
		return ErrBadValue
	}
	assetID, err := vm.pop(true)
	if err != nil {
		return err
	}
	amount, err := vm.popInt64(true)
	if err != nil {
		return err
	}
	if amount < 0 {
		return ErrBadValue
	}
	refdatahash, err := vm.pop(true)
	if err != nil {
		return err
	}
	index, err := vm.popInt64(true)
	if err != nil {
		return err
	}
	if index < 0 {
		return ErrBadValue
	}
	if index > math.MaxUint32 {
		return ErrBadValue // xxx
	}

	// The following is per the discussion at
	// https://chainhq.slack.com/archives/txgraph/p1487964172000960
	if !vm.tx.DestIsMux(vm.inputIndex) {
		return ErrContext // xxx ?
	}

	isRetirement, destAssetID, destAmount, destData, destVMVersion, destCode, err := vm.tx.MuxDest(uint32(index))
	if err != nil {
		return err // xxx ?
	}

	someChecks := func(resAssetID []byte, resAmount uint64, resData []byte) bool {
		if !bytes.Equal(resAssetID, assetID) {
			return false
		}
		if resAmount != uint64(amount) {
			return false
		}
		if len(refdatahash) > 0 && !bytes.Equal(refdatahash, resData) {
			return false
		}
		return true
	}

	ok := someChecks(destAssetID, destAmount, destData)
	if !ok {
		return vm.pushBool(false, true)
	}

	if isRetirement {
		if vmVersion == 1 && len(code) > 0 && code[0] == byte(OP_FAIL) {
			// Special case alert! Old-style retirements were just outputs
			// with a control program beginning [FAIL]. New-style retirements
			// do not have control programs, but for compatibility we allow
			// CHECKOUTPUT to test for them by specifying a programming
			// beginnning with [FAIL].
			return vm.pushBool(true, true)
		}
		return vm.pushBool(false, true)
	}

	if destVMVersion != uint64(vmVersion) {
		return vm.pushBool(false, true)
	}
	if !bytes.Equal(destCode, code) {
		return vm.pushBool(false, true)
	}
	return vm.pushBool(true, true)
}

func opAsset(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	assetID, err := vm.tx.AssetID(vm.inputIndex)
	if err != nil {
		return ErrContext // xxx right?
	}

	return vm.push(assetID, true)
}

func opAmount(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	amount, err := vm.tx.Amount(vm.inputIndex)
	if err != nil {
		return err // xxx ?
	}

	return vm.pushInt64(int64(amount), true)
}

func opProgram(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	return vm.push(vm.mainprog, true)
}

func opMinTime(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	return vm.pushInt64(int64(vm.tx.MinTimeMS()), true)
}

func opMaxTime(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	maxTime := vm.tx.MaxTimeMS()
	if maxTime == 0 || maxTime > math.MaxInt64 {
		maxTime = uint64(math.MaxInt64)
	}

	return vm.pushInt64(int64(maxTime), true)
}

func opRefDataHash(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	data, err := vm.tx.InpData(vm.inputIndex)
	if err != nil {
		return err // xxx ?
	}

	return vm.push(data[:], true)
}

func opTxRefDataHash(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	return vm.push(vm.tx.TxData(), true)
}

func opIndex(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	return vm.pushInt64(int64(vm.inputIndex), true)
}

func opOutputID(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	spentOutputID, err := vm.tx.SpentOutputID(vm.inputIndex)
	if err != nil {
		return ErrContext // xxx ?
	}
	return vm.push(spentOutputID, true)
}

func opNonce(vm *virtualMachine) error {
	if vm.tx == nil {
		return ErrContext
	}

	err := vm.applyCost(1)
	if err != nil {
		return err
	}

	anchorID, err := vm.tx.AnchorID(vm.inputIndex)
	if err != nil {
		return err // xxx ?
	}

	return vm.push(anchorID, true)
}

func opNextProgram(vm *virtualMachine) error {
	if vm.block == nil {
		return ErrContext
	}
	err := vm.applyCost(1)
	if err != nil {
		return err
	}
	return vm.push(vm.block.NextConsensusProgram(), true)
}

func opBlockTime(vm *virtualMachine) error {
	if vm.block == nil {
		return ErrContext
	}
	err := vm.applyCost(1)
	if err != nil {
		return err
	}
	if vm.block.TimestampMS() > math.MaxInt64 {
		return fmt.Errorf("block timestamp out of range")
	}
	return vm.pushInt64(int64(vm.block.TimestampMS()), true)
}
