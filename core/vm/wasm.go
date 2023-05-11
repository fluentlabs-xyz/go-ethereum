package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/holiman/uint256"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/sys"
	zkwasm_wasmi "github.com/wasm0/zkwasm-wasmi"
	"log"
)

type WASMState struct {
	scope *ScopeContext
}

type WASMInterpreter struct {
	// input params
	evm    *EVM
	config Config
	// queue with all WASM contexts
	stateQueue []WASMState
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
	scope := in.stateQueue[in.evm.depth-1].scope
	scope.Memory = newMemoryFromSlice([]byte{}, in)
	return scope
}

func (in *WASMInterpreter) GlobalVariable(relativePc uint64, opcode api.OpCodeInfo, value uint64) {
	var wasmLogger WASMLogger
	var ok bool
	if wasmLogger, ok = in.config.Tracer.(WASMLogger); !ok {
		panic("tracer must implement [WASMLogger] in this mode")
	}
	wasmLogger.CaptureGlobalVariable(relativePc, opcode, value)
}

func (in *WASMInterpreter) BeforeState(relativePc uint64, opcode api.OpCodeInfo, stack []uint64, memory *api.MemoryChangeInfo) {
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

func (in *WASMInterpreter) AfterState(relativePc uint64, opcode api.OpCodeInfo, stack []uint64, memory *api.MemoryChangeInfo) {
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

func (in *WASMInterpreter) resizeMemory(size uint64) {
	// do honest resize memory here
	panic("not implemented")
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

	if in.config.Debug && in.config.Tracer == nil {
		panic("tracer must be configured in debug mode")
	}

	//ctx := context.TODO()

	// create scope
	stack := newstack()
	defer func() {
		returnStack(stack)
	}()

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

	stateQueue := WASMState{scope}
	defer func() {
		in.stateQueue = in.stateQueue[0 : in.evm.depth-1]
	}()
	if len(in.stateQueue) != in.evm.depth-1 {
		panic("state queue len and evm depth mismatch, this is not possible")
	}
	in.stateQueue = append(in.stateQueue, stateQueue)

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

	traceJsonBytes, err := in.wasmEngine.ComputeTrace()
	if err != nil {
		return nil, err
	}
	type traceMemory struct {
		Offset uint32 `json:"offset"`
		Len    uint32 `json:"len"`
		Data   []byte `json:"data"`
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
	}
	type traceStruct struct {
		GlobalMemory []traceMemory  `json:"global_memory"`
		Logs         []traceLog     `json:"logs"`
		FnMetas      []functionMeta `json:"fn_metas"`
	}
	trace := &traceStruct{}
	fmt.Printf("trace from wasmi: %s\n", traceJsonBytes)
	_ = json.Unmarshal(traceJsonBytes, trace)

	if in.config.Debug {
		var wasmLogger WASMLogger
		var ok bool
		if wasmLogger, ok = in.config.Tracer.(WASMLogger); !ok {
			panic("tracer must implement [WASMLogger] in this mode")
		}
		if trace.GlobalMemory != nil {
			globalMemory := make(map[uint32][]byte)
			for _, gm := range trace.GlobalMemory {
				globalMemory[gm.Offset] = gm.Data
			}
			wasmLogger.CaptureGlobalMemoryState(globalMemory)
		}
		for _, fm := range trace.FnMetas {
			wasmLogger.CaptureWasmFunctionCall(fm.FnIndex, fm.MaxStackHeight, fm.NumLocals)
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

	return in.returnData, nil
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

func (in *WASMInterpreter) execEvmOp(opcode OpCode, scope *ScopeContext) error {
	//gasCopy := scope.Contract.Gas
	memory := scope.Memory
	op := in.config.JumpTable[opcode]
	cost := op.constantGas
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
	var pc uint64
	ei := NewEVMInterpreter(in.evm, in.config)
	ei.readOnly = in.readOnly
	//if in.config.Debug {
	//	in.config.Tracer.CaptureState(math.MaxUint64, opcode, gasCopy, cost, scope, in.returnData, in.evm.depth, nil)
	//}
	ret, err := op.execute(&pc, ei, scope)
	// always copy return data, because revert opcode return data with error
	in.returnData = ret
	return err
}

func makeFnTypeWithArgValType(argsCount int, argValType api.ValueType) (result []api.ValueType) {
	result = make([]api.ValueType, 0)
	for i := 0; i < argsCount; i++ {
		result = append(result, argValType)
	}
	return result
}

func makeFnTypeI32(args int) (result []api.ValueType) {
	return makeFnTypeWithArgValType(args, api.ValueTypeI32)
}

var (
	wasmFnTypeArgs0 = makeFnTypeI32(0)
	wasmFnTypeArgs1 = makeFnTypeI32(1)
	wasmFnTypeArgs2 = makeFnTypeI32(2)
	wasmFnTypeArgs3 = makeFnTypeI32(3)
	wasmFnTypeArgs4 = makeFnTypeI32(4)
	wasmFnTypeArgs5 = makeFnTypeI32(5)
	wasmFnTypeArgs6 = makeFnTypeI32(6)
	wasmFnTypeArgs7 = makeFnTypeI32(7)
	wasmFnTypeArgs8 = makeFnTypeI32(8)
)

var wasmFunctionTypes = map[OpCode][]api.ValueType{
	STOP:           wasmFnTypeArgs0,
	SHA3:           wasmFnTypeArgs3,
	ADDRESS:        wasmFnTypeArgs1,
	BALANCE:        wasmFnTypeArgs2,
	ORIGIN:         wasmFnTypeArgs1,
	CALLER:         wasmFnTypeArgs1,
	CALLVALUE:      wasmFnTypeArgs1,
	CALLDATALOAD:   wasmFnTypeArgs2,
	CALLDATASIZE:   wasmFnTypeArgs1,
	CALLDATACOPY:   wasmFnTypeArgs3,
	CODESIZE:       wasmFnTypeArgs1,
	CODECOPY:       wasmFnTypeArgs3,
	GASPRICE:       wasmFnTypeArgs1,
	EXTCODESIZE:    wasmFnTypeArgs2,
	EXTCODECOPY:    wasmFnTypeArgs4,
	EXTCODEHASH:    wasmFnTypeArgs2,
	RETURNDATASIZE: wasmFnTypeArgs1,
	RETURNDATACOPY: wasmFnTypeArgs3,
	BLOCKHASH:      wasmFnTypeArgs2,
	COINBASE:       wasmFnTypeArgs1,
	TIMESTAMP:      wasmFnTypeArgs1,
	NUMBER:         wasmFnTypeArgs1,
	DIFFICULTY:     wasmFnTypeArgs1,
	GASLIMIT:       wasmFnTypeArgs1,
	CHAINID:        wasmFnTypeArgs1,
	SELFBALANCE:    wasmFnTypeArgs1,
	BASEFEE:        wasmFnTypeArgs1,
	SLOAD:          wasmFnTypeArgs2,
	SSTORE:         wasmFnTypeArgs2,
	PC:             wasmFnTypeArgs0,
	MSIZE:          wasmFnTypeArgs0,
	GAS:            wasmFnTypeArgs0,
	LOG0:           wasmFnTypeArgs2,
	LOG1:           wasmFnTypeArgs3,
	LOG2:           wasmFnTypeArgs4,
	LOG3:           wasmFnTypeArgs5,
	LOG4:           wasmFnTypeArgs6,
	CREATE:         wasmFnTypeArgs4,
	CALL:           wasmFnTypeArgs8,
	CALLCODE:       wasmFnTypeArgs8,
	RETURN:         wasmFnTypeArgs2,
	DELEGATECALL:   wasmFnTypeArgs7,
	CREATE2:        wasmFnTypeArgs5,
	STATICCALL:     wasmFnTypeArgs7,
	REVERT:         wasmFnTypeArgs2,
	SELFDESTRUCT:   wasmFnTypeArgs1,
}

func (in *WASMInterpreter) processOpcode(
	input []uint64,
	opcode OpCode,
	finalizer opCodeResultFn,
	inputPreprocessors ...opCodeResultFn,
) {
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
				//_ = mod.CloseWithExitCode(ctx, 1)
				panic(err2)
			}
		}
		returnStack(stack)
	}()
	// call input processors to convert memory offset to stack items for some elements
	for _, inputPreprocessor := range inputPreprocessors {
		if err := inputPreprocessor(input, scope); err != nil {
			//_ = mod.CloseWithExitCode(ctx, 1)
			panic(err)
		}
	}
	// execute EVM opcode by emulating EVM environment with gas calculation
	err := in.execEvmOp(opcode, scope)
	// TODO delete?
	if err == errStopToken {
		err = nil // clear stop token error
	}
	if err != nil {
		//_ = mod.CloseWithExitCode(ctx, 1)
		panic(err)
	}
}

type opCodeResultFn = func(input []uint64, scope *ScopeContext) error

//func (in *WASMInterpreter) defaultOpCodeHandler(
//	host wazero.HostModuleBuilder,
//	fnName string,
//	opcode OpCode,
//	finalizer opCodeResultFn,
//	inputPreprocessors ...opCodeResultFn,
//) wazero.HostModuleBuilder {
//	fnType := wasmFunctionTypes[opcode]
//	if fnType == nil {
//		log.Panicf("there is no function type for opcode (%s)", opcode.String())
//	}
//	return host.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(ctx context.Context, mod api.Module, input []uint64) {
//		// fill stack with input parameters
//		stack := newstack()
//		for i := range input {
//			stack.push(uint256.NewInt(input[len(input)-i-1]))
//		}
//		scope := in.ScopeWithStack(stack)
//		memory := in.Memory()
//		// handle function end to trigger finalizers (RETURN opcode panics errStopToken)
//		defer func() {
//			if finalizer != nil {
//				if err2 := finalizer(input, scope, memory); err2 != nil {
//					_ = mod.CloseWithExitCode(ctx, 1)
//					panic(err2)
//				}
//			}
//			returnStack(stack)
//		}()
//		// call input processors to convert memory offset to stack items for some elements
//		if inputPreprocessors != nil {
//			for _, inputPreprocessor := range inputPreprocessors {
//				if err := inputPreprocessor(input, scope, memory); err != nil {
//					_ = mod.CloseWithExitCode(ctx, 1)
//					panic(err)
//				}
//			}
//		}
//		// execute EVM opcode by emulating EVM environment with gas calculation
//		err := in.execEvmOp(opcode, scope)
//		if err != nil {
//			_ = mod.CloseWithExitCode(ctx, 1)
//			panic(err)
//		}
//	}), fnType, []api.ValueType{}).Export(fnName)
//}

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
			valueBytes, commitMemory := scope.Memory.GetPtr(int64(offset), AddressDestLen)
			itemToReplace.SetBytes(valueBytes)
			commitMemory()
		case Uint256FieldType:
			valueBytes, commitMemory := scope.Memory.GetPtr(int64(offset), Uint256DestLen)
			itemToReplace.SetBytes(valueBytes)
			commitMemory()
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
	var paramsCount = 0
	if valueTypes, ok := wasmFunctionTypes[opcode]; !ok {
		log.Panicf("failed to register fn '%s', function type not found", fnName)
	} else {
		paramsCount = len(valueTypes)
	}
	fnNameInner := fnName
	in.wasmEngine.RegisterHostFnI32(fnName, paramsCount, func(params []int32) int32 {
		if len(params) != paramsCount {
			log.Panicf("host fn '%s' called with params count %d while expected %d\n", fnNameInner, len(params), paramsCount)
		}
		input := make([]uint64, len(params))
		for i, paramValue := range params {
			input[i] = uint64(paramValue)
		}
		in.processOpcode(input, opcode, finalizer, inputPreprocessors...)
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
			Data   []byte `json:"data"`
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
		l := &traceLog{}
		_ = json.Unmarshal([]byte(jsonTrace), l)

		var memoryChange *api.MemoryChangeInfo
		if len(l.MemoryChanges) > 0 {
			if len(l.MemoryChanges) > 1 {
				panic("multiple memory changes are not supported yet")
			}
			memoryChange = &api.MemoryChangeInfo{
				Offset: l.MemoryChanges[0].Offset,
				Value:  l.MemoryChanges[0].Data,
			}
		}
		op := &wasmOpcodeInfo{
			pc:     l.Pc,
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
	in.registerNativeFunction("_evm_sload", SLOAD, copyLastStackItemToMemory(Uint256DestLen))
	in.registerNativeFunction("_evm_sstore", SSTORE, nil,
		replaceMemOffsetWithValueOnStack(0, Uint256FieldType),
		replaceMemOffsetWithValueOnStack(1, Uint256FieldType))
	// system opcodes
	in.registerNativeFunction("_evm_pc", PC, copyLastStackItemToMemory(Uint32DestLen))
	in.registerNativeFunction("_evm_msize", MSIZE, copyLastStackItemToMemory(Uint32DestLen))
	in.registerNativeFunction("_evm_gas", GAS, copyLastStackItemToMemory(Uint64DestLen))
	// log emit opcodes
	in.registerNativeFunction("_evm_log0", LOG0, nil)
	in.registerNativeFunction("_evm_log1", LOG1, nil)
	in.registerNativeFunction("_evm_log2", LOG2, nil)
	in.registerNativeFunction("_evm_log3", LOG3, nil)
	in.registerNativeFunction("_evm_log4", LOG4, nil)
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
			panic(sys.NewExitError(GasImportedFunction, errCodeGasParamsMismatch))
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
			panic(sys.NewExitError(GasImportedFunction, errCodeOutOfGas))
		}
		return int32(zkwasm_wasmi.ComputeTraceErrorCodeOk)
	})
}
