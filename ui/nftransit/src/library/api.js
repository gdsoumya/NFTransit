import { ethers } from "ethers";
const config = require(`./config.json`);
const abi = require(`./abi.json`);

const domain = {
  name: config.name,
  version: config.version,
  chainId: config.chainId,
  verifyingContract: config.contract,
};

export const signMintPayload = async (wallet, uri, address) => {
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

export const burn = async (wallet, id) => {
  console.log(id, typeof id);
  const contract = new ethers.Contract(config.contract, abi, wallet);
  console.log([Number(id)]);
  const tx = await contract.burn([Number(id)]);
  await tx.wait();
  let burnDetails = await fetch(config.backend + "/get_burn", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      tx_hash: tx.hash,
      contract: config.contract,
    }),
  }).then((res) => res.json());
  while (burnDetails.error !== undefined) {
    await wait(5000);
    burnDetails = await fetch(config.backend + "/get_burn", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        tx_hash: tx.hash,
        contract: config.contract,
      }),
    }).then((res) => res.json());
  }
  console.log("details ", burnDetails);
  const sig = await signBurnPayload(wallet, burnDetails.nonce);
  return { hash: tx.hash, sig };
};

export const getNFTS = async (address) => {
  return fetch(config.backend + "/user_tokens", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      address: address,
      contract: config.contract,
    }),
  }).then((res) => res.json());
};

export const signBurnPayload = async (wallet, nonce) => {
  const type = {
    BurnRequestType: [{ name: "nonce", type: "uint256" }],
  };
  const data = {
    nonce,
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
      uris: ["dgfd"], //data._uris,
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
