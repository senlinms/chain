package bc

import "errors"

var (
	errPosition              = errors.New("invalid source or destination position")
	errEntryType             = errors.New("invalid entry type")
	errBadTimeRange          = errors.New("bad time range")
	errEmptyResults          = errors.New("transaction has no results")
	errMismatchedAssetID     = errors.New("mismatched asset id")
	errMismatchedBlock       = errors.New("mismatched block")
	errMismatchedMerkleRoot  = errors.New("mismatched merkle root")
	errMismatchedReference   = errors.New("mismatched reference")
	errMismatchedValue       = errors.New("mismatched value")
	errMisorderedBlockHeight = errors.New("misordered block height")
	errMisorderedBlockTime   = errors.New("misordered block time")
	errNoPrevBlock           = errors.New("no previous block")
	errNoSource              = errors.New("no source for value")
	errNonemptyExtHash       = errors.New("non-empty extension hash")
	errOverflow              = errors.New("arithmetic overflow/underflow")
	errTxVersion             = errors.New("invalid transaction version")
	errUnbalanced            = errors.New("unbalanced")
	errUntimelyTransaction   = errors.New("block timestamp outside transaction time range")
	errVersionRegression     = errors.New("version regression")
	errWrongBlockchain       = errors.New("wrong blockchain")
	errZeroTime              = errors.New("timerange has one or two bounds set to zero")
)

type validationState struct {
	blockVersion   uint64
	initialBlockID Hash

	currentTx *TxEntries

	// Set this to the ID of an entry whenever validating an entry
	currentEntryID Hash

	// Must be defined when validating a valueSource
	sourcePosition uint64

	// Must be defined when validating a valueDestination
	destPosition uint64

	// The block timestamp
	timestampMS       uint64
	prevBlockHeader   *BlockHeaderEntry
	prevBlockHeaderID Hash
	blockTxs          []*TxEntries
	// xxx reachable entries?
}

type blockVMContext struct {
}

func newBlockVMContext(blockEntries *BlockEntries, prog []byte, args [][]byte) *blockVMContext {
	return &blockVMContext{}
}

type txVMContext struct {
}

func newTxVMContext(txEntries *TxEntries, entry Entry, prog Program, args [][]byte) *txVMContext {
	return &txVMContext{}
}
