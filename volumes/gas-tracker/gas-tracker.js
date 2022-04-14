var web3 = require("web3").web3;

web3.currentProvider.sendAsync(
  {
    method: "debug_traceTransaction",
    params: ["deployer_address", {}],
    jsonrpc: "2.0",
    id: "2",
  },
  function (err, res) {
    if (err) console.log(err);
    console.log("PERSONAL SIGNED:" + res.result);
  }
);
