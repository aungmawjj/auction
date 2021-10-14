// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AssetABI is the input ABI used to generate the binding from.
const AssetABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"winer\",\"type\":\"address\"}],\"name\":\"endAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pendingAuction\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"auction\",\"type\":\"address\"}],\"name\":\"startAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

var AssetParsedABI, _ = abi.JSON(strings.NewReader(AssetABI))

// AssetFuncSigs maps the 4-byte function signature to its string representation.
var AssetFuncSigs = map[string]string{
	"2872abbc": "endAuction(address)",
	"8da5cb5b": "owner()",
	"2210ae4b": "pendingAuction()",
	"4075fa0f": "startAuction(address)",
}

// AssetBin is the compiled bytecode used for deploying new contracts.
var AssetBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163217905561016e806100326000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80632210ae4b146100515780632872abbc146100755780634075fa0f1461009d5780638da5cb5b146100c3575b600080fd5b6100596100cb565b604080516001600160a01b039092168252519081900360200190f35b61009b6004803603602081101561008b57600080fd5b50356001600160a01b03166100da565b005b61009b600480360360208110156100b357600080fd5b50356001600160a01b0316610108565b61005961012a565b6001546001600160a01b031681565b600180546001600160a01b0319908116909155600080549091166001600160a01b0392909216919091179055565b600180546001600160a01b0319166001600160a01b0392909216919091179055565b6000546001600160a01b03168156fea265627a7a72315820391cd3b6b877c16e5b0c21e8bc702fee59dc2314afb924b7e6eeff214812651564736f6c63430005110032"

// DeployAsset deploys a new Ethereum contract, binding an instance of Asset to it.
func DeployAsset(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Asset, error) {
	parsed, err := abi.JSON(strings.NewReader(AssetABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AssetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Asset{AssetCaller: AssetCaller{contract: contract}, AssetTransactor: AssetTransactor{contract: contract}, AssetFilterer: AssetFilterer{contract: contract}}, nil
}

// Asset is an auto generated Go binding around an Ethereum contract.
type Asset struct {
	AssetCaller     // Read-only binding to the contract
	AssetTransactor // Write-only binding to the contract
	AssetFilterer   // Log filterer for contract events
}

// AssetCaller is an auto generated read-only Go binding around an Ethereum contract.
type AssetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AssetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AssetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AssetSession struct {
	Contract     *Asset            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AssetCallerSession struct {
	Contract *AssetCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AssetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AssetTransactorSession struct {
	Contract     *AssetTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetRaw is an auto generated low-level Go binding around an Ethereum contract.
type AssetRaw struct {
	Contract *Asset // Generic contract binding to access the raw methods on
}

// AssetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AssetCallerRaw struct {
	Contract *AssetCaller // Generic read-only contract binding to access the raw methods on
}

// AssetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AssetTransactorRaw struct {
	Contract *AssetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAsset creates a new instance of Asset, bound to a specific deployed contract.
func NewAsset(address common.Address, backend bind.ContractBackend) (*Asset, error) {
	contract, err := bindAsset(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Asset{AssetCaller: AssetCaller{contract: contract}, AssetTransactor: AssetTransactor{contract: contract}, AssetFilterer: AssetFilterer{contract: contract}}, nil
}

// NewAssetCaller creates a new read-only instance of Asset, bound to a specific deployed contract.
func NewAssetCaller(address common.Address, caller bind.ContractCaller) (*AssetCaller, error) {
	contract, err := bindAsset(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AssetCaller{contract: contract}, nil
}

// NewAssetTransactor creates a new write-only instance of Asset, bound to a specific deployed contract.
func NewAssetTransactor(address common.Address, transactor bind.ContractTransactor) (*AssetTransactor, error) {
	contract, err := bindAsset(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AssetTransactor{contract: contract}, nil
}

// NewAssetFilterer creates a new log filterer instance of Asset, bound to a specific deployed contract.
func NewAssetFilterer(address common.Address, filterer bind.ContractFilterer) (*AssetFilterer, error) {
	contract, err := bindAsset(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AssetFilterer{contract: contract}, nil
}

// bindAsset binds a generic wrapper to an already deployed contract.
func bindAsset(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AssetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Asset *AssetRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Asset.Contract.AssetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Asset *AssetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Asset.Contract.AssetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Asset *AssetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Asset.Contract.AssetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Asset *AssetCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Asset.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Asset *AssetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Asset.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Asset *AssetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Asset.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Asset *AssetCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Asset.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Asset *AssetSession) Owner() (common.Address, error) {
	return _Asset.Contract.Owner(&_Asset.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Asset *AssetCallerSession) Owner() (common.Address, error) {
	return _Asset.Contract.Owner(&_Asset.CallOpts)
}

// PendingAuction is a free data retrieval call binding the contract method 0x2210ae4b.
//
// Solidity: function pendingAuction() constant returns(address)
func (_Asset *AssetCaller) PendingAuction(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Asset.contract.Call(opts, out, "pendingAuction")
	return *ret0, err
}

// PendingAuction is a free data retrieval call binding the contract method 0x2210ae4b.
//
// Solidity: function pendingAuction() constant returns(address)
func (_Asset *AssetSession) PendingAuction() (common.Address, error) {
	return _Asset.Contract.PendingAuction(&_Asset.CallOpts)
}

// PendingAuction is a free data retrieval call binding the contract method 0x2210ae4b.
//
// Solidity: function pendingAuction() constant returns(address)
func (_Asset *AssetCallerSession) PendingAuction() (common.Address, error) {
	return _Asset.Contract.PendingAuction(&_Asset.CallOpts)
}

// EndAuction is a paid mutator transaction binding the contract method 0x2872abbc.
//
// Solidity: function endAuction(address winer) returns()
func (_Asset *AssetTransactor) EndAuction(opts *bind.TransactOpts, winer common.Address) (*types.Transaction, error) {
	return _Asset.contract.Transact(opts, "endAuction", winer)
}

// EndAuction is a paid mutator transaction binding the contract method 0x2872abbc.
//
// Solidity: function endAuction(address winer) returns()
func (_Asset *AssetSession) EndAuction(winer common.Address) (*types.Transaction, error) {
	return _Asset.Contract.EndAuction(&_Asset.TransactOpts, winer)
}

// EndAuction is a paid mutator transaction binding the contract method 0x2872abbc.
//
// Solidity: function endAuction(address winer) returns()
func (_Asset *AssetTransactorSession) EndAuction(winer common.Address) (*types.Transaction, error) {
	return _Asset.Contract.EndAuction(&_Asset.TransactOpts, winer)
}

// StartAuction is a paid mutator transaction binding the contract method 0x4075fa0f.
//
// Solidity: function startAuction(address auction) returns()
func (_Asset *AssetTransactor) StartAuction(opts *bind.TransactOpts, auction common.Address) (*types.Transaction, error) {
	return _Asset.contract.Transact(opts, "startAuction", auction)
}

// StartAuction is a paid mutator transaction binding the contract method 0x4075fa0f.
//
// Solidity: function startAuction(address auction) returns()
func (_Asset *AssetSession) StartAuction(auction common.Address) (*types.Transaction, error) {
	return _Asset.Contract.StartAuction(&_Asset.TransactOpts, auction)
}

// StartAuction is a paid mutator transaction binding the contract method 0x4075fa0f.
//
// Solidity: function startAuction(address auction) returns()
func (_Asset *AssetTransactorSession) StartAuction(auction common.Address) (*types.Transaction, error) {
	return _Asset.Contract.StartAuction(&_Asset.TransactOpts, auction)
}
