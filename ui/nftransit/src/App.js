import "./App.css";

import {
  Button,
  Container,
  Grid,
  Header,
  Image,
  Modal,
} from "semantic-ui-react";
import { burn, getNFTS } from "./library/api";

import { ScatterBoxLoader } from "react-awesome-loaders";
import { ethers } from "ethers";
import { retrieveJSON } from "./library/ipfs";
import { useState } from "react";

function App() {
  const [cards, setCards] = useState(null);
  const [open, setOpen] = useState(false);
  const [wallet, setWallet] = useState(null);
  const [txData, setTxData] = useState(null);
  const [addr, setAddr] = useState(null);
  const [item, setItem] = useState(null);
  const [loading, setLoader] = useState(false);

  const burnItemPress = async (item) => {
    setLoader(true);
    const data = await burn(wallet, item.id);
    const newCards = cards.filter((card) => card.id !== item.id);
    setCards(newCards);
    setTxData(data);
    setLoader(false);
  };

  const updateTokens = async (addr) => {
    const data = await getNFTS(addr);
    const cardList = [];
    console.log(data);
    for (let id in data) {
      console.log(data[id]);
      if (data[id].startsWith("ipfs://")) {
        const meta = await retrieveJSON(data[id]);
        console.log(meta);
        meta.id = id;
        cardList.push(meta);
      }
    }
    setCards(cardList);
  };

  const connectWalletPress = async () => {
    const provider = new ethers.providers.Web3Provider(window.ethereum);
    await provider.send("eth_requestAccounts", []);
    const signer = provider.getSigner();
    const addr = await signer.getAddress();
    console.log(addr);
    setAddr(addr);
    setWallet(signer);
    await updateTokens(addr);
  };

  const refreshPress = async () => {
    if (addr !== null) {
      await updateTokens(addr);
    }
  };

  const ItemGrid = () => {
    return (
      <div className="items">
        {cards.map((card) => (
          <Image
            className="item"
            src={card.image}
            onClick={() => showModal(card)}
          />
        ))}
      </div>
    );
  };

  const LoadingBurnModal = () => {
    return (
      <Modal
        size="small"
        open={loading}
        style={{
          border: "15px solid",
          color: "midnightblue",
          borderRadius: "20px",
        }}
      >
        <Header
          size="huge"
          style={{
            backgroundColor: "#000000",
            color: "white",
          }}
        >
          BURNING...
        </Header>
        <Modal.Content
          style={{
            backgroundColor: "#000000",
            color: "white",
          }}
        >
          <div
            style={{
              paddingLeft: "115px",
              marginLeft: "115px",
              paddingBottom: "50px",
            }}
          >
            <ScatterBoxLoader primaryColor={"lime"} background={"#000000"} />
          </div>
        </Modal.Content>
      </Modal>
    );
  };

  const ItemModal = ({ item }) => {
    const stats = item.attributes.filter(
      (attr) => attr.trait_type === "stats"
    )[0].value;
    const level = item.attributes.filter(
      (attr) => attr.trait_type === "level"
    )[0].value;
    console.log(item);
    return (
      <div>
        <Modal
          onClose={() => {
            setTxData(null);
            setOpen(false);
          }}
          onOpen={() => setOpen(true)}
          open={open}
          style={{
            border: "15px solid",
            color: "midnightblue",
            borderRadius: "20px",
          }}
        >
          <Header
            size="large"
            style={{
              backgroundColor: "#000000",
              color: "white",
            }}
          >
            BURN?
          </Header>
          {txData === null ? (
            <>
              <Modal.Content image style={{ backgroundColor: "#000" }}>
                <Image size="medium" src={item?.image} wrapped />
                <Modal.Description style={{ width: "100%" }}>
                  <Header size="large" style={{ color: "yellow" }}>
                    {item.name}
                  </Header>
                  <p style={{ color: "red" }}>LEVEL: {level}</p>
                  <p style={{ color: "forestgreen" }}>{stats}</p>
                  <p style={{ color: "white" }}>{item.description}</p>
                </Modal.Description>
              </Modal.Content>
              <Modal.Actions style={{ backgroundColor: "#000" }}>
                <Button
                  color="black"
                  inverted="true"
                  size="huge"
                  onClick={() => {
                    setOpen(false);
                  }}
                >
                  Nope
                </Button>
                <Button
                  content="BURN!"
                  onClick={() => burnItemPress(item)}
                  size="huge"
                  inverted="true"
                  color="green"
                />
              </Modal.Actions>
            </>
          ) : (
            <>
              <Modal.Content image style={{ backgroundColor: "#000" }}>
                <Image size="medium" src={item?.image} wrapped />
                <Modal.Description style={{ width: "100%" }}>
                  <Header size="large" style={{ color: "yellow" }}>
                    {item.name}
                  </Header>
                  <p style={{ color: "red", "overflow-wrap": "anywhere" }}>
                    Tx Hash: {txData.hash}
                  </p>
                  <p
                    style={{
                      color: "forestgreen",
                      "overflow-wrap": "anywhere",
                    }}
                  >
                    Sig : {txData.sig}
                  </p>
                </Modal.Description>
              </Modal.Content>
              <Modal.Actions style={{ backgroundColor: "#000" }}>
                <Button
                  color="black"
                  inverted="true"
                  size="huge"
                  onClick={() => {
                    setTxData(null);
                    setOpen(false);
                  }}
                >
                  Close
                </Button>
              </Modal.Actions>
            </>
          )}
        </Modal>
      </div>
    );
  };

  const showModal = (card) => {
    setItem(card);
    setOpen(true);
  };

  return (
    <Container className="center">
      <div>
        <Button
          inverted="true"
          color="green"
          size="huge"
          onClick={() => refreshPress()}
          className="refresh-button"
        >
          REFRESH
        </Button>
        <Button
          inverted="true"
          color="green"
          size="huge"
          onClick={() => connectWalletPress()}
          className="transfer-button"
        >
          {wallet ? truncate(addr, 15) : "CONNECT WALLET"}
        </Button>
        <Header
          size="huge"
          style={{
            backgroundColor: "#000000",
            color: "white",
            marginBottom: "40px",
          }}
        >
          NFTransit
        </Header>
        {cards ? <ItemGrid /> : null}
        {item ? <ItemModal item={item} /> : null}
        <LoadingBurnModal />
      </div>
    </Container>
  );
}
const truncate = (str, n) => {
  return str.length > n ? str.substr(0, n - 1) + "..." : str;
};
export default App;
