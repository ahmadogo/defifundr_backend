// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package defi

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

// DefiMetaData contains all meta data concerning the Defi contract.
var DefiMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"type\":\"function\"}]",
}

// DefiABI is the input ABI used to generate the binding from.
// Deprecated: Use DefiMetaData.ABI instead.
var DefiABI = DefiMetaData.ABI

// Defi is an auto generated Go binding around an Ethereum contract.
type Defi struct {
	DefiCaller     // Read-only binding to the contract
	DefiTransactor // Write-only binding to the contract
	DefiFilterer   // Log filterer for contract events
}

// DefiCaller is an auto generated read-only Go binding around an Ethereum contract.
type DefiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DefiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DefiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DefiSession struct {
	Contract     *Defi             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DefiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DefiCallerSession struct {
	Contract *DefiCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DefiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DefiTransactorSession struct {
	Contract     *DefiTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DefiRaw is an auto generated low-level Go binding around an Ethereum contract.
type DefiRaw struct {
	Contract *Defi // Generic contract binding to access the raw methods on
}

// DefiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DefiCallerRaw struct {
	Contract *DefiCaller // Generic read-only contract binding to access the raw methods on
}

// DefiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DefiTransactorRaw struct {
	Contract *DefiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDefi creates a new instance of Defi, bound to a specific deployed contract.
func NewDefi(address common.Address, backend bind.ContractBackend) (*Defi, error) {
	contract, err := bindDefi(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Defi{DefiCaller: DefiCaller{contract: contract}, DefiTransactor: DefiTransactor{contract: contract}, DefiFilterer: DefiFilterer{contract: contract}}, nil
}

// NewDefiCaller creates a new read-only instance of Defi, bound to a specific deployed contract.
func NewDefiCaller(address common.Address, caller bind.ContractCaller) (*DefiCaller, error) {
	contract, err := bindDefi(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DefiCaller{contract: contract}, nil
}

// NewDefiTransactor creates a new write-only instance of Defi, bound to a specific deployed contract.
func NewDefiTransactor(address common.Address, transactor bind.ContractTransactor) (*DefiTransactor, error) {
	contract, err := bindDefi(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DefiTransactor{contract: contract}, nil
}

// NewDefiFilterer creates a new log filterer instance of Defi, bound to a specific deployed contract.
func NewDefiFilterer(address common.Address, filterer bind.ContractFilterer) (*DefiFilterer, error) {
	contract, err := bindDefi(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DefiFilterer{contract: contract}, nil
}

// bindDefi binds a generic wrapper to an already deployed contract.
func bindDefi(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DefiMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Defi *DefiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Defi.Contract.DefiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Defi *DefiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Defi.Contract.DefiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Defi *DefiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Defi.Contract.DefiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Defi *DefiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Defi.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Defi *DefiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Defi.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Defi *DefiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Defi.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) returns(uint256 balance)
func (_Defi *DefiCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Defi.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) returns(uint256 balance)
func (_Defi *DefiSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Defi.Contract.BalanceOf(&_Defi.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) returns(uint256 balance)
func (_Defi *DefiCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Defi.Contract.BalanceOf(&_Defi.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() returns(uint8)
func (_Defi *DefiCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Defi.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() returns(uint8)
func (_Defi *DefiSession) Decimals() (uint8, error) {
	return _Defi.Contract.Decimals(&_Defi.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() returns(uint8)
func (_Defi *DefiCallerSession) Decimals() (uint8, error) {
	return _Defi.Contract.Decimals(&_Defi.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string)
func (_Defi *DefiCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Defi.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string)
func (_Defi *DefiSession) Name() (string, error) {
	return _Defi.Contract.Name(&_Defi.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() returns(string)
func (_Defi *DefiCallerSession) Name() (string, error) {
	return _Defi.Contract.Name(&_Defi.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() returns(string)
func (_Defi *DefiCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Defi.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() returns(string)
func (_Defi *DefiSession) Symbol() (string, error) {
	return _Defi.Contract.Symbol(&_Defi.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() returns(string)
func (_Defi *DefiCallerSession) Symbol() (string, error) {
	return _Defi.Contract.Symbol(&_Defi.CallOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Defi *DefiTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Defi.contract.Transact(opts, "transfer", _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Defi *DefiSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Defi.Contract.Transfer(&_Defi.TransactOpts, _to, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _value) returns(bool)
func (_Defi *DefiTransactorSession) Transfer(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Defi.Contract.Transfer(&_Defi.TransactOpts, _to, _value)
}
