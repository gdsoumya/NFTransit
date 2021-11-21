import { ethers } from "ethers";
const config = require(`./config.json`);
const abi = require(`./abi.json`);

const domain = {
  name: config.name,
  version: config.version,
  chainId: config.chainId,
  verifyingContract: config.contract,
};

export const signMintPayload = async (uri, address) => {
  const wallet = new ethers.Wallet(
    config.pk,
    new ethers.providers.JsonRpcProvider(config.rpc)
  );
  const contract = new ethers.Contract(config.contract, abi, wallet);
  const nonce = await contract.mintNonce();
  const type = {
    MintRequestType: [
      { name: "_uris", type: "bytes32[]" },
      { name: "_tos", type: "address[]" },
      { name: "nonce", type: "uint256" },
    ],
  };
  const data = {
    _uris: [ethers.utils.id(uri)],
    _tos: [address],
    nonce: Number(nonce),
  };
  console.log(domain, type, data);
  console.log(wallet._signingKey().privateKey);
  let sig = await wallet._signTypedData(domain, type, data);
  const addr = await wallet.getAddress();
  while (ethers.utils.verifyTypedData(domain, type, data, sig) !== addr) {
    sig = await wallet._signTypedData(domain, type, data);
  }
  return { ...data, _uris: [uri], signature: sig };
};

export const signBurnPayload = async () => {
  const wallet = new ethers.Wallet(
    config.pk,
    new ethers.providers.JsonRpcProvider(config.rpc)
  );
  const contract = new ethers.Contract(config.contract, abi, wallet);
  const nonce = await contract.burnNonce();
  const type = {
    BurnRequestType: [{ name: "nonce", type: "uint256" }],
  };
  const data = {
    nonce: Number(nonce),
  };
  return await wallet._signTypedData(domain, type, data);
};

const wait = function (ms = 1000) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
};

export const mintRequest = async (data) => {
  const resp = await fetch(config.backend + "/queue_mint", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      counter: data.nonce,
      contract: config.contract,
      sender: config.name,
      uris: data._uris,
      to_addrs: data._tos,
      signature: data.signature,
    }),
  });
  await wait(2000);
  resp.json().then(console.log);
  if (resp.status === 200) {
    let result = await fetch(config.backend + "/query_mint", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        counter: data.nonce,
        contract: config.contract,
        sender: config.name,
      }),
    }).then((res) => {
      console.log(res);
      return res.json();
    });
    while (result.status === "pending") {
      await wait(5000);
      result = await fetch(config.backend + "/query_mint", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          counter: data.nonce,
          contract: config.contract,
          sender: config.name,
        }),
      }).then((res) => {
        console.log(res);
        return res.json();
      });
    }
    console.log(result);
    return result;
  }
};

export const verifyBurn = async (txHash, sig) => {
  return await fetch(config.backend + "/verify_burn", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      contract: config.contract,
      tx_hash: txHash,
      signature: sig,
    }),
  }).then((res) => res.json());
};

export const getTokenURI = async (id) => {
  const wallet = new ethers.Wallet(
    config.pk,
    new ethers.providers.JsonRpcProvider(config.rpc)
  );
  const contract = new ethers.Contract(config.contract, abi, wallet);
  return await contract.uri(id);
};
