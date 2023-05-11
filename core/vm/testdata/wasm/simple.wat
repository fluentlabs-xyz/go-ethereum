(module
  (type (;0;) (func))
  (func (;0;) (type 0)
    i32.const 100
    i32.const 20
    i32.const 3
    i32.add
    i32.add
    drop)
  (memory (;0;) 1)
  (export "memory" (memory 0))
  (export "main" (func 0)))
