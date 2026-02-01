// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FlemingAnchorMetaData contains all meta data concerning the FlemingAnchor contract.
var FlemingAnchorMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"initialAnchorer\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchor\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"anchorCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"anchorer\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"anchors\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchAnchor\",\"inputs\":[{\"name\":\"roots\",\"type\":\"bytes32[]\",\"internalType\":\"bytes32[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAnchorAge\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"age\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAnchorTimestamp\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getContractInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"version\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isAnchored\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceAnchorer\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAnchorer\",\"inputs\":[{\"name\":\"newAnchorer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AnchorerUpdated\",\"inputs\":[{\"name\":\"previousAnchorer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newAnchorer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RootAnchored\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"blockNumber\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"anchorer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyAnchored\",\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"Unauthorized\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroRoot\",\"inputs\":[]}]",
}

// FlemingAnchorABI is the input ABI used to generate the binding from.
// Deprecated: Use FlemingAnchorMetaData.ABI instead.
var FlemingAnchorABI = FlemingAnchorMetaData.ABI

// FlemingAnchor is an auto generated Go binding around an Ethereum contract.
type FlemingAnchor struct {
	FlemingAnchorCaller     // Read-only binding to the contract
	FlemingAnchorTransactor // Write-only binding to the contract
	FlemingAnchorFilterer   // Log filterer for contract events
}

// FlemingAnchorCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlemingAnchorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlemingAnchorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlemingAnchorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlemingAnchorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlemingAnchorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlemingAnchorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlemingAnchorSession struct {
	Contract     *FlemingAnchor    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlemingAnchorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlemingAnchorCallerSession struct {
	Contract *FlemingAnchorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// FlemingAnchorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlemingAnchorTransactorSession struct {
	Contract     *FlemingAnchorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// FlemingAnchorRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlemingAnchorRaw struct {
	Contract *FlemingAnchor // Generic contract binding to access the raw methods on
}

// FlemingAnchorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlemingAnchorCallerRaw struct {
	Contract *FlemingAnchorCaller // Generic read-only contract binding to access the raw methods on
}

// FlemingAnchorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlemingAnchorTransactorRaw struct {
	Contract *FlemingAnchorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlemingAnchor creates a new instance of FlemingAnchor, bound to a specific deployed contract.
func NewFlemingAnchor(address common.Address, backend bind.ContractBackend) (*FlemingAnchor, error) {
	contract, err := bindFlemingAnchor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FlemingAnchor{FlemingAnchorCaller: FlemingAnchorCaller{contract: contract}, FlemingAnchorTransactor: FlemingAnchorTransactor{contract: contract}, FlemingAnchorFilterer: FlemingAnchorFilterer{contract: contract}}, nil
}

// NewFlemingAnchorCaller creates a new read-only instance of FlemingAnchor, bound to a specific deployed contract.
func NewFlemingAnchorCaller(address common.Address, caller bind.ContractCaller) (*FlemingAnchorCaller, error) {
	contract, err := bindFlemingAnchor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlemingAnchorCaller{contract: contract}, nil
}

// NewFlemingAnchorTransactor creates a new write-only instance of FlemingAnchor, bound to a specific deployed contract.
func NewFlemingAnchorTransactor(address common.Address, transactor bind.ContractTransactor) (*FlemingAnchorTransactor, error) {
	contract, err := bindFlemingAnchor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlemingAnchorTransactor{contract: contract}, nil
}

// NewFlemingAnchorFilterer creates a new log filterer instance of FlemingAnchor, bound to a specific deployed contract.
func NewFlemingAnchorFilterer(address common.Address, filterer bind.ContractFilterer) (*FlemingAnchorFilterer, error) {
	contract, err := bindFlemingAnchor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlemingAnchorFilterer{contract: contract}, nil
}

// bindFlemingAnchor binds a generic wrapper to an already deployed contract.
func bindFlemingAnchor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FlemingAnchorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlemingAnchor *FlemingAnchorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlemingAnchor.Contract.FlemingAnchorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlemingAnchor *FlemingAnchorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.FlemingAnchorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlemingAnchor *FlemingAnchorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.FlemingAnchorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlemingAnchor *FlemingAnchorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlemingAnchor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlemingAnchor *FlemingAnchorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlemingAnchor *FlemingAnchorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.contract.Transact(opts, method, params...)
}

// AnchorCount is a free data retrieval call binding the contract method 0x34f96c8c.
//
// Solidity: function anchorCount() view returns(uint256)
func (_FlemingAnchor *FlemingAnchorCaller) AnchorCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "anchorCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AnchorCount is a free data retrieval call binding the contract method 0x34f96c8c.
//
// Solidity: function anchorCount() view returns(uint256)
func (_FlemingAnchor *FlemingAnchorSession) AnchorCount() (*big.Int, error) {
	return _FlemingAnchor.Contract.AnchorCount(&_FlemingAnchor.CallOpts)
}

// AnchorCount is a free data retrieval call binding the contract method 0x34f96c8c.
//
// Solidity: function anchorCount() view returns(uint256)
func (_FlemingAnchor *FlemingAnchorCallerSession) AnchorCount() (*big.Int, error) {
	return _FlemingAnchor.Contract.AnchorCount(&_FlemingAnchor.CallOpts)
}

// Anchorer is a free data retrieval call binding the contract method 0x84609921.
//
// Solidity: function anchorer() view returns(address)
func (_FlemingAnchor *FlemingAnchorCaller) Anchorer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "anchorer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Anchorer is a free data retrieval call binding the contract method 0x84609921.
//
// Solidity: function anchorer() view returns(address)
func (_FlemingAnchor *FlemingAnchorSession) Anchorer() (common.Address, error) {
	return _FlemingAnchor.Contract.Anchorer(&_FlemingAnchor.CallOpts)
}

// Anchorer is a free data retrieval call binding the contract method 0x84609921.
//
// Solidity: function anchorer() view returns(address)
func (_FlemingAnchor *FlemingAnchorCallerSession) Anchorer() (common.Address, error) {
	return _FlemingAnchor.Contract.Anchorer(&_FlemingAnchor.CallOpts)
}

// Anchors is a free data retrieval call binding the contract method 0xb01b6d53.
//
// Solidity: function anchors(bytes32 ) view returns(uint256)
func (_FlemingAnchor *FlemingAnchorCaller) Anchors(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "anchors", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Anchors is a free data retrieval call binding the contract method 0xb01b6d53.
//
// Solidity: function anchors(bytes32 ) view returns(uint256)
func (_FlemingAnchor *FlemingAnchorSession) Anchors(arg0 [32]byte) (*big.Int, error) {
	return _FlemingAnchor.Contract.Anchors(&_FlemingAnchor.CallOpts, arg0)
}

// Anchors is a free data retrieval call binding the contract method 0xb01b6d53.
//
// Solidity: function anchors(bytes32 ) view returns(uint256)
func (_FlemingAnchor *FlemingAnchorCallerSession) Anchors(arg0 [32]byte) (*big.Int, error) {
	return _FlemingAnchor.Contract.Anchors(&_FlemingAnchor.CallOpts, arg0)
}

// GetAnchorAge is a free data retrieval call binding the contract method 0xee220119.
//
// Solidity: function getAnchorAge(bytes32 root) view returns(uint256 age)
func (_FlemingAnchor *FlemingAnchorCaller) GetAnchorAge(opts *bind.CallOpts, root [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "getAnchorAge", root)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnchorAge is a free data retrieval call binding the contract method 0xee220119.
//
// Solidity: function getAnchorAge(bytes32 root) view returns(uint256 age)
func (_FlemingAnchor *FlemingAnchorSession) GetAnchorAge(root [32]byte) (*big.Int, error) {
	return _FlemingAnchor.Contract.GetAnchorAge(&_FlemingAnchor.CallOpts, root)
}

// GetAnchorAge is a free data retrieval call binding the contract method 0xee220119.
//
// Solidity: function getAnchorAge(bytes32 root) view returns(uint256 age)
func (_FlemingAnchor *FlemingAnchorCallerSession) GetAnchorAge(root [32]byte) (*big.Int, error) {
	return _FlemingAnchor.Contract.GetAnchorAge(&_FlemingAnchor.CallOpts, root)
}

// GetAnchorTimestamp is a free data retrieval call binding the contract method 0xb34e44e4.
//
// Solidity: function getAnchorTimestamp(bytes32 root) view returns(uint256 timestamp)
func (_FlemingAnchor *FlemingAnchorCaller) GetAnchorTimestamp(opts *bind.CallOpts, root [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "getAnchorTimestamp", root)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAnchorTimestamp is a free data retrieval call binding the contract method 0xb34e44e4.
//
// Solidity: function getAnchorTimestamp(bytes32 root) view returns(uint256 timestamp)
func (_FlemingAnchor *FlemingAnchorSession) GetAnchorTimestamp(root [32]byte) (*big.Int, error) {
	return _FlemingAnchor.Contract.GetAnchorTimestamp(&_FlemingAnchor.CallOpts, root)
}

// GetAnchorTimestamp is a free data retrieval call binding the contract method 0xb34e44e4.
//
// Solidity: function getAnchorTimestamp(bytes32 root) view returns(uint256 timestamp)
func (_FlemingAnchor *FlemingAnchorCallerSession) GetAnchorTimestamp(root [32]byte) (*big.Int, error) {
	return _FlemingAnchor.Contract.GetAnchorTimestamp(&_FlemingAnchor.CallOpts, root)
}

// GetContractInfo is a free data retrieval call binding the contract method 0x7cc1f867.
//
// Solidity: function getContractInfo() view returns(string version, uint256 chainId)
func (_FlemingAnchor *FlemingAnchorCaller) GetContractInfo(opts *bind.CallOpts) (struct {
	Version string
	ChainId *big.Int
}, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "getContractInfo")

	outstruct := new(struct {
		Version string
		ChainId *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Version = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.ChainId = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetContractInfo is a free data retrieval call binding the contract method 0x7cc1f867.
//
// Solidity: function getContractInfo() view returns(string version, uint256 chainId)
func (_FlemingAnchor *FlemingAnchorSession) GetContractInfo() (struct {
	Version string
	ChainId *big.Int
}, error) {
	return _FlemingAnchor.Contract.GetContractInfo(&_FlemingAnchor.CallOpts)
}

// GetContractInfo is a free data retrieval call binding the contract method 0x7cc1f867.
//
// Solidity: function getContractInfo() view returns(string version, uint256 chainId)
func (_FlemingAnchor *FlemingAnchorCallerSession) GetContractInfo() (struct {
	Version string
	ChainId *big.Int
}, error) {
	return _FlemingAnchor.Contract.GetContractInfo(&_FlemingAnchor.CallOpts)
}

// IsAnchored is a free data retrieval call binding the contract method 0x4f0b5801.
//
// Solidity: function isAnchored(bytes32 root) view returns(bool)
func (_FlemingAnchor *FlemingAnchorCaller) IsAnchored(opts *bind.CallOpts, root [32]byte) (bool, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "isAnchored", root)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAnchored is a free data retrieval call binding the contract method 0x4f0b5801.
//
// Solidity: function isAnchored(bytes32 root) view returns(bool)
func (_FlemingAnchor *FlemingAnchorSession) IsAnchored(root [32]byte) (bool, error) {
	return _FlemingAnchor.Contract.IsAnchored(&_FlemingAnchor.CallOpts, root)
}

// IsAnchored is a free data retrieval call binding the contract method 0x4f0b5801.
//
// Solidity: function isAnchored(bytes32 root) view returns(bool)
func (_FlemingAnchor *FlemingAnchorCallerSession) IsAnchored(root [32]byte) (bool, error) {
	return _FlemingAnchor.Contract.IsAnchored(&_FlemingAnchor.CallOpts, root)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FlemingAnchor *FlemingAnchorCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlemingAnchor.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FlemingAnchor *FlemingAnchorSession) Owner() (common.Address, error) {
	return _FlemingAnchor.Contract.Owner(&_FlemingAnchor.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FlemingAnchor *FlemingAnchorCallerSession) Owner() (common.Address, error) {
	return _FlemingAnchor.Contract.Owner(&_FlemingAnchor.CallOpts)
}

// Anchor is a paid mutator transaction binding the contract method 0xeecdf927.
//
// Solidity: function anchor(bytes32 root) returns()
func (_FlemingAnchor *FlemingAnchorTransactor) Anchor(opts *bind.TransactOpts, root [32]byte) (*types.Transaction, error) {
	return _FlemingAnchor.contract.Transact(opts, "anchor", root)
}

// Anchor is a paid mutator transaction binding the contract method 0xeecdf927.
//
// Solidity: function anchor(bytes32 root) returns()
func (_FlemingAnchor *FlemingAnchorSession) Anchor(root [32]byte) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.Anchor(&_FlemingAnchor.TransactOpts, root)
}

// Anchor is a paid mutator transaction binding the contract method 0xeecdf927.
//
// Solidity: function anchor(bytes32 root) returns()
func (_FlemingAnchor *FlemingAnchorTransactorSession) Anchor(root [32]byte) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.Anchor(&_FlemingAnchor.TransactOpts, root)
}

// BatchAnchor is a paid mutator transaction binding the contract method 0x358eb477.
//
// Solidity: function batchAnchor(bytes32[] roots) returns()
func (_FlemingAnchor *FlemingAnchorTransactor) BatchAnchor(opts *bind.TransactOpts, roots [][32]byte) (*types.Transaction, error) {
	return _FlemingAnchor.contract.Transact(opts, "batchAnchor", roots)
}

// BatchAnchor is a paid mutator transaction binding the contract method 0x358eb477.
//
// Solidity: function batchAnchor(bytes32[] roots) returns()
func (_FlemingAnchor *FlemingAnchorSession) BatchAnchor(roots [][32]byte) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.BatchAnchor(&_FlemingAnchor.TransactOpts, roots)
}

// BatchAnchor is a paid mutator transaction binding the contract method 0x358eb477.
//
// Solidity: function batchAnchor(bytes32[] roots) returns()
func (_FlemingAnchor *FlemingAnchorTransactorSession) BatchAnchor(roots [][32]byte) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.BatchAnchor(&_FlemingAnchor.TransactOpts, roots)
}

// RenounceAnchorer is a paid mutator transaction binding the contract method 0x9ff9e2be.
//
// Solidity: function renounceAnchorer() returns()
func (_FlemingAnchor *FlemingAnchorTransactor) RenounceAnchorer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlemingAnchor.contract.Transact(opts, "renounceAnchorer")
}

// RenounceAnchorer is a paid mutator transaction binding the contract method 0x9ff9e2be.
//
// Solidity: function renounceAnchorer() returns()
func (_FlemingAnchor *FlemingAnchorSession) RenounceAnchorer() (*types.Transaction, error) {
	return _FlemingAnchor.Contract.RenounceAnchorer(&_FlemingAnchor.TransactOpts)
}

// RenounceAnchorer is a paid mutator transaction binding the contract method 0x9ff9e2be.
//
// Solidity: function renounceAnchorer() returns()
func (_FlemingAnchor *FlemingAnchorTransactorSession) RenounceAnchorer() (*types.Transaction, error) {
	return _FlemingAnchor.Contract.RenounceAnchorer(&_FlemingAnchor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FlemingAnchor *FlemingAnchorTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlemingAnchor.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FlemingAnchor *FlemingAnchorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FlemingAnchor.Contract.RenounceOwnership(&_FlemingAnchor.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FlemingAnchor *FlemingAnchorTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FlemingAnchor.Contract.RenounceOwnership(&_FlemingAnchor.TransactOpts)
}

// SetAnchorer is a paid mutator transaction binding the contract method 0x0e0ab0a1.
//
// Solidity: function setAnchorer(address newAnchorer) returns()
func (_FlemingAnchor *FlemingAnchorTransactor) SetAnchorer(opts *bind.TransactOpts, newAnchorer common.Address) (*types.Transaction, error) {
	return _FlemingAnchor.contract.Transact(opts, "setAnchorer", newAnchorer)
}

// SetAnchorer is a paid mutator transaction binding the contract method 0x0e0ab0a1.
//
// Solidity: function setAnchorer(address newAnchorer) returns()
func (_FlemingAnchor *FlemingAnchorSession) SetAnchorer(newAnchorer common.Address) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.SetAnchorer(&_FlemingAnchor.TransactOpts, newAnchorer)
}

// SetAnchorer is a paid mutator transaction binding the contract method 0x0e0ab0a1.
//
// Solidity: function setAnchorer(address newAnchorer) returns()
func (_FlemingAnchor *FlemingAnchorTransactorSession) SetAnchorer(newAnchorer common.Address) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.SetAnchorer(&_FlemingAnchor.TransactOpts, newAnchorer)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FlemingAnchor *FlemingAnchorTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FlemingAnchor.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FlemingAnchor *FlemingAnchorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.TransferOwnership(&_FlemingAnchor.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FlemingAnchor *FlemingAnchorTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FlemingAnchor.Contract.TransferOwnership(&_FlemingAnchor.TransactOpts, newOwner)
}

// FlemingAnchorAnchorerUpdatedIterator is returned from FilterAnchorerUpdated and is used to iterate over the raw logs and unpacked data for AnchorerUpdated events raised by the FlemingAnchor contract.
type FlemingAnchorAnchorerUpdatedIterator struct {
	Event *FlemingAnchorAnchorerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlemingAnchorAnchorerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlemingAnchorAnchorerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlemingAnchorAnchorerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlemingAnchorAnchorerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlemingAnchorAnchorerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlemingAnchorAnchorerUpdated represents a AnchorerUpdated event raised by the FlemingAnchor contract.
type FlemingAnchorAnchorerUpdated struct {
	PreviousAnchorer common.Address
	NewAnchorer      common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterAnchorerUpdated is a free log retrieval operation binding the contract event 0x2aa923b6c5e7fed8e6d79877345c501d2cd5ed651aa3d3aff9127a6b7f44f34b.
//
// Solidity: event AnchorerUpdated(address indexed previousAnchorer, address indexed newAnchorer)
func (_FlemingAnchor *FlemingAnchorFilterer) FilterAnchorerUpdated(opts *bind.FilterOpts, previousAnchorer []common.Address, newAnchorer []common.Address) (*FlemingAnchorAnchorerUpdatedIterator, error) {

	var previousAnchorerRule []interface{}
	for _, previousAnchorerItem := range previousAnchorer {
		previousAnchorerRule = append(previousAnchorerRule, previousAnchorerItem)
	}
	var newAnchorerRule []interface{}
	for _, newAnchorerItem := range newAnchorer {
		newAnchorerRule = append(newAnchorerRule, newAnchorerItem)
	}

	logs, sub, err := _FlemingAnchor.contract.FilterLogs(opts, "AnchorerUpdated", previousAnchorerRule, newAnchorerRule)
	if err != nil {
		return nil, err
	}
	return &FlemingAnchorAnchorerUpdatedIterator{contract: _FlemingAnchor.contract, event: "AnchorerUpdated", logs: logs, sub: sub}, nil
}

// WatchAnchorerUpdated is a free log subscription operation binding the contract event 0x2aa923b6c5e7fed8e6d79877345c501d2cd5ed651aa3d3aff9127a6b7f44f34b.
//
// Solidity: event AnchorerUpdated(address indexed previousAnchorer, address indexed newAnchorer)
func (_FlemingAnchor *FlemingAnchorFilterer) WatchAnchorerUpdated(opts *bind.WatchOpts, sink chan<- *FlemingAnchorAnchorerUpdated, previousAnchorer []common.Address, newAnchorer []common.Address) (event.Subscription, error) {

	var previousAnchorerRule []interface{}
	for _, previousAnchorerItem := range previousAnchorer {
		previousAnchorerRule = append(previousAnchorerRule, previousAnchorerItem)
	}
	var newAnchorerRule []interface{}
	for _, newAnchorerItem := range newAnchorer {
		newAnchorerRule = append(newAnchorerRule, newAnchorerItem)
	}

	logs, sub, err := _FlemingAnchor.contract.WatchLogs(opts, "AnchorerUpdated", previousAnchorerRule, newAnchorerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlemingAnchorAnchorerUpdated)
				if err := _FlemingAnchor.contract.UnpackLog(event, "AnchorerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAnchorerUpdated is a log parse operation binding the contract event 0x2aa923b6c5e7fed8e6d79877345c501d2cd5ed651aa3d3aff9127a6b7f44f34b.
//
// Solidity: event AnchorerUpdated(address indexed previousAnchorer, address indexed newAnchorer)
func (_FlemingAnchor *FlemingAnchorFilterer) ParseAnchorerUpdated(log types.Log) (*FlemingAnchorAnchorerUpdated, error) {
	event := new(FlemingAnchorAnchorerUpdated)
	if err := _FlemingAnchor.contract.UnpackLog(event, "AnchorerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlemingAnchorOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FlemingAnchor contract.
type FlemingAnchorOwnershipTransferredIterator struct {
	Event *FlemingAnchorOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlemingAnchorOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlemingAnchorOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlemingAnchorOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlemingAnchorOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlemingAnchorOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlemingAnchorOwnershipTransferred represents a OwnershipTransferred event raised by the FlemingAnchor contract.
type FlemingAnchorOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FlemingAnchor *FlemingAnchorFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FlemingAnchorOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FlemingAnchor.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FlemingAnchorOwnershipTransferredIterator{contract: _FlemingAnchor.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FlemingAnchor *FlemingAnchorFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FlemingAnchorOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FlemingAnchor.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlemingAnchorOwnershipTransferred)
				if err := _FlemingAnchor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FlemingAnchor *FlemingAnchorFilterer) ParseOwnershipTransferred(log types.Log) (*FlemingAnchorOwnershipTransferred, error) {
	event := new(FlemingAnchorOwnershipTransferred)
	if err := _FlemingAnchor.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlemingAnchorRootAnchoredIterator is returned from FilterRootAnchored and is used to iterate over the raw logs and unpacked data for RootAnchored events raised by the FlemingAnchor contract.
type FlemingAnchorRootAnchoredIterator struct {
	Event *FlemingAnchorRootAnchored // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlemingAnchorRootAnchoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlemingAnchorRootAnchored)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlemingAnchorRootAnchored)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlemingAnchorRootAnchoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlemingAnchorRootAnchoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlemingAnchorRootAnchored represents a RootAnchored event raised by the FlemingAnchor contract.
type FlemingAnchorRootAnchored struct {
	Root        [32]byte
	Timestamp   *big.Int
	BlockNumber *big.Int
	Anchorer    common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRootAnchored is a free log retrieval operation binding the contract event 0xeab437c5ee6f90e75abb8e8f688a46b43d963cddd68c980867e9c5e2b8cd2b06.
//
// Solidity: event RootAnchored(bytes32 indexed root, uint256 timestamp, uint256 blockNumber, address indexed anchorer)
func (_FlemingAnchor *FlemingAnchorFilterer) FilterRootAnchored(opts *bind.FilterOpts, root [][32]byte, anchorer []common.Address) (*FlemingAnchorRootAnchoredIterator, error) {

	var rootRule []interface{}
	for _, rootItem := range root {
		rootRule = append(rootRule, rootItem)
	}

	var anchorerRule []interface{}
	for _, anchorerItem := range anchorer {
		anchorerRule = append(anchorerRule, anchorerItem)
	}

	logs, sub, err := _FlemingAnchor.contract.FilterLogs(opts, "RootAnchored", rootRule, anchorerRule)
	if err != nil {
		return nil, err
	}
	return &FlemingAnchorRootAnchoredIterator{contract: _FlemingAnchor.contract, event: "RootAnchored", logs: logs, sub: sub}, nil
}

// WatchRootAnchored is a free log subscription operation binding the contract event 0xeab437c5ee6f90e75abb8e8f688a46b43d963cddd68c980867e9c5e2b8cd2b06.
//
// Solidity: event RootAnchored(bytes32 indexed root, uint256 timestamp, uint256 blockNumber, address indexed anchorer)
func (_FlemingAnchor *FlemingAnchorFilterer) WatchRootAnchored(opts *bind.WatchOpts, sink chan<- *FlemingAnchorRootAnchored, root [][32]byte, anchorer []common.Address) (event.Subscription, error) {

	var rootRule []interface{}
	for _, rootItem := range root {
		rootRule = append(rootRule, rootItem)
	}

	var anchorerRule []interface{}
	for _, anchorerItem := range anchorer {
		anchorerRule = append(anchorerRule, anchorerItem)
	}

	logs, sub, err := _FlemingAnchor.contract.WatchLogs(opts, "RootAnchored", rootRule, anchorerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlemingAnchorRootAnchored)
				if err := _FlemingAnchor.contract.UnpackLog(event, "RootAnchored", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRootAnchored is a log parse operation binding the contract event 0xeab437c5ee6f90e75abb8e8f688a46b43d963cddd68c980867e9c5e2b8cd2b06.
//
// Solidity: event RootAnchored(bytes32 indexed root, uint256 timestamp, uint256 blockNumber, address indexed anchorer)
func (_FlemingAnchor *FlemingAnchorFilterer) ParseRootAnchored(log types.Log) (*FlemingAnchorRootAnchored, error) {
	event := new(FlemingAnchorRootAnchored)
	if err := _FlemingAnchor.contract.UnpackLog(event, "RootAnchored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
