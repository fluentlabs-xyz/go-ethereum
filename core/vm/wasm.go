package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/holiman/uint256"
	"github.com/scroll-tech/go-ethereum/common/hexutil"
	"github.com/scroll-tech/go-ethereum/common/math"
	zkwasm_wasmi "github.com/wasm0/zkwasm-wasmi"
	"log"
	"strings"
)

type WASMInterpreter struct {
	// input params
	evm    *EVM
	config Config
	// queue with all WASM contexts
	stateQueue []*ScopeContext
	// stateless params
	readOnly   bool
	returnData []byte
	wasmEngine *zkwasm_wasmi.WasmEngine
}

func NewWASMInterpreter(
	evm *EVM,
	config Config,
) VirtualInterpreter {
	instance := &WASMInterpreter{
		evm:        evm,
		config:     config,
		wasmEngine: zkwasm_wasmi.NewWasmEngine(),
	}
	instance.registerNativeFunctions()
	instance.registerLogsCallback()
	return instance
}

func (in *WASMInterpreter) ScopeWithStack(stack *Stack) *ScopeContext {
	scope := in.Scope()
	if stack == nil {
		return scope
	}
	return &ScopeContext{Memory: scope.Memory, Stack: stack, Contract: scope.Contract}
}

func (in *WASMInterpreter) Scope() *ScopeContext {
	if in.evm.depth-1 >= len(in.stateQueue) {
		panic("context queue is empty")
	}
	scope := in.stateQueue[in.evm.depth-1]
	scope.Memory = newMemoryFromSlice([]byte{}, in)
	return scope
}

func (in *WASMInterpreter) GlobalVariable(relativePc uint64, opcode OpCodeInfo, value uint64) {
	var wasmLogger WASMLogger
	var ok bool
	if wasmLogger, ok = in.config.Tracer.(WASMLogger); !ok {
		panic("tracer must implement [WASMLogger] in this mode")
	}
	wasmLogger.CaptureGlobalVariable(relativePc, opcode, value)
}

func (in *WASMInterpreter) BeforeState(relativePc uint64, opcode OpCodeInfo, stack []uint64, memory *MemoryChangeInfo) {
	evmStack := newstack()
	defer func() {
		returnStack(evmStack)
	}()
	for i := range stack {
		evmStack.push(uint256.NewInt(stack[i]))
	}
	scope := in.ScopeWithStack(evmStack)
	var wasmLogger WASMLogger
	var ok bool
	if wasmLogger, ok = in.config.Tracer.(WASMLogger); !ok {
		panic("tracer must implement [WASMLogger] in this mode")
	}
	wasmLogger.CaptureWasmState(relativePc, opcode, memory, scope, in.evm.depth, 0, 0)
}

func (in *WASMInterpreter) AfterState(relativePc uint64, opcode OpCodeInfo, stack []uint64, memory *MemoryChangeInfo) {
}

func (in *WASMInterpreter) writeMemory(offset, size uint64, value []byte) {
	// do honest write memory here
	_ = in.wasmEngine.TraceMemoryChange(uint32(offset), uint32(size), value)
}

func (in *WASMInterpreter) readMemory(offset, size uint64) []byte {
	// do honest read memory here
	data, _ := in.wasmEngine.MemoryData()
	dataChunk := data[offset : offset+size]
	return dataChunk
}

func (in *WASMInterpreter) rawData() []byte {
	data, _ := in.wasmEngine.MemoryData()
	return data
}

func (in *WASMInterpreter) resizeMemory(size uint64) {
	// do honest resize memory here
	panic("not implemented")
}

func (in *WASMInterpreter) memorySize() uint64 {
	data, _ := in.wasmEngine.MemoryData()
	return uint64(len(data))
}

func (in *WASMInterpreter) Run(
	contract *Contract,
	input []byte,
	readOnly bool,
) (ret []byte, err error) {
	// Increment the call depth which is restricted to 1024
	in.evm.depth++
	defer func() { in.evm.depth-- }()

	// Make sure the readOnly is only set if we aren't in readOnly yet.
	// This also makes sure that the readOnly flag isn't removed for child calls.
	if readOnly && !in.readOnly {
		in.readOnly = true
		defer func() { in.readOnly = false }()
	}

	// Reset the previous call's return data. It's unimportant to preserve the old buffer
	// as every returning call will return new data anyway.
	in.returnData = nil

	// Don't bother with the execution if there's no code.
	if len(contract.Code) == 0 {
		return nil, nil
	}

	in.readOnly = readOnly

	var wasmLogger WASMLogger
	var ok bool
	if in.config.Debug && in.config.Tracer == nil {
		panic("tracer must be configured in debug mode")
	} else if wasmLogger, ok = in.config.Tracer.(WASMLogger); !ok {
		panic("tracer must implement [WASMLogger] in this mode")
	}

	//ctx := context.TODO()

	// create scope
	stack := newstack()
	defer func() {
		returnStack(stack)
	}()

	contract.Input = input

	//var runtimeConfig wazero.RuntimeConfig
	//if in.config.Debug {
	//	runtimeConfig = wazero.NewRuntimeConfigWithTracer(in)
	//} else {
	//	runtimeConfig = wazero.NewRuntimeConfigInterpreter()
	//}
	//runtime := wazero.NewRuntimeWithConfig(ctx, runtimeConfig)
	//defer func() {
	//	_ = runtime.Close(ctx)
	//}()

	// if contract deployment then check contract code is safe (do injection)
	if input == nil || len(input) == 0 {
		injectResult, err := injectGasComputationAndStackProtection(contract.Code)
		if err != nil {
			return nil, fmt.Errorf("failed to check contract deployment code: %s", err)
		}
		contract.Code = injectResult
	}

	//hostModuleBuilder := runtime.NewHostModuleBuilder("env")
	//hostModuleBuilder = in.registerNativeFunctions(hostModuleBuilder)
	//_, err = hostModuleBuilder.Instantiate(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//mod, err := runtime.InstantiateModuleFromBinary(ctx, contract.Code)
	//if err != nil {
	//	return nil, err
	//}

	// create new wasm context
	scope := &ScopeContext{Stack: stack, Contract: contract}

	//memory := mod.ExportedMemory("memory")
	//if memory != nil {
	//	scope.Memory = newMemoryFromSlice(memory.RawBuffer())
	//} else {
	//	scope.Memory = newMemoryFromSlice([]byte{})
	//}

	in.wasmEngine.SetWasmBinary(contract.Code)

	if len(in.stateQueue) != in.evm.depth-1 {
		panic("state queue len and evm depth mismatch, this is not possible")
	}
	defer func() {
		in.stateQueue = in.stateQueue[0 : in.evm.depth-1]
	}()
	in.stateQueue = append(in.stateQueue, scope)

	// capture global memory state
	//if in.config.Debug {
	//	var wasmLogger WASMLogger
	//	var ok bool
	//	if wasmLogger, ok = in.config.Tracer.(WASMLogger); !ok {
	//		panic("tracer must implement [WASMLogger] in this mode")
	//	}
	//	wasmLogger.CaptureGlobalMemoryState(mod.DataInstances())
	//}

	//main := mod.ExportedFunction("main")
	//if main == nil {
	//	return nil, ErrEntrypointNotFound
	//}
	//_, err = main.Call(ctx)
	//if err != nil {
	//	if exitErr, ok := err.(*sys.ExitError); ok {
	//		switch exitErr.ExitCode() {
	//		case errCodeGasParamsMismatch:
	//			return nil, ErrBadInputParams
	//		case errCodeOutOfGas:
	//			return nil, ErrOutOfGas
	//		}
	//		return nil, exitErr
	//	}
	//	if errors.Unwrap(err) == errStopToken {
	//		// don't return error on stop token
	//	} else {
	//		return nil, tryUnwrapError(err)
	//	}
	//}

	res, err := func() (res int32, err error) {
		defer func() {
			if err2 := recover(); err2 != nil {
				err = err2.(error)
			}
		}()
		res, err = in.wasmEngine.ComputeResult()
		if err != nil {
			return res, err
		}
		return res, nil
	}()

	if zkwasm_wasmi.ComputeTraceErrorCode(res) == zkwasm_wasmi.ComputeTraceErrorCodeOutOfGas {
		err = ErrOutOfGas
	} else if zkwasm_wasmi.ComputeTraceErrorCode(res) == zkwasm_wasmi.ComputeTraceErrorCodeExecutionReverted ||
		zkwasm_wasmi.ComputeTraceErrorCode(res) == zkwasm_wasmi.ComputeTraceErrorCodeUnknown {
		err = ErrExecutionReverted
	} else if zkwasm_wasmi.ComputeTraceErrorCode(res) == zkwasm_wasmi.ComputeTraceErrorCodeStopToken {
		err = errStopToken
	}

	type traceMemory struct {
		Offset uint32 `json:"offset"`
		Len    uint32 `json:"len"`
		Data   string `json:"data"`
	}
	type traceLog struct {
		Pc            uint32        `json:"pc"`
		SourcePc      uint32        `json:"source-pc"`
		Opcode        uint8         `json:"opcode"`
		Name          string        `json:"name"`
		StackDrop     uint32        `json:"stack_drop,omitempty"`
		StackKeep     uint32        `json:"stack_keep,omitempty"`
		Params        []uint64      `json:"params,omitempty"`
		MemoryChanges []traceMemory `json:"memory_changes,omitempty"`
		Stack         []uint64      `json:"stack,omitempty"`
	}
	type functionMeta struct {
		FnIndex        uint32 `json:"fn_index"`
		MaxStackHeight uint32 `json:"max_stack_height"`
		NumLocals      uint32 `json:"num_locals"`
		FnName         string `json:"fn_name"`
	}
	type globalVariable struct {
		Index uint64 `json:"index"`
		Value uint64 `json:"value"`
	}
	type traceStruct struct {
		GlobalMemory    []traceMemory    `json:"global_memory"`
		Logs            []traceLog       `json:"logs"`
		GlobalVariables []globalVariable `json:"global_variables"`
		FnMetas         []functionMeta   `json:"fn_metas"`
	}

	if in.config.Debug {
		traceJsonBytes, err := in.wasmEngine.DumpTrace()
		if err != nil {
			fmt.Printf("can't calculate trace")
			panic(err)
		}
		trace := &traceStruct{}
		err = json.Unmarshal(traceJsonBytes, trace)
		if err != nil {
			fmt.Printf("received bad json from wasmi")
			panic(err)
		}
		if trace.GlobalMemory != nil {
			globalMemory := make(map[uint32][]byte)
			for _, gm := range trace.GlobalMemory {
				data := gm.Data
				if !strings.HasPrefix(data, "0x") {
					data = "0x" + data
				}
				globalMemory[gm.Offset] = hexutil.MustDecode(data)
			}
			wasmLogger.CaptureGlobalMemoryState(globalMemory)
		}
		for _, fm := range trace.GlobalVariables {
			wasmLogger.CaptureGlobalVariable(fm.Index, nil, fm.Value)
		}
		for _, fm := range trace.FnMetas {
			wasmLogger.CaptureWasmFunctionCall(fm.FnIndex, fm.MaxStackHeight, fm.NumLocals, fm.FnName)
		}
	}

	// pass callback for host functions (EVM)
	// pass write/read memory ops
	// return execution trace

	// if contract deployment then inject validation bytecode
	//if input == nil || len(input) == 0 {
	//	injectResult, err := injectGasComputationAndStackProtection(in.returnData)
	//	if err != nil {
	//		return nil, fmt.Errorf("failed to check contract deployment code: %s", err)
	//	}
	//	in.returnData = injectResult
	//}

	if err == errStopToken {
		err = nil
	}

	return in.returnData, err
}

type wasmOpcodeInfo struct {
	pc     uint32
	opcode byte
	name   string
	params []uint64
}

func (i *wasmOpcodeInfo) String() string {
	return i.name
}
func (i *wasmOpcodeInfo) Code() byte {
	return i.opcode
}
func (i *wasmOpcodeInfo) GetParams() []uint64 {
	return i.params
}
func (i *wasmOpcodeInfo) Pc() uint64 {
	return uint64(i.pc)
}

func tryUnwrapError(err error) error {
	err2 := errors.Unwrap(err)
	if _, ok := evmUnwrapableErrors[err2]; ok {
		return err2
	}
	return err
}

func (in *WASMInterpreter) execEvmOp(opcode OpCode, scope *ScopeContext) (err error) {
	gasCopy := scope.Contract.Gas
	memory := scope.Memory
	op := in.config.JumpTable[opcode]
	if op == nil {
		in.config.JumpTable = newLondonInstructionSet()
		op = in.config.JumpTable[opcode]
	}
	cost := op.constantGas
	defer func() {
		err2 := err
		if err2 == errStopToken {
			err2 = nil
		}
		if in.config.Debug {
			in.config.Tracer.CaptureState(math.MaxUint64, opcode, gasCopy, cost, scope, in.returnData, in.evm.depth, err2)
		}
	}()
	if !scope.Contract.UseGas(cost) {
		return ErrOutOfGas
	}
	if op.dynamicGas != nil {
		var memorySize uint64
		if op.memorySize != nil {
			memSize, overflow := op.memorySize(scope.Stack)
			if overflow {
				return ErrGasUintOverflow
			}
			if memorySize, overflow = math.SafeMul(toWordSize(memSize), 32); overflow {
				return ErrGasUintOverflow
			}
		}
		var dynamicCost uint64
		dynamicCost, err := op.dynamicGas(in.evm, scope.Contract, scope.Stack, memory, memorySize)
		cost += dynamicCost // for tracing
		if err != nil || !scope.Contract.UseGas(dynamicCost) {
			return ErrOutOfGas
		}
	}

	pc, _ := in.wasmEngine.GetLastPc()
	pc_u64 := uint64(pc)

	ei := NewEVMInterpreter(in.evm, in.config)
	ei.readOnly = in.readOnly
	ret, err := op.execute(&pc_u64, ei, scope)
	// always copy return data, because revert opcode return data with error
	in.returnData = ret
	return err
}

var wasmFunctionTypes = map[OpCode]int{
	STOP:           0,
	SHA3:           3,
	ADDRESS:        1,
	BALANCE:        2,
	ORIGIN:         1,
	CALLER:         1,
	CALLVALUE:      1,
	CALLDATALOAD:   2,
	CALLDATASIZE:   1,
	CALLDATACOPY:   3,
	CODESIZE:       1,
	CODECOPY:       3,
	GASPRICE:       1,
	EXTCODESIZE:    2,
	EXTCODECOPY:    4,
	EXTCODEHASH:    2,
	RETURNDATASIZE: 1,
	RETURNDATACOPY: 3,
	BLOCKHASH:      2,
	COINBASE:       1,
	TIMESTAMP:      1,
	NUMBER:         1,
	DIFFICULTY:     1,
	GASLIMIT:       1,
	CHAINID:        1,
	SELFBALANCE:    1,
	BASEFEE:        1,
	SLOAD:          2,
	SSTORE:         2,
	PC:             1,
	MSIZE:          1,
	GAS:            1,
	LOG0:           2,
	LOG1:           3,
	LOG2:           4,
	LOG3:           5,
	LOG4:           6,
	CREATE:         4,
	CALL:           8,
	CALLCODE:       8,
	RETURN:         2,
	DELEGATECALL:   7,
	CREATE2:        5,
	STATICCALL:     7,
	REVERT:         2,
	SELFDESTRUCT:   1,
}

func (in *WASMInterpreter) processOpcode(
	input []uint64,
	opcode OpCode,
	finalizer opCodeResultFn,
	inputPreprocessors ...opCodeResultFn,
) error {
	// fill stack with input parameters
	stack := newstack()
	for i := range input {
		value := uint256.NewInt(input[len(input)-i-1])
		stack.push(value)
	}
	scope := in.ScopeWithStack(stack)
	// handle function end to trigger finalizers (RETURN opcode panics errStopToken)
	defer func() {
		if finalizer != nil {
			if err2 := finalizer(input, scope); err2 != nil {
				panic(err2)
			}
		}
		returnStack(stack)
	}()
	// call input processors to convert memory offset to stack items for some elements
	for _, inputPreprocessor := range inputPreprocessors {
		if err := inputPreprocessor(input, scope); err != nil {
			return err
		}
	}
	// execute EVM opcode by emulating EVM environment with gas calculation
	err := in.execEvmOp(opcode, scope)
	if err != nil {
		return err
	}
	return nil
}

type opCodeResultFn = func(input []uint64, scope *ScopeContext) error

const (
	AddressDestLen = 20
	SizeDestLen    = 4
	Uint256DestLen = 32
	Uint32DestLen  = 4
	Uint64DestLen  = 8
	HashDestLen    = 32
	BoolDestLen    = 1
)

func copyLastStackItemToMemory(destLen int) opCodeResultFn {
	return func(input []uint64, scope *ScopeContext) error {
		if len(input) == 0 {
			return fmt.Errorf("last function param must be a destination pointer")
		}
		destOffset := input[len(input)-1]
		lastBytes := scope.Stack.peek().Bytes()
		res := make([]byte, destLen)
		copy(res[destLen-len(lastBytes):], lastBytes)
		scope.Memory.Set(destOffset, uint64(len(res)), res)
		return nil
	}
}

type FieldType int32

const (
	UnknownFieldType FieldType = iota
	AddressFieldType
	Uint256FieldType
)

func replaceMemOffsetWithValueOnStack(inputOffsetIndex int, fieldType FieldType) opCodeResultFn {
	return func(input []uint64, scope *ScopeContext) error {
		stack := scope.Stack
		expectedMinLength := inputOffsetIndex + 1
		if len(input) < expectedMinLength {
			return fmt.Errorf("input length is too small (%d) to contain the expected value", len(input))
		} else if stack.len() < expectedMinLength {
			return fmt.Errorf("stack length is too small (%d) to contain the expected value", stack.len())
		} else if inputOffsetIndex < 0 {
			return fmt.Errorf("input offset index can't be less 0")
		}
		offset := input[inputOffsetIndex]
		itemToReplace := stack.Back(inputOffsetIndex)
		switch fieldType {
		case AddressFieldType:
			valueBytes := scope.Memory.GetCopy(int64(offset), AddressDestLen)
			itemToReplace.SetBytes(valueBytes)
		case Uint256FieldType:
			valueBytes := scope.Memory.GetCopy(int64(offset), Uint256DestLen)
			itemToReplace.SetBytes(valueBytes)
		default:
			return fmt.Errorf("unsupported field type provided (%d)", fieldType)
		}
		return nil
	}
}

func (in *WASMInterpreter) registerNativeFunction(
	fnName string,
	opcode OpCode,
	finalizer opCodeResultFn,
	inputPreprocessors ...opCodeResultFn,
) {
	var paramsCount int
	var ok bool
	if paramsCount, ok = wasmFunctionTypes[opcode]; !ok {
		log.Panicf("failed to register fn '%s', function type not found", fnName)
	}
	fnNameInner := fnName
	in.wasmEngine.RegisterHostFnI32(fnName, paramsCount, func(params []int32) int32 {
		if len(params) != paramsCount {
			log.Printf("host fn '%s' called with params count %d while expected %d\n", fnNameInner, len(params), paramsCount)
			return int32(zkwasm_wasmi.ComputeTraceErrorCodeUnknown)
		}
		input := make([]uint64, len(params))
		for i, paramValue := range params {
			input[i] = uint64(paramValue)
		}
		err := in.processOpcode(input, opcode, finalizer, inputPreprocessors...)
		if err == errStopToken {
			return int32(zkwasm_wasmi.ComputeTraceErrorCodeStopToken)
		} else if err != nil {
			panic(err)
		}
		return int32(zkwasm_wasmi.ComputeTraceErrorCodeOk)
	})
}

func (in *WASMInterpreter) registerLogsCallback() {
	in.wasmEngine.RegisterCallbackOnAfterItemAddedToLogs(func(jsonTrace string) {
		var ok bool
		wasmLogger, ok := in.config.Tracer.(WASMLogger)
		if !ok {
			panic("tracer must implement [WASMLogger] in this mode")
		}
		type traceMemory struct {
			Offset uint32 `json:"offset"`
			Len    uint32 `json:"len"`
			Data   string `json:"data"`
		}
		type traceLog struct {
			Pc            uint32        `json:"pc"`
			SourcePc      uint32        `json:"source_pc"`
			Opcode        uint8         `json:"opcode"`
			Name          string        `json:"name"`
			StackDrop     uint32        `json:"stack_drop,omitempty"`
			StackKeep     uint32        `json:"stack_keep,omitempty"`
			Params        []uint64      `json:"params,omitempty"`
			MemoryChanges []traceMemory `json:"memory_changes,omitempty"`
			Stack         []uint64      `json:"stack,omitempty"`
		}
		l := &traceLog{}
		_ = json.Unmarshal([]byte(jsonTrace), l)

		var memoryChange *MemoryChangeInfo
		if len(l.MemoryChanges) > 0 {
			if len(l.MemoryChanges) > 1 {
				panic("multiple memory changes are not supported yet")
			}
			data := l.MemoryChanges[0].Data
			if !strings.HasPrefix(data, "0x") {
				data = "0x" + data
			}
			memoryChange = &MemoryChangeInfo{
				Offset: l.MemoryChanges[0].Offset,
				Value:  hexutil.MustDecode(data),
			}
		}
		op := &wasmOpcodeInfo{
			pc:     l.SourcePc,
			opcode: l.Opcode,
			name:   l.Name,
			params: l.Params,
		}
		evmStack := newstack()
		for i := range l.Stack {
			evmStack.push(uint256.NewInt(l.Stack[i]))
		}
		scope := in.ScopeWithStack(evmStack)
		wasmLogger.CaptureWasmState(uint64(l.Pc), op, memoryChange, scope, in.evm.depth, l.StackDrop, l.StackKeep)
		returnStack(evmStack)
	})
}

func (in *WASMInterpreter) registerNativeFunctions() {
	in.registerNativeFunction("_evm_return", RETURN, nil)
	in.registerNativeFunction("_evm_address", ADDRESS, copyLastStackItemToMemory(AddressDestLen))
	in.registerNativeFunction("_evm_stop", STOP, nil)
	in.registerNativeFunction("_evm_keccak256", SHA3, copyLastStackItemToMemory(HashDestLen))
	in.registerNativeFunction("_evm_balance", BALANCE, copyLastStackItemToMemory(Uint256DestLen),
		replaceMemOffsetWithValueOnStack(0, AddressFieldType))
	in.registerNativeFunction("_evm_origin", ORIGIN, copyLastStackItemToMemory(AddressDestLen))
	in.registerNativeFunction("_evm_caller", CALLER, copyLastStackItemToMemory(AddressDestLen))
	in.registerNativeFunction("_evm_callvalue", CALLVALUE, copyLastStackItemToMemory(Uint256DestLen))
	in.registerNativeFunction("_evm_calldataload", CALLDATALOAD, copyLastStackItemToMemory(HashDestLen),
		replaceMemOffsetWithValueOnStack(0, Uint256FieldType))
	in.registerNativeFunction("_evm_calldatasize", CALLDATASIZE, copyLastStackItemToMemory(SizeDestLen))
	in.registerNativeFunction("_evm_calldatacopy", CALLDATACOPY, nil)
	in.registerNativeFunction("_evm_codesize", CODESIZE, copyLastStackItemToMemory(SizeDestLen))
	in.registerNativeFunction("_evm_codecopy", CODECOPY, nil)
	in.registerNativeFunction("_evm_gasprice", GASPRICE, copyLastStackItemToMemory(Uint256DestLen))
	in.registerNativeFunction("_evm_extcodesize", EXTCODESIZE, copyLastStackItemToMemory(SizeDestLen),
		replaceMemOffsetWithValueOnStack(0, AddressFieldType))
	in.registerNativeFunction("_evm_extcodecopy", EXTCODECOPY, nil,
		replaceMemOffsetWithValueOnStack(0, AddressFieldType))
	in.registerNativeFunction("_evm_extcodehash", EXTCODEHASH, copyLastStackItemToMemory(HashDestLen),
		replaceMemOffsetWithValueOnStack(0, AddressFieldType))
	in.registerNativeFunction("_evm_returndatasize", RETURNDATASIZE, copyLastStackItemToMemory(SizeDestLen))
	in.registerNativeFunction("_evm_returndatacopy", RETURNDATACOPY, nil)
	in.registerNativeFunction("_evm_blockhash", BLOCKHASH, copyLastStackItemToMemory(HashDestLen))
	in.registerNativeFunction("_evm_coinbase", COINBASE, copyLastStackItemToMemory(AddressDestLen))
	in.registerNativeFunction("_evm_timestamp", TIMESTAMP, copyLastStackItemToMemory(Uint64DestLen))
	in.registerNativeFunction("_evm_number", NUMBER, copyLastStackItemToMemory(Uint64DestLen))
	in.registerNativeFunction("_evm_difficulty", DIFFICULTY, copyLastStackItemToMemory(Uint256DestLen))
	in.registerNativeFunction("_evm_gaslimit", GASLIMIT, copyLastStackItemToMemory(Uint64DestLen))
	in.registerNativeFunction("_evm_chainid", CHAINID, copyLastStackItemToMemory(Uint256DestLen))
	in.registerNativeFunction("_evm_selfbalance", SELFBALANCE, copyLastStackItemToMemory(Uint256DestLen))
	in.registerNativeFunction("_evm_basefee", BASEFEE, copyLastStackItemToMemory(Uint256DestLen))
	// storage
	in.registerNativeFunction("_evm_sload", SLOAD, copyLastStackItemToMemory(Uint256DestLen),
		replaceMemOffsetWithValueOnStack(0, Uint256FieldType))
	in.registerNativeFunction("_evm_sstore", SSTORE, nil,
		replaceMemOffsetWithValueOnStack(0, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(1, Uint256FieldType))
	// system opcodes
	in.registerNativeFunction("_evm_pc", PC, copyLastStackItemToMemory(Uint64DestLen))
	in.registerNativeFunction("_evm_msize", MSIZE, copyLastStackItemToMemory(Uint64DestLen))
	in.registerNativeFunction("_evm_gas", GAS, copyLastStackItemToMemory(Uint64DestLen))
	// log emit opcodes
	in.registerNativeFunction("_evm_log0", LOG0, nil)
	in.registerNativeFunction("_evm_log1", LOG1, nil,
		replaceMemOffsetWithValueOnStack(2, Uint256FieldType))
	in.registerNativeFunction("_evm_log2", LOG2, nil,
		replaceMemOffsetWithValueOnStack(2, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(3, Uint256FieldType))
	in.registerNativeFunction("_evm_log3", LOG3, nil,
		replaceMemOffsetWithValueOnStack(2, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(3, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(4, Uint256FieldType))
	in.registerNativeFunction("_evm_log4", LOG4, nil,
		replaceMemOffsetWithValueOnStack(2, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(3, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(4, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(5, Uint256FieldType))
	// call & create opcodes
	in.registerNativeFunction("_evm_create", CREATE, copyLastStackItemToMemory(AddressDestLen),
		replaceMemOffsetWithValueOnStack(0, Uint256FieldType))
	in.registerNativeFunction("_evm_call", CALL, copyLastStackItemToMemory(BoolDestLen),
		replaceMemOffsetWithValueOnStack(1, AddressFieldType),
		replaceMemOffsetWithValueOnStack(2, Uint256FieldType))
	in.registerNativeFunction("_evm_callcode", CALLCODE, copyLastStackItemToMemory(BoolDestLen),
		replaceMemOffsetWithValueOnStack(1, AddressFieldType),
		replaceMemOffsetWithValueOnStack(2, Uint256FieldType))
	in.registerNativeFunction("_evm_delegatecall", DELEGATECALL, copyLastStackItemToMemory(BoolDestLen),
		replaceMemOffsetWithValueOnStack(1, Uint256FieldType))
	in.registerNativeFunction("_evm_create2", CREATE2, copyLastStackItemToMemory(AddressDestLen),
		replaceMemOffsetWithValueOnStack(0, Uint256FieldType))
	in.registerNativeFunction("_evm_staticcall", STATICCALL, copyLastStackItemToMemory(BoolDestLen),
		replaceMemOffsetWithValueOnStack(1, AddressFieldType))
	in.registerNativeFunction("_evm_revert", REVERT, nil)
	in.registerNativeFunction("_evm_selfdestruct", SELFDESTRUCT, nil)

	in.registerGasCheckFunction()
}

func (in *WASMInterpreter) registerGasCheckFunction() {
	paramsCount := 1
	in.wasmEngine.RegisterHostFnI64(GasImportedFunction, paramsCount, func(params []int64) int32 {
		if len(params) != paramsCount {
			panic(ErrOutOfGas)
		}
		scope := in.Scope()
		input := make([]uint64, len(params))
		for i, paramValue := range params {
			input[i] = uint64(paramValue)
		}
		val := int64(input[0])
		gasSpend := uint64(val)
		if in.config.Debug {
			scope := &ScopeContext{
				Contract: scope.Contract,
			}
			wasmLogger := in.config.Tracer.(WASMLogger)
			if scope.Contract.Gas < gasSpend {
				wasmLogger.CaptureGasState(gasSpend, scope, in.evm.depth, ErrOutOfGas)
				return int32(zkwasm_wasmi.ComputeTraceErrorCodeOutOfGas)
			} else {
				wasmLogger.CaptureGasState(gasSpend, scope, in.evm.depth, nil)
			}
		}
		if !scope.Contract.UseGas(gasSpend) {
			panic(ErrOutOfGas)
		}
		return int32(zkwasm_wasmi.ComputeTraceErrorCodeOk)
	})
}
