(module
  (type (;0;) (func (param i32)))
  (type (;1;) (func (param i32 i32)))
  (type (;2;) (func (param i32) (result i32)))
  (type (;3;) (func (param i64)))
  (import "env" "_evm_address" (func (;0;) (type 0)))
  (import "env" "_evm_callvalue" (func (;1;) (type 0)))
  (import "env" "_evm_timestamp" (func (;2;) (type 0)))
  (import "env" "_evm_balance" (func (;3;) (type 1)))
  (import "env" "_evm_origin" (func (;4;) (type 0)))
  (import "env" "_evm_caller" (func (;5;) (type 0)))
  (import "env" "gas" (func (;6;) (type 3)))
  (func (;7;) (type 2) (param i32) (result i32)
    (local i32 i32 i32 i32 i32)
    i64.const 2750
    call 6
    global.get 0
    i32.const 64
    i32.sub
    local.tee 1
    global.set 0
    i32.const 0
    local.set 2
    block  ;; label = @1
      local.get 0
      i32.const 1
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      i64.const 7086
      call 6
      local.get 1
      i32.const 48
      i32.add
      local.tee 2
      i32.const 0
      i32.store
      local.get 1
      i32.const 40
      i32.add
      local.tee 3
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=32
      local.get 1
      i32.const 32
      i32.add
      call 0
      local.get 1
      i32.load8_u offset=32
      local.get 1
      i32.load8_u offset=33
      i32.add
      local.get 1
      i32.load8_u offset=34
      i32.add
      local.get 1
      i32.load8_u offset=35
      i32.add
      local.get 1
      i32.load8_u offset=36
      i32.add
      local.get 1
      i32.load8_u offset=37
      i32.add
      local.get 1
      i32.load8_u offset=38
      i32.add
      local.get 1
      i32.load8_u offset=39
      i32.add
      local.get 3
      i32.load8_u
      i32.add
      local.get 1
      i32.load8_u offset=41
      i32.add
      local.get 1
      i32.load8_u offset=42
      i32.add
      local.get 1
      i32.load8_u offset=43
      i32.add
      local.get 1
      i32.load8_u offset=44
      i32.add
      local.get 1
      i32.load8_u offset=45
      i32.add
      local.get 1
      i32.load8_u offset=46
      i32.add
      local.get 1
      i32.load8_u offset=47
      i32.add
      local.get 2
      i32.load8_u
      i32.add
      local.get 1
      i32.load8_u offset=49
      i32.add
      local.get 1
      i32.load8_u offset=50
      i32.add
      local.get 1
      i32.load8_u offset=51
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 2
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      i64.const 4409
      call 6
      local.get 1
      i32.const 32
      i32.add
      i32.const 24
      i32.add
      local.tee 3
      i64.const 0
      i64.store
      local.get 1
      i32.const 32
      i32.add
      i32.const 16
      i32.add
      local.tee 4
      i64.const 0
      i64.store
      local.get 1
      i32.const 32
      i32.add
      i32.const 8
      i32.add
      local.tee 5
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=32
      local.get 1
      i32.const 32
      i32.add
      call 1
      local.get 1
      i32.const 24
      i32.add
      local.get 3
      i64.load
      i64.store
      local.get 1
      i32.const 16
      i32.add
      local.get 4
      i64.load
      i64.store
      local.get 1
      i32.const 8
      i32.add
      local.get 5
      i64.load
      i64.store
      local.get 1
      local.get 1
      i64.load offset=32
      i64.store
      local.get 1
      global.get 3
      i32.const 8
      i32.add
      global.set 3
      global.get 3
      i32.const 1024
      i32.gt_u
      if  ;; label = @2
        unreachable
      end
      call 8
      global.get 3
      i32.const 8
      i32.sub
      global.set 3
      local.get 2
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 4
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      i64.const 1022
      call 6
      local.get 1
      i64.const 0
      i64.store offset=32
      local.get 1
      i32.const 32
      i32.add
      call 2
      local.get 2
      local.get 1
      i32.load offset=32
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 8
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      i64.const 12105
      call 6
      local.get 1
      i32.const 16
      i32.add
      i32.const 16843009
      i32.store
      local.get 1
      i32.const 8
      i32.add
      i64.const 72340172838076673
      i64.store
      local.get 1
      i64.const 72340172838076673
      i64.store
      local.get 1
      i32.const 56
      i32.add
      local.tee 3
      i64.const 0
      i64.store
      local.get 1
      i32.const 32
      i32.add
      i32.const 16
      i32.add
      local.tee 4
      i64.const 0
      i64.store
      local.get 1
      i32.const 32
      i32.add
      i32.const 8
      i32.add
      local.tee 5
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=32
      local.get 1
      local.get 1
      i32.const 32
      i32.add
      call 3
      local.get 1
      i32.load8_u offset=32
      local.get 1
      i32.load8_u offset=33
      i32.add
      local.get 1
      i32.load8_u offset=34
      i32.add
      local.get 1
      i32.load8_u offset=35
      i32.add
      local.get 1
      i32.load8_u offset=36
      i32.add
      local.get 1
      i32.load8_u offset=37
      i32.add
      local.get 1
      i32.load8_u offset=38
      i32.add
      local.get 1
      i32.load8_u offset=39
      i32.add
      local.get 5
      i32.load8_u
      i32.add
      local.get 1
      i32.load8_u offset=41
      i32.add
      local.get 1
      i32.load8_u offset=42
      i32.add
      local.get 1
      i32.load8_u offset=43
      i32.add
      local.get 1
      i32.load8_u offset=44
      i32.add
      local.get 1
      i32.load8_u offset=45
      i32.add
      local.get 1
      i32.load8_u offset=46
      i32.add
      local.get 1
      i32.load8_u offset=47
      i32.add
      local.get 4
      i32.load8_u
      i32.add
      local.get 1
      i32.load8_u offset=49
      i32.add
      local.get 1
      i32.load8_u offset=50
      i32.add
      local.get 1
      i32.load8_u offset=51
      i32.add
      local.get 1
      i32.load8_u offset=52
      i32.add
      local.get 1
      i32.load8_u offset=53
      i32.add
      local.get 1
      i32.load8_u offset=54
      i32.add
      local.get 1
      i32.load8_u offset=55
      i32.add
      local.get 3
      i32.load8_u
      i32.add
      local.get 1
      i32.load8_u offset=57
      i32.add
      local.get 1
      i32.load8_u offset=58
      i32.add
      local.get 1
      i32.load8_u offset=59
      i32.add
      local.get 1
      i32.load8_u offset=60
      i32.add
      local.get 1
      i32.load8_u offset=61
      i32.add
      local.get 1
      i32.load8_u offset=62
      i32.add
      local.get 1
      i32.load8_u offset=63
      i32.add
      local.get 2
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 16
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      i64.const 6196
      call 6
      local.get 1
      i32.const 32
      i32.add
      global.get 3
      i32.const 4
      i32.add
      global.set 3
      global.get 3
      i32.const 1024
      i32.gt_u
      if  ;; label = @2
        unreachable
      end
      call 9
      global.get 3
      i32.const 4
      i32.sub
      global.set 3
      local.get 1
      i32.load8_u offset=32
      local.get 1
      i32.load8_u offset=33
      i32.add
      local.get 1
      i32.load8_u offset=34
      i32.add
      local.get 1
      i32.load8_u offset=35
      i32.add
      local.get 1
      i32.load8_u offset=36
      i32.add
      local.get 1
      i32.load8_u offset=37
      i32.add
      local.get 1
      i32.load8_u offset=38
      i32.add
      local.get 1
      i32.load8_u offset=39
      i32.add
      local.get 1
      i32.load8_u offset=40
      i32.add
      local.get 1
      i32.load8_u offset=41
      i32.add
      local.get 1
      i32.load8_u offset=42
      i32.add
      local.get 1
      i32.load8_u offset=43
      i32.add
      local.get 1
      i32.load8_u offset=44
      i32.add
      local.get 1
      i32.load8_u offset=45
      i32.add
      local.get 1
      i32.load8_u offset=46
      i32.add
      local.get 1
      i32.load8_u offset=47
      i32.add
      local.get 1
      i32.load8_u offset=48
      i32.add
      local.get 1
      i32.load8_u offset=49
      i32.add
      local.get 1
      i32.load8_u offset=50
      i32.add
      local.get 1
      i32.load8_u offset=51
      i32.add
      local.get 2
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 32
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      i64.const 7251
      call 6
      local.get 1
      i32.const 48
      i32.add
      local.tee 0
      i32.const 0
      i32.store
      local.get 1
      i32.const 40
      i32.add
      local.tee 3
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=32
      local.get 1
      i32.const 32
      i32.add
      call 4
      local.get 1
      i32.load8_u offset=32
      local.get 1
      i32.load8_u offset=33
      i32.add
      local.get 1
      i32.load8_u offset=34
      i32.add
      local.get 1
      i32.load8_u offset=35
      i32.add
      local.get 1
      i32.load8_u offset=36
      i32.add
      local.get 1
      i32.load8_u offset=37
      i32.add
      local.get 1
      i32.load8_u offset=38
      i32.add
      local.get 1
      i32.load8_u offset=39
      i32.add
      local.get 3
      i32.load8_u
      i32.add
      local.get 1
      i32.load8_u offset=41
      i32.add
      local.get 1
      i32.load8_u offset=42
      i32.add
      local.get 1
      i32.load8_u offset=43
      i32.add
      local.get 1
      i32.load8_u offset=44
      i32.add
      local.get 1
      i32.load8_u offset=45
      i32.add
      local.get 1
      i32.load8_u offset=46
      i32.add
      local.get 1
      i32.load8_u offset=47
      i32.add
      local.get 0
      i32.load8_u
      i32.add
      local.get 1
      i32.load8_u offset=49
      i32.add
      local.get 1
      i32.load8_u offset=50
      i32.add
      local.get 1
      i32.load8_u offset=51
      i32.add
      local.get 2
      i32.add
      local.set 2
    end
    local.get 1
    i32.const 64
    i32.add
    global.set 0
    local.get 2)
  (func (;8;) (type 2) (param i32) (result i32)
    (local i32 i32 i32)
    i64.const 997
    call 6
    local.get 0
    i32.load8_u offset=31
    local.tee 1
    i32.clz
    i32.const -24
    i32.add
    local.set 2
    block  ;; label = @1
      block  ;; label = @2
        local.get 1
        i32.eqz
        br_if 0 (;@2;)
        i64.const 376
        call 6
        local.get 2
        i32.const 255
        i32.and
        local.set 1
        br 1 (;@1;)
      end
      i64.const 992
      call 6
      local.get 0
      i32.load8_u offset=30
      local.tee 1
      i32.clz
      local.get 2
      i32.add
      i32.const -24
      i32.add
      local.set 2
      block  ;; label = @2
        local.get 1
        i32.eqz
        br_if 0 (;@2;)
        i64.const 376
        call 6
        local.get 2
        i32.const 255
        i32.and
        local.set 1
        br 1 (;@1;)
      end
      i64.const 1038
      call 6
      local.get 0
      i32.load8_u offset=29
      i32.const 255
      i32.and
      local.tee 1
      i32.clz
      local.get 2
      i32.add
      i32.const -24
      i32.add
      local.set 2
      block  ;; label = @2
        local.get 1
        i32.eqz
        br_if 0 (;@2;)
        i64.const 376
        call 6
        local.get 2
        i32.const 255
        i32.and
        local.set 1
        br 1 (;@1;)
      end
      i64.const 1038
      call 6
      local.get 0
      i32.load8_u offset=28
      i32.const 255
      i32.and
      local.tee 3
      i32.clz
      local.get 2
      i32.add
      i32.const -24
      i32.add
      local.set 1
      block  ;; label = @2
        local.get 3
        i32.eqz
        br_if 0 (;@2;)
        i64.const 376
        call 6
        local.get 1
        i32.const 255
        i32.and
        local.set 1
        br 1 (;@1;)
      end
      i64.const 1038
      call 6
      local.get 0
      i32.load8_u offset=27
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      local.get 1
      i32.add
      i32.const -24
      i32.add
      local.set 1
      block  ;; label = @2
        local.get 2
        i32.eqz
        br_if 0 (;@2;)
        i64.const 376
        call 6
        local.get 1
        i32.const 255
        i32.and
        local.set 1
        br 1 (;@1;)
      end
      i64.const 1038
      call 6
      local.get 0
      i32.load8_u offset=26
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      local.get 1
      i32.add
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=25
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=24
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=23
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=22
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=21
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=20
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=19
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=18
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=17
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=16
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=15
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=14
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=13
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=12
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=11
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=10
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=9
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=8
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=7
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=6
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=5
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=4
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=3
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=2
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 1038
      call 6
      local.get 1
      local.get 0
      i32.load8_u offset=1
      i32.const 255
      i32.and
      local.tee 2
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
      local.get 2
      br_if 0 (;@1;)
      i64.const 708
      call 6
      local.get 1
      local.get 0
      i32.load8_u
      i32.const 255
      i32.and
      i32.clz
      i32.const -24
      i32.add
      i32.const 255
      i32.and
      i32.add
      local.set 1
    end
    i32.const 256
    local.get 1
    i32.sub)
  (func (;9;) (type 0) (param i32)
    i64.const 1025
    call 6
    local.get 0
    i64.const 0
    i64.store align=1
    local.get 0
    i32.const 16
    i32.add
    i32.const 0
    i32.store align=1
    local.get 0
    i32.const 8
    i32.add
    i64.const 0
    i64.store align=1
    local.get 0
    call 5)
  (func (;10;) (type 2) (param i32) (result i32)
    local.get 0
    global.get 3
    i32.const 10
    i32.add
    global.set 3
    global.get 3
    i32.const 1024
    i32.gt_u
    if  ;; label = @1
      unreachable
    end
    call 7
    global.get 3
    i32.const 10
    i32.sub
    global.set 3)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 16)
  (global (;0;) (mut i32) (i32.const 1048576))
  (global (;1;) i32 (i32.const 1048576))
  (global (;2;) i32 (i32.const 1048576))
  (global (;3;) (mut i32) (i32.const 0))
  (export "memory" (memory 0))
  (export "main" (func 10))
  (export "__data_end" (global 1))
  (export "__heap_base" (global 2)))
