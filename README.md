# Bridge between goja and wazero
This project implements WebAssembly JS API in goja using wazero (a WASM runtime in pure-Go).  
This is currently more a working PoC than an actual prod-ready project.  
Feel free to open a PR for optimizations/bug fixes/implementation of furthermore WASM API functions.

# Warning
This project uses my own fork of wazero, due to the original one not exposing table API.  
If you want to use the original repo, feel free to open a PR in my fork to add all the comments :)