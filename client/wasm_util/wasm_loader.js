import "./wasm_exec.js";

const go = new Go();
WebAssembly.instantiateStreaming(fetch("wasm.wasm"), go.importObject)
  .then(result => {
    go.run(result.instance);
  })
  .catch(console.error);
