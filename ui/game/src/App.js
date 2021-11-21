import "./App.css";

import {
  Button,
  Container,
  Grid,
  Header,
  Image,
  Input,
  Modal,
} from "semantic-ui-react";
import {
  getTokenURI,
  mintRequest,
  signMintPayload,
  verifyBurn,
} from "./library/api";
import { retrieveJSON, upload } from "./library/ipfs";
import { useEffect, useState } from "react";

import { ScatterBoxLoader } from "react-awesome-loaders";
import { gameItems } from "./items.js";

function App() {
  const [cards, setCards] = useState(null);
  const [open, setOpen] = useState(false);
  const [showTransfer, setTransfer] = useState(false);
  const [item, setItem] = useState(null);
  const [loading, setLoader] = useState(false);
  const [loadingBringBack, setLoadingBringBack] = useState(false);
  const [counter, setCounter] = useState(false);
  const [burnTxs, setBurnTxs] = useState([]);

  const shuffleCards = () => {
    let count = 0;
    const shuffledCards = [...gameItems]
      .sort(() => Math.random() - 0.5)
      .map((card) => ({ ...card, id: count++ }));
    setCounter(count);
    setCards(shuffledCards);
  };

  const transferItemPress = (item) => {
    setTransfer(true);
  };

  const mintItemPress = async (item, mintaddress) => {
    setOpen(false);
    setLoader(true);
    const id = item["id"];
    delete item["id"];
    const uri = await upload(JSON.stringify(item));
    const res = await signMintPayload("ipfs://" + uri, mintaddress).then(
      mintRequest
    );
    console.log(res);
    item["id"] = id;
    if (res.status === "completed") {
      const newCards = cards.filter((card) => card.id !== id);
      setCards(newCards);
    } else {
      alert(res.error);
    }
    setLoader(false);
  };

  const bringBackItemPress = async (txhash, signature) => {
    setTransfer(false);
    setLoadingBringBack(true);

    if (burnTxs.indexOf(txhash) !== -1) {
      alert("burn already redeemed");
    } else {
      const data = await verifyBurn(txhash, signature);

      if (data.valid) {
        const uri = await getTokenURI(data.ids[0]);
        const newItem = await retrieveJSON(uri);
        newItem.id = counter + 1;
        setCounter(newItem.id);
        setCards([newItem, ...cards]);
        setBurnTxs([...burnTxs, txhash]);
      } else {
        alert("signature invalid");
      }
    }
    setLoadingBringBack(false);
  };

  const ItemGrid = () => {
    return (
      <Grid centered={true}>
        <Grid.Row columns={4}>
          {cards.map((card) => (
            <Grid.Column key={card.id}>
              <Image
                className="item"
                src={card.image}
                onClick={() => showModal(card)}
              />
            </Grid.Column>
          ))}
        </Grid.Row>
      </Grid>
    );
  };

  const LoadingMintModal = () => {
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
          MINTING...
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

  const LoadingBringBackModal = () => {
    return (
      <Modal
        size="small"
        open={loadingBringBack}
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
          RETRIEVING ITEM...
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
    let value = "";
    return (
      <div>
        <Modal
          onClose={() => setOpen(false)}
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
            Mint into NFT?
          </Header>
          <Modal.Content image style={{ backgroundColor: "#000" }}>
            <Image size="medium" src={item.image} wrapped />
            <Modal.Description style={{ width: "100%" }}>
              <Header size="large" style={{ color: "yellow" }}>
                {item.name}
              </Header>
              <p style={{ color: "red" }}>LEVEL: {level}</p>
              <p style={{ color: "forestgreen" }}>{stats}</p>
              <p style={{ color: "white" }}>{item.description}</p>
              <Input
                size="large"
                placeholder="ETH Address"
                style={{
                  width: "100%",
                  paddingTop: "150px",
                  backgroundColor: "black !important",
                  color: "lime !important",
                }}
                onChange={(e) => {
                  value = e.target.value;
                }}
              />
            </Modal.Description>
          </Modal.Content>
          <Modal.Actions style={{ backgroundColor: "#000" }}>
            <Button
              color="black"
              inverted="true"
              size="huge"
              onClick={() => setOpen(false)}
            >
              Nope
            </Button>
            <Button
              content="MINT!"
              onClick={() => mintItemPress(item, value)}
              size="huge"
              inverted="true"
              color="green"
            />
          </Modal.Actions>
        </Modal>
      </div>
    );
  };

  const TransferModal = () => {
    let txHash = "",
      sig = "";
    return (
      <div>
        <Modal
          size="small"
          onClose={() => setTransfer(false)}
          onOpen={() => setTransfer(true)}
          open={showTransfer}
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
            Bring Item Back
          </Header>
          <Modal.Content image style={{ backgroundColor: "#000" }}>
            <Modal.Description style={{ width: "100%" }}>
              <Input
                size="large"
                placeholder="Tx Hash"
                style={{
                  width: "100%",
                  backgroundColor: "black !important",
                  color: "lime !important",
                }}
                onChange={(e) => {
                  txHash = e.target.value;
                }}
              />
              <Input
                size="large"
                placeholder="Signature"
                style={{
                  width: "100%",
                  paddingTop: "20px",
                  backgroundColor: "black !important",
                  color: "lime !important",
                }}
                onChange={(e) => {
                  sig = e.target.value;
                }}
              />
            </Modal.Description>
          </Modal.Content>
          <Modal.Actions style={{ backgroundColor: "#000" }}>
            <Button
              color="black"
              inverted="true"
              size="huge"
              onClick={() => setTransfer(false)}
            >
              Nope
            </Button>
            <Button
              content="TRANSFER!"
              onClick={() => bringBackItemPress(txHash, sig)}
              size="huge"
              inverted="true"
              color="green"
            />
          </Modal.Actions>
        </Modal>
      </div>
    );
  };

  const showModal = (card) => {
    setItem(card);
    setOpen(true);
    //console.log(card);
  };

  // Shuffle items
  useEffect(() => {
    shuffleCards();
  }, []);

  return (
    <Container className="center">
      <div>
        <Button
          inverted="true"
          color="green"
          size="huge"
          onClick={() => transferItemPress(item)}
          className="transfer-button"
        >
          TRANSFER NFT
        </Button>
        <Header
          as="h1"
          attached="top"
          style={{
            marginBottom: "40px",
            width: "100%",
            color: "white",
            backgroundColor: "black",
            border: "0px",
            borderRadius: "25px",
          }}
        >
          Your Inventory
        </Header>
        {cards ? <ItemGrid /> : null}
        {item ? <ItemModal item={item} /> : null}
        <LoadingMintModal />
        <TransferModal />
        <LoadingBringBackModal />
      </div>
    </Container>
  );
}

export default App;
