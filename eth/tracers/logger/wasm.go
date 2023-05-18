package logger

import (
	"encoding/json"
	"fmt"
	"github.com/holiman/uint256"
	"github.com/sbinet/wasm"
	"github.com/scroll-tech/go-ethereum/common"
	"github.com/scroll-tech/go-ethereum/common/hexutil"
	"github.com/scroll-tech/go-ethereum/common/math"
	"github.com/scroll-tech/go-ethereum/core/vm"
	"math/big"
	"strings"
	"sync/atomic"
	"time"
)

//go:generate go run github.com/fjl/gencodec -type WasmLog -field-override wasmLogMarshaling -out gen_wasmlog.go

type OpCodeFamily int

const (
	OpCodeFamilyUnknown OpCodeFamily = iota
	OpCodeFamilyEVM
	OpCodeFamilyWASM
	OpCodeFamilyGAS
)

func (f OpCodeFamily) String() string {
	switch f {
	case OpCodeFamilyWASM:
		return "WASM"
	case OpCodeFamilyEVM:
		return "EVM"
	case OpCodeFamilyGAS:
		return "GAS"
	}
	return "UNKNOWN"
}

type WasmGlobal struct {
	Pc     uint64   `json:"pc"`
	Index  uint64   `json:"index"`
	Op     string   `json:"op"`
	Params []uint64 `json:"params,omitempty"`
	Value  uint64   `json:"value"`
}

// WasmLog is emitted to the EVM each cycle and lists information about the current internal state
// prior to the execution of the statement.
type WasmLog struct {
	Pc            uint64                      `json:"pc"`
	OpFamily      OpCodeFamily                `json:"opcodeFamily,omitempty"`
	Op            vm.OpCodeInfo               `json:"op"`
	Params        []uint64                    `json:"params,omitempty"`
	Gas           uint64                      `json:"gas"`
	GasCost       uint64                      `json:"gasCost"`
	Memory        []byte                      `json:"memory,omitempty"`
	MemoryOffset  uint32                      `json:"memoryOffset,omitempty"`
	MemorySize    uint32                      `json:"memSize"`
	Stack         []uint256.Int               `json:"stack"`
	ReturnData    []byte                      `json:"returnData,omitempty"`
	Storage       map[common.Hash]common.Hash `json:"-"`
	Depth         int                         `json:"depth"`
	RefundCounter uint64                      `json:"refund"`
	Err           error                       `json:"-"`
	Keep          uint32                      `json:"keep"`
	Drop          uint32                      `json:"drop"`
}

type WasmFnCallLog struct {
	FnIndex        uint32 `json:"fnIndex"`
	MaxStackHeight uint32 `json:"maxStackHeight"`
	NumLocals      uint32 `json:"numLocals"`
	FnName         string `json:"fnName"`
}

// overrides for gencodec
type wasmLogMarshaling struct {
	Gas         math.HexOrDecimal64
	GasCost     math.HexOrDecimal64
	Memory      hexutil.Bytes
	ReturnData  hexutil.Bytes
	OpName      string `json:"opName"`          // adds call to OpName() in MarshalJSON
	ErrorString string `json:"error,omitempty"` // adds call to ErrorString() in MarshalJSON
}

// OpName formats the operand name in a human-readable format.
func (s *WasmLog) OpName() string {
	if s.OpFamily == OpCodeFamilyWASM {
		opName, ok := vm.WasmOpCodeToName[wasm.Opcode(s.Op.Code())]
		if !ok {
			return fmt.Sprintf("unknown opcode name: (code=0x%x, family=%s)", s.Op.Code(), s.Op.String())
		}
		return opName
	} else if s.OpFamily == OpCodeFamilyEVM {
		return fmt.Sprintf("evm_%s", strings.ToLower(s.Op.String()))
	} else if s.OpFamily == OpCodeFamilyGAS {
		return strings.ToLower(s.Op.String())
	} else {
		panic("unknown family, not possible")
	}
}

func (s *WasmLog) OpParams() []uint64 {
	params := s.Op.GetParams()
	// somehow wazero returns -1 for end opcode in params
	if wasm.Opcode(s.Op.Code()) == wasm.Op_end {
		params = []uint64{}
	}
	return params
}

// ErrorString formats the log's error as a string.
func (s *WasmLog) ErrorString() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}

type WebAssemblyLogger struct {
	cfg Config
	env *vm.EVM

	storage       map[common.Address]Storage
	logs          []WasmLog
	functionCalls []WasmFnCallLog
	globals       []WasmGlobal
	output        []byte
	globalMemory  map[uint32][]byte
	err           error
	gasLimit      uint64
	usedGas       uint64

	interrupt uint32
	reason    error
}

func NewWebAssemblyLogger(cfg *Config) *WebAssemblyLogger {
	logger := &WebAssemblyLogger{
		storage:       make(map[common.Address]Storage),
		logs:          make([]WasmLog, 0),
		functionCalls: make([]WasmFnCallLog, 0),
		globals:       make([]WasmGlobal, 0),
	}
	if cfg != nil {
		logger.cfg = *cfg
	}
	return logger
}

// Reset clears the data held by the logger.
func (l *WebAssemblyLogger) Reset() {
	l.storage = make(map[common.Address]Storage)
	l.output = make([]byte, 0)
	l.logs = l.logs[:0]
	l.err = nil
}

// CaptureStart implements the EVMLogger interface to initialize the tracing operation.
func (l *WebAssemblyLogger) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	l.env = env
}

func (l *WebAssemblyLogger) CaptureGlobalVariable(index uint64, op vm.OpCodeInfo, value uint64) {
	global := WasmGlobal{
		Pc:     op.Pc(),
		Index:  index,
		Op:     vm.WasmOpCodeToName[wasm.Opcode(op.Code())],
		Params: op.GetParams(),
		Value:  value,
	}
	l.globals = append(l.globals, global)
}

func (l *WebAssemblyLogger) CaptureGlobalMemoryState(globalMemory map[uint32][]byte) {
	if l.globalMemory == nil {
		l.globalMemory = make(map[uint32][]byte)
	}
	for k, v := range globalMemory {
		l.globalMemory[k] = v
	}
}

func (l *WebAssemblyLogger) CaptureWasmState(
	pc uint64,
	op vm.OpCodeInfo,
	memory *vm.MemoryChangeInfo,
	scope *vm.ScopeContext,
	depth int,
	drop,
	keep uint32,
) {
	// If tracing was interrupted, set the error and stop
	if atomic.LoadUint32(&l.interrupt) > 0 {
		l.env.Cancel()
		return
	}
	// check if already accumulated the specified number of logs
	if l.cfg.Limit != 0 && l.cfg.Limit <= len(l.logs) {
		return
	}
	stack := scope.Stack
	// Copy a snapshot of the current memory state to a new buffer
	var memData []byte
	var memOffset uint32
	var memLen uint32
	if l.cfg.EnableMemory && memory != nil {
		memData = make([]byte, len(memory.Value))
		copy(memData, memory.Value)
		memOffset = memory.Offset
		memLen = uint32(len(memory.Value))
	}
	// Copy a snapshot of the current stack state to a new buffer
	var stck []uint256.Int
	if !l.cfg.DisableStack {
		stck = make([]uint256.Int, len(stack.Data()))
		for i, item := range stack.Data() {
			stck[i] = item
		}
	}

	log := WasmLog{
		op.Pc(),
		OpCodeFamilyWASM,
		op,
		[]uint64{},
		scope.Contract.Gas,
		0,
		memData,
		memOffset,
		memLen,
		stck,
		nil,
		nil,
		depth,
		l.env.StateDB.GetRefund(),
		nil,
		keep,
		drop,
	}

	l.logs = append(l.logs, log)
}

func (l *WebAssemblyLogger) CaptureGasState(gasCost uint64, scope *vm.ScopeContext, depth int, err error) {
	// If tracing was interrupted, set the error and stop
	if atomic.LoadUint32(&l.interrupt) > 0 {
		l.env.Cancel()
		return
	}
	// last log must be a call of host function
	lastLog := l.logs[len(l.logs)-1]
	if lastLog.OpFamily != OpCodeFamilyWASM || wasm.Opcode(lastLog.Op.Code()) != wasm.Op_call {
		panic("trace order is corrupted")
	}
	// create a new snapshot of the EVM.
	log := WasmLog{
		lastLog.Pc,
		OpCodeFamilyGAS,
		GasOpCodeInfo(0),
		[]uint64{},
		scope.Contract.Gas,
		// total gas that is consumed by next WASM operations
		gasCost,
		// copy memory state from last log
		lastLog.Memory,
		lastLog.MemoryOffset,
		lastLog.MemorySize,
		// remove/replace last stack item because it contains gas spent that we copied to the gas cost section
		lastLog.Stack,
		[]byte{},
		nil,
		depth,
		lastLog.RefundCounter,
		err,
		0,
		0,
	}
	l.logs[len(l.logs)-1] = log
}

func (l *WebAssemblyLogger) CaptureWasmFunctionCall(fnIndex, maxStackHeight, numLocals uint32, fnName string) {
	l.functionCalls = append(l.functionCalls, WasmFnCallLog{
		FnIndex:        fnIndex,
		MaxStackHeight: maxStackHeight,
		NumLocals:      numLocals,
		FnName:         fnName,
	})
}

type GasOpCodeInfo byte

func (i GasOpCodeInfo) String() string {
	return "gas"
}

func (i GasOpCodeInfo) Code() byte {
	return byte(i)
}

func (i GasOpCodeInfo) GetParams() []uint64 {
	return []uint64{}
}

func (i GasOpCodeInfo) Pc() uint64 {
	return 0
}

type EvmOpCodeInfo byte

func (i EvmOpCodeInfo) String() string {
	return vm.OpCode(i).String()
}

func (i EvmOpCodeInfo) Code() byte {
	return byte(i)
}

func (i EvmOpCodeInfo) GetParams() []uint64 {
	return []uint64{}
}

func (i EvmOpCodeInfo) Pc() uint64 {
	return 0
}

// CaptureState logs a new structured log message and pushes it out to the environment
//
// CaptureState also tracks SLOAD/SSTORE ops to track storage change.
func (l *WebAssemblyLogger) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	// If tracing was interrupted, set the error and stop
	if atomic.LoadUint32(&l.interrupt) > 0 {
		l.env.Cancel()
		return
	}
	// check if already accumulated the specified number of logs
	if l.cfg.Limit != 0 && l.cfg.Limit <= len(l.logs) {
		return
	}

	memory := scope.Memory
	stack := scope.Stack
	contract := scope.Contract
	// Copy a snapshot of the current memory state to a new buffer
	var mem []byte
	if l.cfg.EnableMemory && len(memory.Data()) > 0 {
		mem = make([]byte, len(memory.Data()))
		copy(mem, memory.Data())
	}
	// Copy a snapshot of the current stack state to a new buffer
	var stck []uint256.Int
	if !l.cfg.DisableStack {
		stck = make([]uint256.Int, len(stack.Data()))
		for i, item := range stack.Data() {
			stck[i] = item
		}
	}
	stackData := stack.Data()
	stackLen := len(stackData)
	// Copy a snapshot of the current storage to a new container
	var storage Storage
	if !l.cfg.DisableStorage && (op == vm.SLOAD || op == vm.SSTORE) {
		// initialise new changed values storage container for this contract
		// if not present.
		if l.storage[contract.Address()] == nil {
			l.storage[contract.Address()] = make(Storage)
		}
		// capture SLOAD opcodes and record the read entry in the local storage
		if op == vm.SLOAD && stackLen >= 1 {
			var (
				address = common.Hash(stackData[stackLen-1].Bytes32())
				value   = l.env.StateDB.GetState(contract.Address(), address)
			)
			l.storage[contract.Address()][address] = value
			storage = l.storage[contract.Address()].Copy()
		} else if op == vm.SSTORE && stackLen >= 2 {
			// capture SSTORE opcodes and record the written entry in the local storage.
			var (
				value   = common.Hash(stackData[stackLen-2].Bytes32())
				address = common.Hash(stackData[stackLen-1].Bytes32())
			)
			l.storage[contract.Address()][address] = value
			storage = l.storage[contract.Address()].Copy()
		}
	}
	var rdata []byte
	if l.cfg.EnableReturnData {
		rdata = make([]byte, len(rData))
		copy(rdata, rData)
	}
	// last log must be a call of host function
	if len(l.logs) == 0 {
		panic("trace order is corrupted")
	}
	lastLog := l.logs[len(l.logs)-1]
	if lastLog.OpFamily != OpCodeFamilyWASM || wasm.Opcode(lastLog.Op.Code()) != wasm.Op_call {
		panic("trace order is corrupted")
	}
	l.logs = l.logs[0 : len(l.logs)-1]
	// create a new snapshot of the EVM.
	log := WasmLog{
		lastLog.Pc,
		OpCodeFamilyEVM,
		EvmOpCodeInfo(op),
		[]uint64{},
		gas,
		cost,
		// copy memory from last state
		lastLog.Memory,
		lastLog.MemoryOffset,
		lastLog.MemorySize,
		lastLog.Stack,
		rdata,
		storage,
		depth,
		l.env.StateDB.GetRefund(),
		err,
		0,
		0,
	}
	l.logs = append(l.logs, log)
}

func (l *WebAssemblyLogger) CaptureStateAfter(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
}

// CaptureFault implements the EVMLogger interface to trace an execution fault
// while running an opcode.
func (l *WebAssemblyLogger) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (l *WebAssemblyLogger) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	l.output = output
	l.err = err
	if l.cfg.Debug {
		//fmt.Printf("%#x\n", output)
		if err != nil {
			fmt.Printf(" error: %v\n", err)
		}
	}
}

func (l *WebAssemblyLogger) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

func (l *WebAssemblyLogger) CaptureExit(output []byte, gasUsed uint64, err error) {
}

func (l *WebAssemblyLogger) GetResult() (json.RawMessage, error) {
	// Tracing aborted
	if l.reason != nil {
		return nil, l.reason
	}
	failed := l.err != nil
	returnData := common.CopyBytes(l.output)
	// Return data when successful and revert reason when reverted, otherwise empty.
	returnVal := fmt.Sprintf("%x", returnData)
	if failed && l.err != vm.ErrExecutionReverted {
		returnVal = ""
	}
	globalMemory := make(map[uint32]string)
	if l.globalMemory != nil {
		for k, v := range l.globalMemory {
			globalMemory[k] = strings.TrimPrefix(hexutil.Encode(v), "0x")
		}
	}
	return json.Marshal(&WasmExecutionResult{
		Gas:           l.usedGas,
		Failed:        failed,
		GlobalMemory:  globalMemory,
		ReturnValue:   returnVal,
		StructLogs:    FormatWasmLogs(l.WasmLogs()),
		Globals:       l.globals,
		FunctionCalls: l.functionCalls,
	})
}

// Stop terminates execution of the tracer at the first opportune moment.
func (l *WebAssemblyLogger) Stop(err error) {
	l.reason = err
	atomic.StoreUint32(&l.interrupt, 1)
}

func (l *WebAssemblyLogger) CaptureTxStart(gasLimit uint64) {
	l.gasLimit = gasLimit
}

func (l *WebAssemblyLogger) CaptureTxEnd(restGas uint64) {
	l.usedGas = l.gasLimit - restGas
}

// WasmLogs returns the captured log entries.
func (l *WebAssemblyLogger) WasmLogs() []WasmLog { return l.logs }

// Error returns the VM error captured by the trace.
func (l *WebAssemblyLogger) Error() error { return l.err }

// Output returns the VM return value captured by the trace.
func (l *WebAssemblyLogger) Output() []byte { return l.output }

// WasmExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
// execution status, the amount of gas used and the return value
type WasmExecutionResult struct {
	Gas           uint64            `json:"gas"`
	InternalError string            `json:"internalError,omitempty"`
	Failed        bool              `json:"failed"`
	GlobalMemory  map[uint32]string `json:"globalMemory,omitempty"`
	ReturnValue   string            `json:"returnValue"`
	StructLogs    []WasmLogRes      `json:"structLogs"`
	Globals       []WasmGlobal      `json:"globals,omitempty"`
	FunctionCalls []WasmFnCallLog   `json:"functionCalls"`
}

// WasmLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type WasmLogRes struct {
	Pc            uint64             `json:"pc"`
	OpFamily      string             `json:"opcodeFamily"`
	Params        []uint64           `json:"params,omitempty"`
	Op            string             `json:"op"`
	Gas           uint64             `json:"gas"`
	GasCost       uint64             `json:"gasCost"`
	Depth         int                `json:"depth"`
	Error         string             `json:"error,omitempty"`
	Stack         *[]string          `json:"stack,omitempty"`
	MemoryChanges *map[uint32]string `json:"memoryChanges,omitempty"`
	Storage       *map[string]string `json:"storage,omitempty"`
	RefundCounter uint64             `json:"refund,omitempty"`
	Drop          uint32             `json:"drop,omitempty"`
}

// FormatWasmLogs formats EVM returned structured logs for json output
func FormatWasmLogs(logs []WasmLog) []WasmLogRes {
	formatted := make([]WasmLogRes, len(logs))
	for index, trace := range logs {
		formatted[index] = WasmLogRes{
			Pc:            trace.Pc,
			OpFamily:      trace.OpFamily.String(),
			Op:            trace.OpName(),
			Params:        trace.OpParams(),
			Gas:           trace.Gas,
			GasCost:       trace.GasCost,
			Depth:         trace.Depth,
			Error:         trace.ErrorString(),
			RefundCounter: trace.RefundCounter,
			Drop:          trace.Drop,
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = stackValue.Hex()
			}
			formatted[index].Stack = &stack
		}
		if trace.Memory != nil {
			mem := hexutil.Encode(trace.Memory)
			memoryChanges := make(map[uint32]string)
			memoryChanges[trace.MemoryOffset] = mem
			formatted[index].MemoryChanges = &memoryChanges
		}
		if trace.Storage != nil {
			storage := make(map[string]string)
			for i, storageValue := range trace.Storage {
				storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
			}
			formatted[index].Storage = &storage
		}
	}
	return formatted
}
