# use default go compiler to build wasm -- note this tends to be a lot bigger
GOOS=js GOARCH=wasm go build -o wasm/bin/wasm.wasm ./wasm/*.go
