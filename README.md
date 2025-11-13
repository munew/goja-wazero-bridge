# Bridge between goja and wazero
This project implements WebAssembly JS API in goja using wazero (a WASM runtime in pure-Go).  
This is currently more a working PoC than an actual prod-ready project.  
Feel free to open a PR for optimizations/bug fixes/implementation of furthermore WASM API functions.

# Example
```go
import (
    "context"
    "github.com/dop251/goja"
    "github.com/tetratelabs/wazero"
    wbridge "github.com/munew/goja-wazero-bridge"
)

func main() {
    vm := goja.New()
    r := wazero.NewRuntime(context.Background())
    wbridge.Install(vm, r)

    // your code goes here....
}
```

# Warning
This project uses my own fork of wazero, due to the original one not exposing table API.  
This means that you must include the 'replace ...' line (last one in go.mod) in every projects using this repo.  
If you want to use the original repo, feel free to open a PR in my fork to add all the comments :)