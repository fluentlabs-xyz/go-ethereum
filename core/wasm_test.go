package core

import (
	_ "embed"
	"fmt"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/eth/tracers/logger"
	"github.com/stretchr/testify/assert"
	zkwasm_wasmi "github.com/wasm0/zkwasm-wasmi"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	"github.com/wasmerio/wasmer-go/wasmer"
)

func newWasmMachine() (*vm.EVM, *logger.WebAssemblyLogger) {
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	config := params.AllEthashProtocolChanges
	config.WebAssemblyBlock = big.NewInt(0)
	blockCtx := vm.BlockContext{
		Transfer: func(
			vm.StateDB,
			common.Address,
			common.Address,
			*big.Int,
		) {
		},
		BlockNumber: big.NewInt(1),
		Time:        big.NewInt(2),
		Difficulty:  big.NewInt(3),
		BaseFee:     big.NewInt(4),
	}
	txCtx := vm.TxContext{}
	tracer := logger.NewWebAssemblyLogger(&logger.Config{
		EnableMemory:     false,
		DisableStack:     false,
		DisableStorage:   false,
		EnableReturnData: true,
		Debug:            true,
		Limit:            0,
	})
	evm := vm.NewEVM(
		blockCtx, txCtx, statedb, config, vm.Config{
			Tracer: tracer,
			Debug:  true,
		},
	)
	return evm, tracer
}

func newWasmContract(evm *vm.EVM, addr common.Address, watCode string) {
	wasmCode, err := wasmer.Wat2Wasm(watCode)
	if err != nil {
		panic(err)
	}
	evm.StateDB.SetCode(addr, wasmCode)
}

//go:embed vm/testdata/wasm/hello.wasm
var wasmTestHello string

//go:embed vm/testdata/wasm/hello_injected.wat
var watTestHelloInjected string

//go:embed vm/testdata/wasm/greeting.wasm
var wasmTestGreeting string

//go:embed vm/testdata/wasm/greeting.wat
var watTestGreeting string

//go:embed vm/testdata/wasm/deploy.wasm
var wasmTestDeploy string

//go:embed vm/testdata/wasm/deploy.wat
var watTestDeploy string

//go:embed vm/testdata/wasm/simple.wasm
var wasmTestSimple string

//go:embed vm/testdata/wasm/simple.wat
var watTestSimple string

const AddressFunctionFlag = 1 << 0
const CallValueFunctionFlag = 1 << 1
const TimestampFunctionFlag = 1 << 2
const BalanceFunctionFlag = 1 << 3
const CallerFunctionFlag = 1 << 4
const OriginFunctionFlag = 1 << 5
const PanickingFunctionFlag = 1 << 6

func expectGasLeft(t *testing.T, interpreter vm.VirtualInterpreter, expectedGasLeft uint64, msgAndArgs ...interface{}) {
	wasmInterpreter, ok := interpreter.(*vm.WASMInterpreter)
	assert.Truef(t, ok, "expected evm.Interpreter() to be of type *vm.WASMInterpreter")
	assert.Equal(t, expectedGasLeft, wasmInterpreter.Scope().Contract.Gas, msgAndArgs)
}

func TestWASMInterpreter_Hello(t *testing.T) {
	{
		evm, tracer := newWasmMachine()
		newWasmContract(evm, common.Address{100, 20, 3}, wasmTestHello)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{100, 20, 3}, []byte{AddressFunctionFlag}, 10_000_000, big.NewInt(0))
		require.EqualError(t, err, "exit return code: 123")
		msg, _ := tracer.GetResult()
		msg, _ = msg.MarshalJSON()
		println(string(msg))
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, wasmTestHello)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{CallValueFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, wasmTestHello)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{TimestampFunctionFlag}, 10_000_000, big.NewInt(0))
		require.EqualError(t, err, "exit return code: 2")
	}

	{
		evm, tracer := newWasmMachine()
		newWasmContract(evm, common.Address{}, wasmTestHello)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{BalanceFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
		msg, _ := tracer.GetResult()
		msg, _ = msg.MarshalJSON()
		println(string(msg))
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, wasmTestHello)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{CallerFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, wasmTestHello)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{OriginFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
	}

	{
		//evm, _ := newWasmMachine()
		//newWasmContract(evm, common.Address{}, wasmTestHello)
		//_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{PanickingFunctionFlag}, 10_000_000, big.NewInt(0))
		//require.NoError(t, err)
	}
}

func TestWASMInterpreter_Hello_injected(t *testing.T) {
	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{100, 20, 3}, watTestHelloInjected)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{100, 20, 3}, []byte{AddressFunctionFlag}, 10_000_000, big.NewInt(0))
		require.EqualError(t, err, "exit return code: 123")
		expectGasLeft(t, evm.Interpreter(), 0x987012, "AddressFunctionFlag")
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, watTestHelloInjected)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{CallValueFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
		expectGasLeft(t, evm.Interpreter(), 0x97fa68, "CallValueFunctionFlag")
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, watTestHelloInjected)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{TimestampFunctionFlag}, 10_000_000, big.NewInt(0))
		require.EqualError(t, err, "exit return code: 2")
		expectGasLeft(t, evm.Interpreter(), 0x9887c2, "TimestampFunctionFlag")
	}

	{
		evm, tracer := newWasmMachine()
		newWasmContract(evm, common.Address{}, watTestHelloInjected)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{BalanceFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
		msg, _ := tracer.GetResult()
		msg, _ = msg.MarshalJSON()
		expectGasLeft(t, evm.Interpreter(), 0x985251, "BalanceFunctionFlag")
		println(string(msg))
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, watTestHelloInjected)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{CallerFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
		expectGasLeft(t, evm.Interpreter(), 0x986f8b, "CallerFunctionFlag")
	}

	{
		evm, _ := newWasmMachine()
		newWasmContract(evm, common.Address{}, watTestHelloInjected)
		_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{OriginFunctionFlag}, 10_000_000, big.NewInt(0))
		require.NoError(t, err)
		expectGasLeft(t, evm.Interpreter(), 0x986f6d, "OriginFunctionFlag")
	}
}

func TestWASMInterpreter_Deploy(t *testing.T) {
	evm, _ := newWasmMachine()
	newWasmContract(evm, common.Address{}, watTestDeploy)
	ret, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{}, 10_000_000, big.NewInt(0))
	require.NoError(t, err)
	expected := "0061736d01000000010b0260027f7f0060017f017f02130103656e760b5f65766d5f72657475726e0000030201010405017001010105030100110619037f01418080c0000b7f00418c80c0000b7f00419080c0000b072c04066d656d6f72790200046d61696e00010a5f5f646174615f656e6403010b5f5f686561705f6261736503020a0f010d00418080c000410c100041000b0b150100418080c0000b0c48656c6c6f2c20576f726c64"
	fmt.Printf("returned: %x\n", expected)
	fmt.Printf("expected: %x\n", ret)
	require.Equal(
		t,
		ret,
		hexutil.MustDecode("0x"+expected),
	)
	expectGasLeft(t, evm.Interpreter(), 0x989623)
}

func TestWASMInterpreter_Greeting(t *testing.T) {
	evm, tracer := newWasmMachine()
	newWasmContract(evm, common.Address{}, watTestGreeting)
	ret, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{}, 10_000_000, big.NewInt(0))
	msg, _ := tracer.GetResult()
	msg, _ = msg.MarshalJSON()
	println(string(msg))
	require.NoError(t, err)
	require.Equal(t, ret, []byte("Hello, World"))
	//expectGasLeft(t, evm.Interpreter(), 0x989623)
}

func TestWASMInterpreter_SimpleWasmFile(t *testing.T) {
	evm, tracer := newWasmMachine()
	newWasmContract(evm, common.Address{}, watTestSimple)
	_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{}, 10_000_000, big.NewInt(0))
	require.NoError(t, err)
	msg, _ := tracer.GetResult()
	msg, _ = msg.MarshalJSON()
	println(string(msg))
}

func TestWASMInterpreter_SimpleWasmFile__out_of_gas(t *testing.T) {
	evm, tracer := newWasmMachine()
	newWasmContract(evm, common.Address{}, watTestSimple)
	_, _, err := evm.Call(vm.AccountRef(common.Address{}), common.Address{}, []byte{}, 10, big.NewInt(0))
	require.Error(t, err, zkwasm_wasmi.ComputeTraceErrorCodeOutOfGas)
	msg, _ := tracer.GetResult()
	msg, _ = msg.MarshalJSON()
	println(string(msg))
}
