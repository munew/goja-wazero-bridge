package wbridge

import (
	"fmt"
	"math/big"

	"github.com/dop251/goja"
	"github.com/tetratelabs/wazero/api"
)

func bytesFromJS(vm *goja.Runtime, v goja.Value) ([]byte, error) {
	if b, ok := v.Export().([]byte); ok {
		return b, nil
	}
	if obj := v.ToObject(vm); obj != nil {
		if bytes, ok := obj.Export().([]byte); ok {
			return bytes, nil
		}
	}
	return nil, fmt.Errorf("expected bytes")
}

func wazeroToGoja(erStorage *externrefStorage, vm *goja.Runtime, t api.ValueType, raw uint64) goja.Value {
	switch t {
	case api.ValueTypeExternref:
		if ref, ok := erStorage.Get(raw); ok {
			return ref
		}
		return goja.Undefined()
	case api.ValueTypeF32:
		return vm.ToValue(api.DecodeF32(raw))
	case api.ValueTypeF64:
		return vm.ToValue(api.DecodeF64(raw))
	case api.ValueTypeI32:
		return vm.ToValue(api.DecodeI32(raw))
	case api.ValueTypeI64:
		return vm.ToValue(big.NewInt(int64(raw)))
	}

	panic(fmt.Sprintf("unsupported type: %v", t))
}

func gojaToWazero(erStorage *externrefStorage, t api.ValueType, v goja.Value) uint64 {
	switch t {
	case api.ValueTypeExternref:
		return api.EncodeExternref(erStorage.Set(v))
	case api.ValueTypeF32:
		return api.EncodeF32(v.Export().(float32))
	case api.ValueTypeF64:
		return api.EncodeF64(v.Export().(float64))
	case api.ValueTypeI32:
		return api.EncodeI32(int32(v.ToInteger()))
	case api.ValueTypeI64:
		if exported := v.Export().(*big.Int); exported != nil {
			return api.EncodeI64(exported.Int64())
		}
		return api.EncodeI64(int64(v.Export().(int64)))
	}

	panic(fmt.Sprintf("unsupported type: %v", t))
}
