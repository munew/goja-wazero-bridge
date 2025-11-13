package wbridge

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dop251/goja"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type importsFn map[string]map[string]goja.Callable

func resolveImports(vm *goja.Runtime, erStorage *externrefStorage, rt wazero.Runtime, compiled wazero.CompiledModule, imports *goja.Object) (map[string]wazero.HostModuleBuilder, importsFn, error) {
	hostBuilders := make(map[string]wazero.HostModuleBuilder)
	fns := make(map[string]map[string]goja.Callable)

	for _, f := range compiled.ImportedFunctions() {
		modName, name, isImport := f.Import()
		if !isImport {
			continue
		}

		if _, ok := fns[modName]; !ok {
			fns[modName] = make(map[string]goja.Callable)
		}

		if _, ok := hostBuilders[modName]; !ok {
			hostBuilders[modName] = rt.NewHostModuleBuilder(modName)
		}

		modVal := imports.Get(modName)
		if modVal == nil || goja.IsNull(modVal) || goja.IsUndefined(modVal) {
			return nil, nil, fmt.Errorf("module %s not found in imports", modName)
		}

		modObject := modVal.ToObject(vm)
		importFunction := modObject.Get(name)
		callable, ok := goja.AssertFunction(importFunction)
		if !ok {
			return nil, nil, fmt.Errorf("import %s.%s is not a function", modName, name)
		}

		fns[modName][name] = callable

		params := f.ParamTypes()
		results := f.ResultTypes()
		thunk := wazeroHostThunk(vm, erStorage, callable, params, results)
		hostBuilders[modName].NewFunctionBuilder().WithGoModuleFunction(thunk, params, results).Export(name)
	}
	return hostBuilders, fns, nil
}

func wazeroHostThunk(vm *goja.Runtime, erStorage *externrefStorage, jsFn goja.Callable, params, results []api.ValueType) api.GoModuleFunction {
	return api.GoModuleFunc(func(ctx context.Context, mod api.Module, stack []uint64) {
		args := make([]goja.Value, len(params))
		for i, t := range params {
			args[i] = wazeroToGoja(erStorage, vm, t, stack[i])
		}

		res, err := jsFn(goja.Undefined(), args...)
		if err != nil {
			panic(err)
		}

		switch len(results) {
		case 0:
			return
		case 1:
			stack[0] = gojaToWazero(erStorage, results[0], res)
		default:
			resObj := res.ToObject(vm)
			for i, t := range results {
				stack[i] = gojaToWazero(erStorage, t, resObj.Get(strconv.Itoa(i)))
			}
		}
	})
}
