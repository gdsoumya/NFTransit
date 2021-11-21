import { create } from "ipfs-http-client";
const config = require(`./config.json`);

export const upload = async (data) => {
  const url = config.ipfs.uploadAddr;
  const ipfs = create({ url });
  const res = await ipfs.add(data);
  return res.path;
};

export const retrieveJSON = async (hash) => {
  hash = hash.replace("ipfs://", "");
  const url = config.ipfs.retrieveAddr + hash;
  let response = await fetch(url).then((res) => res.json());
  return response;
};

export const retrieveFile = (hash) => {
  if (hash == null) return "";
  hash = hash.replace("ipfs://", "");
  return config.ipfs.retrieveAddr + hash;
};
