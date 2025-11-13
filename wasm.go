package wbridge

import (
	"context"

	"github.com/dop251/goja"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type wasmModule struct {
	rt     wazero.Runtime
	module wazero.CompiledModule
}

type wasmInstance struct {
	mod     api.Module
	module  *wasmModule
}

func setupWebAssemblyObject(vm *goja.Runtime, rt wazero.Runtime) {
	wasm := vm.NewObject()
	externrefStorage := newExternrefStorage()

	_ = wasm.Set("Module", func(call goja.ConstructorCall) *goja.Object {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("WebAssembly.Module: expected bytes"))
		}
		bytes, err := bytesFromJS(vm, call.Argument(0))
		if err != nil {
			panic(vm.ToValue("WebAssembly.Module: " + err.Error()))
		}
		mod, err := rt.CompileModule(context.Background(), bytes)
		if err != nil {
			panic(vm.ToValue("WebAssembly.Module: " + err.Error()))
		}

		jsMod := &wasmModule{
			rt:     rt,
			module: mod,
		}

		o := vm.NewObject()
		_ = o.Set("__go_mod", jsMod)
		return o
	})

	_ = wasm.Set("Instance", func(call goja.ConstructorCall) *goja.Object {
		if len(call.Arguments) < 1 {
			panic(vm.ToValue("WebAssembly.Instance: expected module"))
		}
		rawMod := call.Argument(0).ToObject(vm).Get("__go_mod")
		if rawMod == nil || goja.IsNull(rawMod) || goja.IsUndefined(rawMod) {
			panic(vm.ToValue("WebAssembly.Instance: expected module"))
		}
		m := rawMod.Export().(*wasmModule)

		var importsObj *goja.Object
		if len(call.Arguments) > 1 {
			importsObj = call.Argument(1).ToObject(vm)
		}

		hostBuilders, fns, err := resolveImports(vm, externrefStorage, rt, m.module, importsObj)
		if err != nil {
			panic(vm.ToValue("WebAssembly.Instance: " + err.Error()))
		}

		for modName, hb := range hostBuilders {
			if _, err := hb.Instantiate(context.Background()); err != nil {
				panic(vm.ToValue("WebAssembly.Instance host builder: " + modName + ": " + err.Error()))
			}
		}

		mod, err := m.rt.InstantiateModule(context.Background(), m.module, wazero.NewModuleConfig())
		if err != nil {
			panic(vm.ToValue("WebAssembly.Instance: " + err.Error()))
		}

		o := vm.NewObject()
		_ = o.Set("exports", resolveExports(vm, externrefStorage, mod, fns))
		_ = o.Set("module", call.Argument(0))
		_ = o.Set("__go_instance", &wasmInstance{
			mod:     mod,
			module:  m,
		})
		return o
	})

	_ = vm.Set("WebAssembly", wasm)
}
