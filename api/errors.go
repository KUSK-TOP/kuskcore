package api

import (
	"context"

	"kuskcore/account"
	"kuskcore/asset"
	"kuskcore/blockchain/pseudohsm"
	"kuskcore/blockchain/rpc"
	"kuskcore/blockchain/signers"
	"kuskcore/blockchain/txbuilder"
	"kuskcore/contract"
	"kuskcore/errors"
	"kuskcore/net/http/httperror"
	"kuskcore/net/http/httpjson"
	"kuskcore/protocol/validation"
	"kuskcore/protocol/vm"
)

var (
	// ErrDefault is default Kusk API Error
	ErrDefault = errors.New("Kusk API Error")
)

func isTemporary(info httperror.Info, err error) bool {
	switch info.ChainCode {
	case "KUSK000": // internal server error
		return true
	case "KUSK001": // request timed out
		return true
	case "KUSK761": // outputs currently reserved
		return true
	case "KUSK706": // 1 or more action errors
		errs := errors.Data(err)["actions"].([]httperror.Response)
		temp := true
		for _, actionErr := range errs {
			temp = temp && isTemporary(actionErr.Info, nil)
		}
		return temp
	default:
		return false
	}
}

var respErrFormatter = map[error]httperror.Info{
	ErrDefault: {500, "KUSK000", "Kusk API Error"},

	// Signers error namespace (2xx)
	signers.ErrBadQuorum: {400, "KUSK200", "Quorum must be greater than or equal to 1, and must be less than or equal to the length of xpubs"},
	signers.ErrBadXPub:   {400, "KUSK201", "Invalid xpub format"},
	signers.ErrNoXPubs:   {400, "KUSK202", "At least one xpub is required"},
	signers.ErrDupeXPub:  {400, "KUSK203", "Root XPubs cannot contain the same key more than once"},

	// Contract error namespace (3xx)
	contract.ErrContractDuplicated: {400, "KUSK302", "Contract is duplicated"},
	contract.ErrContractNotFound:   {400, "KUSK303", "Contract not found"},

	// Transaction error namespace (7xx)
	// Build transaction error namespace (70x ~ 72x)
	account.ErrInsufficient:         {400, "KUSK700", "Funds of account are insufficient"},
	account.ErrImmature:             {400, "KUSK701", "Available funds of account are immature"},
	account.ErrReserved:             {400, "KUSK702", "Available UTXOs of account have been reserved"},
	account.ErrMatchUTXO:            {400, "KUSK703", "UTXO with given hash not found"},
	ErrBadActionType:                {400, "KUSK704", "Invalid action type"},
	ErrBadAction:                    {400, "KUSK705", "Invalid action object"},
	ErrBadActionConstruction:        {400, "KUSK706", "Invalid action construction"},
	txbuilder.ErrMissingFields:      {400, "KUSK707", "One or more fields are missing"},
	txbuilder.ErrBadAmount:          {400, "KUSK708", "Invalid asset amount"},
	account.ErrFindAccount:          {400, "KUSK709", "Account not found"},
	asset.ErrFindAsset:              {400, "KUSK710", "Asset not found"},
	txbuilder.ErrBadContractArgType: {400, "KUSK711", "Invalid contract argument type"},
	txbuilder.ErrOrphanTx:           {400, "KUSK712", "Transaction input UTXO not found"},
	txbuilder.ErrExtTxFee:           {400, "KUSK713", "Transaction fee exceeded max limit"},
	txbuilder.ErrNoGasInput:         {400, "KUSK714", "Transaction has no gas input"},

	// Submit transaction error namespace (73x ~ 79x)
	// Validation error (73x ~ 75x)
	validation.ErrTxVersion:                 {400, "KUSK730", "Invalid transaction version"},
	validation.ErrWrongTransactionSize:      {400, "KUSK731", "Invalid transaction size"},
	validation.ErrBadTimeRange:              {400, "KUSK732", "Invalid transaction time range"},
	validation.ErrNotStandardTx:             {400, "KUSK733", "Not standard transaction"},
	validation.ErrWrongCoinbaseTransaction:  {400, "KUSK734", "Invalid coinbase transaction"},
	validation.ErrWrongCoinbaseAsset:        {400, "KUSK735", "Invalid coinbase assetID"},
	validation.ErrCoinbaseArbitraryOversize: {400, "KUSK736", "Invalid coinbase arbitrary size"},
	validation.ErrEmptyResults:              {400, "KUSK737", "No results in the transaction"},
	validation.ErrMismatchedAssetID:         {400, "KUSK738", "Mismatched assetID"},
	validation.ErrMismatchedPosition:        {400, "KUSK739", "Mismatched value source/dest position"},
	validation.ErrMismatchedReference:       {400, "KUSK740", "Mismatched reference"},
	validation.ErrMismatchedValue:           {400, "KUSK741", "Mismatched value"},
	validation.ErrMissingField:              {400, "KUSK742", "Missing required field"},
	validation.ErrNoSource:                  {400, "KUSK743", "No source for value"},
	validation.ErrOverflow:                  {400, "KUSK744", "Arithmetic overflow/underflow"},
	validation.ErrPosition:                  {400, "KUSK745", "Invalid source or destination position"},
	validation.ErrUnbalanced:                {400, "KUSK746", "Unbalanced asset amount between input and output"},
	validation.ErrOverGasCredit:             {400, "KUSK747", "Gas credit has been spent"},
	validation.ErrGasCalculate:              {400, "KUSK748", "Gas usage calculate got a math error"},

	// VM error (76x ~ 78x)
	vm.ErrAltStackUnderflow:  {400, "KUSK760", "Alt stack underflow"},
	vm.ErrBadValue:           {400, "KUSK761", "Bad value"},
	vm.ErrContext:            {400, "KUSK762", "Wrong context"},
	vm.ErrDataStackUnderflow: {400, "KUSK763", "Data stack underflow"},
	vm.ErrDisallowedOpcode:   {400, "KUSK764", "Disallowed opcode"},
	vm.ErrDivZero:            {400, "KUSK765", "Division by zero"},
	vm.ErrFalseVMResult:      {400, "KUSK766", "False result for executing VM"},
	vm.ErrLongProgram:        {400, "KUSK767", "Program size exceeds max int32"},
	vm.ErrRange:              {400, "KUSK768", "Arithmetic range error"},
	vm.ErrReturn:             {400, "KUSK769", "RETURN executed"},
	vm.ErrRunLimitExceeded:   {400, "KUSK770", "Run limit exceeded because the KUSK Fee is insufficient"},
	vm.ErrShortProgram:       {400, "KUSK771", "Unexpected end of program"},
	vm.ErrToken:              {400, "KUSK772", "Unrecognized token"},
	vm.ErrUnexpected:         {400, "KUSK773", "Unexpected error"},
	vm.ErrUnsupportedVM:      {400, "KUSK774", "Unsupported VM because the version of VM is mismatched"},
	vm.ErrVerifyFailed:       {400, "KUSK775", "VERIFY failed"},

	// Mock HSM error namespace (8xx)
	pseudohsm.ErrDuplicateKeyAlias: {400, "KUSK800", "Key Alias already exists"},
	pseudohsm.ErrLoadKey:           {400, "KUSK801", "Key not found or wrong password"},
	pseudohsm.ErrDecrypt:           {400, "KUSK802", "Could not decrypt key with given passphrase"},
}

// Map error values to standard kusk error codes. Missing entries
// will map to internalErrInfo.
//
// TODO(jackson): Share one error table across Chain
// products/services so that errors are consistent.
var errorFormatter = httperror.Formatter{
	Default:     httperror.Info{500, "KUSK000", "Kusk API Error"},
	IsTemporary: isTemporary,
	Errors: map[error]httperror.Info{
		// General error namespace (0xx)
		context.DeadlineExceeded: {408, "KUSK001", "Request timed out"},
		httpjson.ErrBadRequest:   {400, "KUSK002", "Invalid request body"},
		rpc.ErrWrongNetwork:      {502, "KUSK103", "A peer core is operating on a different blockchain network"},

		//accesstoken authz err namespace (86x)
		errNotAuthenticated: {401, "KUSK860", "Request could not be authenticated"},
	},
}
