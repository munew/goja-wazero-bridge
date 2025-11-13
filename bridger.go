package wbridge

import (
	"github.com/dop251/goja"
	"github.com/tetratelabs/wazero"
)

func Install(vm *goja.Runtime, rt wazero.Runtime) {
	setupWebAssemblyObject(vm, rt)
}
