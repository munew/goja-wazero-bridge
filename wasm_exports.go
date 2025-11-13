package wbridge

import (
	"context"
	"strconv"

	"github.com/dop251/goja"
	"github.com/tetratelabs/wazero/api"
)

func resolveExports(vm *goja.Runtime, erStorage *externrefStorage, m api.Module, imports importsFn) *goja.Object {
	exports := vm.NewObject()

	exports.Set("memory", resolveMemory(vm, m))
	resolveTables(vm, m, exports)

	for name, def := range m.ExportedFunctionDefinitions() {
		if importModName, importName, isImport := def.Import(); isImport {
			jsFn := imports[importModName][importName]
			exports.Set(name, jsFn)
			continue
		}

		f := m.ExportedFunction(name)
		if f == nil {
			panic("function " + name + "not found in exports")
		}

		exports.Set(name, func(call goja.FunctionCall) goja.Value {
			params := make([]uint64, len(def.ParamTypes()))
			for i, t := range def.ParamTypes() {
				params[i] = gojaToWazero(erStorage, t, call.Argument(i))
			}

			res, err := f.Call(context.Background(), params...)
			if err != nil {
				panic(err)
			}

			switch len(def.ResultTypes()) {
				case 0:
					return goja.Undefined()
				case 1:
					return wazeroToGoja(erStorage, vm, def.ResultTypes()[0], res[0])
				default:
					results := vm.NewArray()
					for i, t := range def.ResultTypes() {
						results.Set(strconv.Itoa(i), wazeroToGoja(erStorage, vm, t, res[i]))
					}
					return results
			}
		})
	}

	return exports
}

func resolveMemory(vm *goja.Runtime, mod api.Module) *goja.Object {
	mem := mod.Memory()
	if mem == nil {
		panic("no memory table")
	}

	obj := vm.NewObject()

	// Here, we use a getter because buffer ptr changes when memory grows
	_ = obj.DefineAccessorProperty(
        "buffer",
        vm.ToValue(func(goja.FunctionCall) goja.Value {
            buf := mem.GetBuffer()
            return vm.ToValue(vm.NewArrayBuffer(buf))
        }),
        nil,
        goja.FLAG_FALSE,
        goja.FLAG_FALSE,
    )

	_ = obj.Set("grow", func(call goja.FunctionCall) goja.Value {
		_, ok := mem.Grow(uint32(call.Argument(0).ToInteger()))
		if !ok {
			panic("memory grow failed")
		}
		return nil
	})
	return obj
}

func resolveTables(vm *goja.Runtime, mod api.Module, dst *goja.Object) {
	for name, table := range mod.ExportedTables() {
		object := vm.NewObject()
		object.Set("grow", func(call goja.FunctionCall) goja.Value {
			_ = table.Grow(uint32(call.Argument(0).ToInteger()), 0)
			return nil
		})
		object.Set("get", func(call goja.FunctionCall) goja.Value {
			index := uint32(call.Argument(0).ToInteger())
			return vm.ToValue(table.Get(index))
		})
		object.Set("set", func(call goja.FunctionCall) goja.Value {
			index := uint32(call.Argument(0).ToInteger())
			value := call.Argument(1).ToInteger()
			table.Set(index, api.Reference(value))
			return nil
		})
		dst.Set(name, object)
	}
}
