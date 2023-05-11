(module
  (type (;0;) (func (param i32 i32)))
  (type (;1;) (func))
  (import "env" "_evm_return" (func (;0;) (type 0)))
  (func (;1;) (type 1)
    i32.const 1048576
    i32.const 160
    call 0)
  (memory (;0;) 17)
  (global (;0;) (mut i32) (i32.const 1048576))
  (global (;1;) i32 (i32.const 1048736))
  (global (;2;) i32 (i32.const 1048736))
  (export "memory" (memory 0))
  (export "main" (func 1))
  (export "__data_end" (global 1))
  (export "__heap_base" (global 2))
  (data (;0;) (i32.const 1048576) "\00asm\01\00\00\00\01\09\02`\02\7f\7f\00`\00\00\02\13\01\03env\0b_evm_return\00\00\03\02\01\01\05\03\01\00\11\06\19\03\7f\01A\80\80\c0\00\0b\7f\00A\8c\80\c0\00\0b\7f\00A\90\80\c0\00\0b\07,\04\06memory\02\00\04main\00\01\0a__data_end\03\01\0b__heap_base\03\02\0a\0d\01\0b\00A\80\80\c0\00A\0c\10\00\0b\0b\15\01\00A\80\80\c0\00\0b\0cHello, World"))
