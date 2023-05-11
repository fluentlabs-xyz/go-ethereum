(module
  (type (;0;) (func (param i32 i32 i32) (result i32)))
  (type (;1;) (func (param i32 i32) (result i32)))
  (type (;2;) (func (param i32)))
  (type (;3;) (func (param i32 i32)))
  (type (;4;) (func (param i32) (result i32)))
  (type (;5;) (func (param i32) (result i64)))
  (type (;6;) (func (param i32 i32 i32 i32) (result i32)))
  (type (;7;) (func))
  (import "env" "_evm_address" (func (;0;) (type 2)))
  (import "env" "_evm_call_value" (func (;1;) (type 2)))
  (import "env" "_evm_timestamp" (func (;2;) (type 2)))
  (import "env" "_evm_balance" (func (;3;) (type 3)))
  (import "env" "_evm_caller" (func (;4;) (type 2)))
  (import "env" "_evm_origin" (func (;5;) (type 2)))
  (func (;6;) (type 4) (param i32) (result i32)
    (local i32 i32 i32 i32 i32 i32 i32)
    global.get 0
    i32.const 128
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
      i32.const 0
      local.set 2
      local.get 1
      i32.const 56
      i32.add
      i32.const 16
      i32.add
      local.tee 3
      i32.const 0
      i32.store
      local.get 1
      i32.const 64
      i32.add
      local.tee 4
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=56
      local.get 1
      i32.const 56
      i32.add
      call 0
      local.get 1
      i32.const 88
      i32.add
      i32.const 16
      i32.add
      local.get 4
      i64.load
      i64.store
      local.get 1
      i32.const 112
      i32.add
      local.get 3
      i32.load
      i32.store
      local.get 1
      i64.const 85899345920
      i64.store offset=88
      local.get 1
      local.get 1
      i64.load offset=56
      i64.store offset=96
      loop  ;; label = @2
        local.get 1
        i32.const 24
        i32.add
        local.get 1
        i32.const 88
        i32.add
        call 7
        local.get 1
        i32.load8_u offset=24
        i32.const 1
        i32.and
        i32.eqz
        br_if 1 (;@1;)
        local.get 2
        local.get 1
        i32.load8_u offset=25
        i32.const 255
        i32.and
        i32.add
        local.set 2
        br 0 (;@2;)
      end
    end
    block  ;; label = @1
      local.get 0
      i32.const 2
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      local.get 1
      i32.const 88
      i32.add
      i32.const 24
      i32.add
      i64.const 0
      i64.store
      local.get 1
      i32.const 88
      i32.add
      i32.const 16
      i32.add
      i64.const 0
      i64.store
      local.get 1
      i32.const 88
      i32.add
      i32.const 8
      i32.add
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=88
      local.get 1
      i32.const 88
      i32.add
      call 1
      local.get 1
      i32.const 56
      i32.add
      i32.const 24
      i32.add
      i64.const 0
      i64.store
      local.get 1
      i32.const 56
      i32.add
      i32.const 16
      i32.add
      i64.const 0
      i64.store
      local.get 1
      i32.const 56
      i32.add
      i32.const 8
      i32.add
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=56
      i32.const 0
      local.set 5
      block  ;; label = @2
        loop  ;; label = @3
          local.get 5
          i32.const 32
          i32.eq
          br_if 1 (;@2;)
          i32.const 1
          local.set 3
          local.get 5
          i32.const 1
          i32.add
          local.set 6
          i32.const 0
          local.set 7
          i32.const 0
          local.set 4
          block  ;; label = @4
            loop  ;; label = @5
              local.get 3
              i32.const 1
              i32.and
              i32.eqz
              br_if 1 (;@4;)
              block  ;; label = @6
                local.get 7
                local.get 5
                i32.add
                local.tee 3
                i32.const 31
                i32.gt_u
                br_if 0 (;@6;)
                local.get 1
                i32.const 88
                i32.add
                local.get 3
                i32.add
                i32.load8_u
                local.get 4
                i32.or
                local.set 4
                i32.const -1
                local.set 7
                i32.const 0
                local.set 3
                br 1 (;@5;)
              end
            end
            local.get 3
            i32.const 1049016
            call 8
            unreachable
          end
          local.get 1
          i32.const 56
          i32.add
          local.get 5
          i32.add
          local.get 4
          i32.store8
          local.get 6
          local.set 5
          br 0 (;@3;)
        end
      end
      local.get 1
      i32.const 88
      i32.add
      i32.const 24
      i32.add
      local.get 1
      i32.const 56
      i32.add
      i32.const 24
      i32.add
      i64.load
      i64.store
      local.get 1
      i32.const 88
      i32.add
      i32.const 16
      i32.add
      local.get 1
      i32.const 56
      i32.add
      i32.const 16
      i32.add
      i64.load
      i64.store
      local.get 1
      i32.const 88
      i32.add
      i32.const 8
      i32.add
      local.get 1
      i32.const 56
      i32.add
      i32.const 8
      i32.add
      i64.load
      i64.store
      local.get 1
      local.get 1
      i64.load offset=56
      i64.store offset=88
      i32.const 0
      local.set 4
      i32.const 31
      local.set 3
      block  ;; label = @2
        loop  ;; label = @3
          local.get 3
          i32.const -1
          i32.eq
          br_if 1 (;@2;)
          local.get 4
          local.get 1
          i32.const 88
          i32.add
          local.get 3
          i32.add
          i32.load8_u
          local.tee 7
          i32.clz
          i32.const -24
          i32.add
          i32.const 255
          i32.and
          i32.add
          local.set 4
          local.get 3
          i32.const -1
          i32.add
          local.set 3
          local.get 7
          i32.eqz
          br_if 0 (;@3;)
        end
      end
      local.get 2
      local.get 4
      i32.sub
      i32.const 256
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 4
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      local.get 1
      i64.const 0
      i64.store offset=88
      local.get 1
      i32.const 88
      i32.add
      call 2
      local.get 2
      local.get 1
      i32.load offset=88
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 8
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      local.get 1
      i32.const 32
      i32.add
      i32.const 16
      i32.add
      i32.const 16843009
      i32.store
      local.get 1
      i32.const 32
      i32.add
      i32.const 8
      i32.add
      i64.const 72340172838076673
      i64.store
      local.get 1
      i64.const 72340172838076673
      i64.store offset=32
      local.get 1
      i32.const 56
      i32.add
      i32.const 24
      i32.add
      local.tee 3
      i64.const 0
      i64.store
      local.get 1
      i32.const 56
      i32.add
      i32.const 16
      i32.add
      local.tee 4
      i64.const 0
      i64.store
      local.get 1
      i32.const 56
      i32.add
      i32.const 8
      i32.add
      local.tee 7
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=56
      local.get 1
      i32.const 32
      i32.add
      local.get 1
      i32.const 56
      i32.add
      call 3
      local.get 1
      i32.const 88
      i32.add
      i32.const 16
      i32.add
      local.get 7
      i64.load
      i64.store
      local.get 1
      i32.const 88
      i32.add
      i32.const 24
      i32.add
      local.get 4
      i64.load
      i64.store
      local.get 1
      i32.const 120
      i32.add
      local.get 3
      i64.load
      i64.store
      local.get 1
      i64.const 137438953472
      i64.store offset=88
      local.get 1
      local.get 1
      i64.load offset=56
      i64.store offset=96
      local.get 1
      i32.const 88
      i32.add
      i32.const 8
      i32.add
      local.set 4
      i32.const 0
      local.set 3
      block  ;; label = @2
        loop  ;; label = @3
          local.get 1
          i32.const 16
          i32.add
          local.get 1
          i32.const 88
          i32.add
          call 9
          local.get 1
          i32.load offset=16
          i32.eqz
          br_if 1 (;@2;)
          local.get 3
          local.get 4
          local.get 1
          i32.load offset=20
          i32.add
          i32.load8_u
          i32.add
          local.set 3
          br 0 (;@3;)
        end
      end
      local.get 3
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
      i32.const 0
      local.set 3
      local.get 1
      i32.const 56
      i32.add
      i32.const 16
      i32.add
      local.tee 4
      i32.const 0
      i32.store
      local.get 1
      i32.const 64
      i32.add
      local.tee 7
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=56
      local.get 1
      i32.const 56
      i32.add
      call 4
      local.get 1
      i32.const 88
      i32.add
      i32.const 16
      i32.add
      local.get 7
      i64.load
      i64.store
      local.get 1
      i32.const 112
      i32.add
      local.get 4
      i32.load
      i32.store
      local.get 1
      i64.const 85899345920
      i64.store offset=88
      local.get 1
      local.get 1
      i64.load offset=56
      i64.store offset=96
      block  ;; label = @2
        loop  ;; label = @3
          local.get 1
          i32.const 8
          i32.add
          local.get 1
          i32.const 88
          i32.add
          call 7
          local.get 1
          i32.load8_u offset=8
          i32.const 1
          i32.and
          i32.eqz
          br_if 1 (;@2;)
          local.get 3
          local.get 1
          i32.load8_u offset=9
          i32.const 255
          i32.and
          i32.add
          local.set 3
          br 0 (;@3;)
        end
      end
      local.get 3
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
      i32.const 0
      local.set 3
      local.get 1
      i32.const 56
      i32.add
      i32.const 16
      i32.add
      local.tee 4
      i32.const 0
      i32.store
      local.get 1
      i32.const 64
      i32.add
      local.tee 7
      i64.const 0
      i64.store
      local.get 1
      i64.const 0
      i64.store offset=56
      local.get 1
      i32.const 56
      i32.add
      call 5
      local.get 1
      i32.const 88
      i32.add
      i32.const 16
      i32.add
      local.get 7
      i64.load
      i64.store
      local.get 1
      i32.const 112
      i32.add
      local.get 4
      i32.load
      i32.store
      local.get 1
      i64.const 85899345920
      i64.store offset=88
      local.get 1
      local.get 1
      i64.load offset=56
      i64.store offset=96
      block  ;; label = @2
        loop  ;; label = @3
          local.get 1
          local.get 1
          i32.const 88
          i32.add
          call 7
          local.get 1
          i32.load8_u
          i32.const 1
          i32.and
          i32.eqz
          br_if 1 (;@2;)
          local.get 3
          local.get 1
          i32.load8_u offset=1
          i32.const 255
          i32.and
          i32.add
          local.set 3
          br 0 (;@3;)
        end
      end
      local.get 3
      local.get 2
      i32.add
      local.set 2
    end
    block  ;; label = @1
      local.get 0
      i32.const 64
      i32.and
      i32.eqz
      br_if 0 (;@1;)
      local.get 1
      i32.const 108
      i32.add
      i32.const 0
      i32.store
      local.get 1
      i32.const 1049032
      i32.store offset=104
      local.get 1
      i64.const 1
      i64.store offset=92 align=4
      local.get 1
      i32.const 1048596
      i32.store offset=88
      local.get 1
      i32.const 88
      i32.add
      i32.const 1048616
      call 10
      unreachable
    end
    local.get 1
    i32.const 128
    i32.add
    global.set 0
    local.get 2)
  (func (;7;) (type 3) (param i32 i32)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 2
    global.set 0
    local.get 2
    i32.const 8
    i32.add
    local.get 1
    call 9
    block  ;; label = @1
      block  ;; label = @2
        local.get 2
        i32.load offset=8
        local.tee 3
        br_if 0 (;@2;)
        br 1 (;@1;)
      end
      local.get 1
      local.get 2
      i32.load offset=12
      i32.add
      i32.const 8
      i32.add
      i32.load8_u
      local.set 1
    end
    local.get 0
    local.get 1
    i32.store8 offset=1
    local.get 0
    local.get 3
    i32.const 0
    i32.ne
    i32.store8
    local.get 2
    i32.const 16
    i32.add
    global.set 0)
  (func (;8;) (type 3) (param i32 i32)
    (local i32)
    global.get 0
    i32.const 48
    i32.sub
    local.tee 2
    global.set 0
    local.get 2
    i32.const 32
    i32.store offset=4
    local.get 2
    local.get 0
    i32.store
    local.get 2
    i32.const 28
    i32.add
    i32.const 2
    i32.store
    local.get 2
    i32.const 44
    i32.add
    i32.const 1
    i32.store
    local.get 2
    i64.const 2
    i64.store offset=12 align=4
    local.get 2
    i32.const 1048684
    i32.store offset=8
    local.get 2
    i32.const 1
    i32.store offset=36
    local.get 2
    local.get 2
    i32.const 32
    i32.add
    i32.store offset=24
    local.get 2
    local.get 2
    i32.store offset=40
    local.get 2
    local.get 2
    i32.const 4
    i32.add
    i32.store offset=32
    local.get 2
    i32.const 8
    i32.add
    local.get 1
    call 10
    unreachable)
  (func (;9;) (type 3) (param i32 i32)
    (local i32 i32)
    i32.const 0
    local.set 2
    block  ;; label = @1
      local.get 1
      i32.load offset=4
      local.get 1
      i32.load
      local.tee 3
      i32.eq
      br_if 0 (;@1;)
      i32.const 1
      local.set 2
      local.get 1
      local.get 3
      i32.const 1
      i32.add
      i32.store
    end
    local.get 0
    local.get 3
    i32.store offset=4
    local.get 0
    local.get 2
    i32.store)
  (func (;10;) (type 3) (param i32 i32)
    (local i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 2
    global.set 0
    local.get 2
    i32.const 1
    i32.store8 offset=24
    local.get 2
    local.get 1
    i32.store offset=20
    local.get 2
    local.get 0
    i32.store offset=16
    local.get 2
    i32.const 1048700
    i32.store offset=12
    local.get 2
    i32.const 1049032
    i32.store offset=8
    local.get 2
    i32.const 8
    i32.add
    call 11
    unreachable)
  (func (;11;) (type 2) (param i32)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 1
    global.set 0
    block  ;; label = @1
      local.get 0
      i32.load offset=8
      local.tee 2
      br_if 0 (;@1;)
      call 16
      unreachable
    end
    local.get 1
    local.get 0
    i32.load offset=12
    i32.store offset=8
    local.get 1
    local.get 0
    i32.store offset=4
    local.get 1
    local.get 2
    i32.store
    local.get 1
    call 18
    unreachable)
  (func (;12;) (type 2) (param i32))
  (func (;13;) (type 5) (param i32) (result i64)
    i64.const -2339565740727300315)
  (func (;14;) (type 1) (param i32 i32) (result i32)
    (local i32 i32 i64 i64 i32 i32 i32 i32 i32 i32 i32 i32 i32 i32)
    global.get 0
    i32.const 48
    i32.sub
    local.tee 2
    global.set 0
    local.get 0
    i32.load
    local.tee 3
    i64.extend_i32_u
    local.set 4
    i32.const 39
    local.set 0
    block  ;; label = @1
      block  ;; label = @2
        local.get 3
        i32.const 10000
        i32.ge_u
        br_if 0 (;@2;)
        local.get 4
        local.set 5
        br 1 (;@1;)
      end
      i32.const 39
      local.set 0
      loop  ;; label = @2
        local.get 2
        i32.const 9
        i32.add
        local.get 0
        i32.add
        local.tee 3
        i32.const -4
        i32.add
        local.get 4
        i64.const 10000
        i64.div_u
        local.tee 5
        i64.const 55536
        i64.mul
        local.get 4
        i64.add
        i32.wrap_i64
        local.tee 6
        i32.const 65535
        i32.and
        i32.const 100
        i32.div_u
        local.tee 7
        i32.const 1
        i32.shl
        i32.const 1048716
        i32.add
        i32.load16_u align=1
        i32.store16 align=1
        local.get 3
        i32.const -2
        i32.add
        local.get 7
        i32.const -100
        i32.mul
        local.get 6
        i32.add
        i32.const 65535
        i32.and
        i32.const 1
        i32.shl
        i32.const 1048716
        i32.add
        i32.load16_u align=1
        i32.store16 align=1
        local.get 0
        i32.const -4
        i32.add
        local.set 0
        local.get 4
        i64.const 99999999
        i64.gt_u
        local.set 3
        local.get 5
        local.set 4
        local.get 3
        br_if 0 (;@2;)
      end
    end
    block  ;; label = @1
      local.get 5
      i32.wrap_i64
      local.tee 3
      i32.const 99
      i32.le_u
      br_if 0 (;@1;)
      local.get 2
      i32.const 9
      i32.add
      local.get 0
      i32.const -2
      i32.add
      local.tee 0
      i32.add
      local.get 5
      i32.wrap_i64
      local.tee 6
      i32.const 65535
      i32.and
      i32.const 100
      i32.div_u
      local.tee 3
      i32.const -100
      i32.mul
      local.get 6
      i32.add
      i32.const 65535
      i32.and
      i32.const 1
      i32.shl
      i32.const 1048716
      i32.add
      i32.load16_u align=1
      i32.store16 align=1
    end
    block  ;; label = @1
      block  ;; label = @2
        local.get 3
        i32.const 10
        i32.lt_u
        br_if 0 (;@2;)
        local.get 2
        i32.const 9
        i32.add
        local.get 0
        i32.const -2
        i32.add
        local.tee 0
        i32.add
        local.get 3
        i32.const 1
        i32.shl
        i32.const 1048716
        i32.add
        i32.load16_u align=1
        i32.store16 align=1
        br 1 (;@1;)
      end
      local.get 2
      i32.const 9
      i32.add
      local.get 0
      i32.const -1
      i32.add
      local.tee 0
      i32.add
      local.get 3
      i32.const 48
      i32.add
      i32.store8
    end
    i32.const 39
    local.get 0
    i32.sub
    local.set 8
    i32.const 1
    local.set 3
    i32.const 43
    i32.const 1114112
    local.get 1
    i32.load
    local.tee 6
    i32.const 1
    i32.and
    local.tee 7
    select
    local.set 9
    local.get 6
    i32.const 29
    i32.shl
    i32.const 31
    i32.shr_s
    i32.const 1049032
    i32.and
    local.set 10
    local.get 2
    i32.const 9
    i32.add
    local.get 0
    i32.add
    local.set 11
    block  ;; label = @1
      block  ;; label = @2
        local.get 1
        i32.load offset=8
        br_if 0 (;@2;)
        local.get 1
        i32.const 24
        i32.add
        i32.load
        local.tee 0
        local.get 1
        i32.const 28
        i32.add
        i32.load
        local.tee 6
        local.get 9
        local.get 10
        call 15
        br_if 1 (;@1;)
        local.get 0
        local.get 11
        local.get 8
        local.get 6
        i32.load offset=12
        call_indirect (type 0)
        local.set 3
        br 1 (;@1;)
      end
      block  ;; label = @2
        block  ;; label = @3
          block  ;; label = @4
            block  ;; label = @5
              block  ;; label = @6
                local.get 1
                i32.const 12
                i32.add
                i32.load
                local.tee 12
                local.get 8
                local.get 7
                i32.add
                local.tee 3
                i32.le_u
                br_if 0 (;@6;)
                local.get 6
                i32.const 8
                i32.and
                br_if 4 (;@2;)
                local.get 12
                local.get 3
                i32.sub
                local.tee 3
                local.set 12
                i32.const 1
                local.get 1
                i32.load8_u offset=32
                local.tee 0
                local.get 0
                i32.const 3
                i32.eq
                select
                i32.const 3
                i32.and
                local.tee 0
                br_table 3 (;@3;) 1 (;@5;) 2 (;@4;) 3 (;@3;)
              end
              i32.const 1
              local.set 3
              local.get 1
              i32.const 24
              i32.add
              i32.load
              local.tee 0
              local.get 1
              i32.const 28
              i32.add
              i32.load
              local.tee 6
              local.get 9
              local.get 10
              call 15
              br_if 4 (;@1;)
              local.get 0
              local.get 11
              local.get 8
              local.get 6
              i32.load offset=12
              call_indirect (type 0)
              local.set 3
              br 4 (;@1;)
            end
            i32.const 0
            local.set 12
            local.get 3
            local.set 0
            br 1 (;@3;)
          end
          local.get 3
          i32.const 1
          i32.shr_u
          local.set 0
          local.get 3
          i32.const 1
          i32.add
          i32.const 1
          i32.shr_u
          local.set 12
        end
        local.get 0
        i32.const 1
        i32.add
        local.set 0
        local.get 1
        i32.const 28
        i32.add
        i32.load
        local.set 7
        local.get 1
        i32.const 24
        i32.add
        i32.load
        local.set 13
        local.get 1
        i32.load offset=4
        local.set 6
        block  ;; label = @3
          loop  ;; label = @4
            local.get 0
            i32.const -1
            i32.add
            local.tee 0
            i32.eqz
            br_if 1 (;@3;)
            local.get 13
            local.get 6
            local.get 7
            i32.load offset=16
            call_indirect (type 1)
            i32.eqz
            br_if 0 (;@4;)
          end
          i32.const 1
          local.set 3
          br 2 (;@1;)
        end
        i32.const 1
        local.set 3
        local.get 6
        i32.const 1114112
        i32.eq
        br_if 1 (;@1;)
        local.get 13
        local.get 7
        local.get 9
        local.get 10
        call 15
        br_if 1 (;@1;)
        local.get 13
        local.get 11
        local.get 8
        local.get 7
        i32.load offset=12
        call_indirect (type 0)
        br_if 1 (;@1;)
        i32.const 0
        local.set 0
        block  ;; label = @3
          loop  ;; label = @4
            block  ;; label = @5
              local.get 12
              local.get 0
              i32.ne
              br_if 0 (;@5;)
              local.get 12
              local.set 0
              br 2 (;@3;)
            end
            local.get 0
            i32.const 1
            i32.add
            local.set 0
            local.get 13
            local.get 6
            local.get 7
            i32.load offset=16
            call_indirect (type 1)
            i32.eqz
            br_if 0 (;@4;)
          end
          local.get 0
          i32.const -1
          i32.add
          local.set 0
        end
        local.get 0
        local.get 12
        i32.lt_u
        local.set 3
        br 1 (;@1;)
      end
      local.get 1
      i32.load offset=4
      local.set 14
      local.get 1
      i32.const 48
      i32.store offset=4
      local.get 1
      i32.load8_u offset=32
      local.set 15
      i32.const 1
      local.set 3
      local.get 1
      i32.const 1
      i32.store8 offset=32
      local.get 1
      i32.const 24
      i32.add
      i32.load
      local.tee 6
      local.get 1
      i32.const 28
      i32.add
      i32.load
      local.tee 13
      local.get 9
      local.get 10
      call 15
      br_if 0 (;@1;)
      local.get 0
      local.get 12
      i32.add
      local.get 7
      i32.sub
      i32.const -38
      i32.add
      local.set 0
      block  ;; label = @2
        loop  ;; label = @3
          local.get 0
          i32.const -1
          i32.add
          local.tee 0
          i32.eqz
          br_if 1 (;@2;)
          local.get 6
          i32.const 48
          local.get 13
          i32.load offset=16
          call_indirect (type 1)
          i32.eqz
          br_if 0 (;@3;)
          br 2 (;@1;)
        end
      end
      local.get 6
      local.get 11
      local.get 8
      local.get 13
      i32.load offset=12
      call_indirect (type 0)
      br_if 0 (;@1;)
      local.get 1
      local.get 15
      i32.store8 offset=32
      local.get 1
      local.get 14
      i32.store offset=4
      i32.const 0
      local.set 3
    end
    local.get 2
    i32.const 48
    i32.add
    global.set 0
    local.get 3)
  (func (;15;) (type 6) (param i32 i32 i32 i32) (result i32)
    (local i32)
    block  ;; label = @1
      block  ;; label = @2
        block  ;; label = @3
          local.get 2
          i32.const 1114112
          i32.eq
          br_if 0 (;@3;)
          i32.const 1
          local.set 4
          local.get 0
          local.get 2
          local.get 1
          i32.load offset=16
          call_indirect (type 1)
          br_if 1 (;@2;)
        end
        local.get 3
        br_if 1 (;@1;)
        i32.const 0
        local.set 4
      end
      local.get 4
      return
    end
    local.get 0
    local.get 3
    i32.const 0
    local.get 1
    i32.load offset=12
    call_indirect (type 0))
  (func (;16;) (type 7)
    (local i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 0
    global.set 0
    local.get 0
    i32.const 20
    i32.add
    i32.const 0
    i32.store
    local.get 0
    i32.const 1049032
    i32.store offset=16
    local.get 0
    i64.const 1
    i64.store offset=4 align=4
    local.get 0
    i32.const 43
    i32.store offset=28
    local.get 0
    i32.const 1049032
    i32.store offset=24
    local.get 0
    local.get 0
    i32.const 24
    i32.add
    i32.store
    local.get 0
    i32.const 1049104
    call 10
    unreachable)
  (func (;17;) (type 7)
    unreachable
    unreachable)
  (func (;18;) (type 2) (param i32)
    (local i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 1
    global.set 0
    local.get 1
    i32.const 8
    i32.add
    local.get 0
    i32.const 8
    i32.add
    i32.load
    i32.store
    local.get 1
    local.get 0
    i64.load align=4
    i64.store
    local.get 1
    call 19
    unreachable)
  (func (;19;) (type 2) (param i32)
    (local i32 i32)
    local.get 0
    i32.load
    local.tee 1
    i32.const 20
    i32.add
    i32.load
    local.set 2
    block  ;; label = @1
      block  ;; label = @2
        local.get 1
        i32.load offset=4
        br_table 0 (;@2;) 0 (;@2;) 1 (;@1;)
      end
      local.get 2
      br_if 0 (;@1;)
      local.get 0
      i32.load offset=4
      i32.load8_u offset=16
      call 20
      unreachable
    end
    local.get 0
    i32.load offset=4
    i32.load8_u offset=16
    call 20
    unreachable)
  (func (;20;) (type 2) (param i32)
    (local i32 i32)
    i32.const 0
    i32.const 0
    i32.load offset=1049124
    local.tee 1
    i32.const 1
    i32.add
    i32.store offset=1049124
    i32.const 0
    i32.const 0
    i32.load offset=1049128
    i32.const 1
    i32.add
    local.tee 2
    i32.store offset=1049128
    block  ;; label = @1
      local.get 1
      i32.const 0
      i32.lt_s
      br_if 0 (;@1;)
      local.get 2
      i32.const 2
      i32.gt_u
      br_if 0 (;@1;)
      i32.const 0
      i32.load offset=1049120
      i32.const -1
      i32.le_s
      br_if 0 (;@1;)
      local.get 2
      i32.const 1
      i32.gt_u
      br_if 0 (;@1;)
      local.get 0
      i32.eqz
      br_if 0 (;@1;)
      call 17
      unreachable
    end
    unreachable
    unreachable)
  (table (;0;) 4 4 funcref)
  (memory (;0;) 17)
  (global (;0;) (mut i32) (i32.const 1048576))
  (global (;1;) i32 (i32.const 1049132))
  (global (;2;) i32 (i32.const 1049136))
  (export "memory" (memory 0))
  (export "main" (func 6))
  (export "__data_end" (global 1))
  (export "__heap_base" (global 2))
  (elem (;0;) (i32.const 1) func 14 12 13)
  (data (;0;) (i32.const 1048576) "some panic message\00\00\00\00\10\00\12\00\00\00src/lib.rs\00\00\1c\00\10\00\0a\00\00\00;\00\00\00\05\00\00\00index out of bounds: the len is  but the index is \00\008\00\10\00 \00\00\00X\00\10\00\12\00\00\00\02\00\00\00\00\00\00\00\01\00\00\00\03\00\00\0000010203040506070809101112131415161718192021222324252627282930313233343536373839404142434445464748495051525354555657585960616263646566676869707172737475767778798081828384858687888990919293949596979899/Users/dmitry/.cargo/registry/src/github.com-1ecc6299db9ec823/fixed-bigint-0.1.9/src/fixeduint.rs\00\00\00T\01\10\00a\00\00\00K\00\00\00\1c\00\00\00called `Option::unwrap()` on a `None` valuelibrary/std/src/panicking.rs\00\f3\01\10\00\1c\00\00\00+\02\00\00\0f\00\00\00"))
