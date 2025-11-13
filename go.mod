module github.com/munew/goja-wazero-bridge

go 1.25.4

require (
	github.com/dop251/goja v0.0.0-20251103141225-af2ceb9156d7
	github.com/tetratelabs/wazero v1.10.1
)

require (
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/google/pprof v0.0.0-20230207041349-798e818bf904 // indirect
	golang.org/x/text v0.3.8 // indirect
)

replace github.com/tetratelabs/wazero => github.com/munew/wazero v0.0.0-20251113044456-8f672ba88daf
