import { Box, Button, Modal, Typography } from "@mui/material";
import { TextField, ThemeProvider, createTheme } from "@mui/material";
import { burn, getNFTS } from "./library/api";

import { ArtWork } from "./skinsPageArtwork";
import { BoltLoader } from "react-awesome-loaders";
import GameItem from "./components/GameItem";
import Header from "./components/Header";
import React from "react";
import { ethers } from "ethers";
import { retrieveJSON } from "./library/ipfs";
import { useState } from "react";
import { weaponItems } from "./data";
import { withStyles } from "@mui/styles";

const StyledTextField = withStyles({
	root: {
		padding: "3px",
		"& .MuiOutlinedInput-root": {
			"& fieldset": {
				borderColor: "#FFFFFF",
			},
			"&:hover fieldset": {
				borderColor: "#FFFFFF",
			},
			"&.Mui-focused fieldset": {
				borderColor: "#FFFFFF",
			},
		},
		"& .MuiFormLabel-root.Mui-disabled": {
			color: "#FFFFFF",
		},
	},
})(TextField);

function App() {
	const theme = createTheme({
		typography: {
			fontFamily: "VALORANT",
		},
	});
	const [agentArtwork, setAgentArtwork] = useState();
	const [cards, setCards] = useState(null);
	const [open, setOpen] = useState(false);
	const [wallet, setWallet] = useState(null);
	const [txData, setTxData] = useState(null);
	const [addr, setAddr] = useState(null);
	const [item, setItem] = useState(null);
	const [loading, setLoader] = useState(false);

	React.useEffect(() => {
		const selectedArt = ArtWork[(Math.random() * ArtWork.length).toFixed()];
		setAgentArtwork(selectedArt);
	}, []);

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
			<div
				style={{
					marginTop: 10,
					display: "flex",
					flexDirection: "row",
					flexWrap: "wrap",
					height: "80vh",
					overflowY: "scroll",
				}}>
				{cards.map((card) => (
					<GameItem
						imgPath={card.image}
						onCardClick={() => showModal(card)}
					/>
				))}
			</div>
		);
	};

	const ItemModal = ({ item }) => {
		return (
			<Modal
				open={open}
				style={{
					display: "flex",
					alignItems: "center",
					justifyContent: "center",
				}}>
				<Box
					style={{
						width: "800px",
						height: "600px",
						boxShadow: "0px 0px 10px #fff",
						borderRadius: 10,
						shadowRadius: 2,
						elevation: 2,
						padding: 20,
						backgroundImage: `url(${"/assets/valorant-modal.png"})`,
						backgroundSize: "cover",
						backgroundPosition: "center",
						color: "#FFFFFF",
					}}>
					{loading ? (
						<div
							style={{
								display: "flex",
								alignItems: "center",
								justifyContent: "center",
								flexDirection: "column",
								padding: 40,
							}}>
							<BoltLoader
								boltColor={"#FFFFFF"}
								backgroundBlurColor={"#E0E7FF"}
							/>
							<Typography variant="h4" style={{ marginTop: 120 }}>
								Burning your NFTs
							</Typography>
						</div>
					) : txData === null ? (
						<div>
							<div>
								<img
									src={item?.image}
									style={{
										display: "block",
										margin: "auto",
										maxWidth: 350,
										padding: 50,
									}}
								/>
								<Typography
									style={{
										color: "yellow",
										textAlign: "center",
										fontSize: 50,
										marginTop: 20,
									}}>
									{item.name}
								</Typography>
								<p
									style={{
										color: "white",
										textAlign: "center",
										fontSize: 30,
										marginTop: 20,
										fontFamily: "Roboto, sans-serif",
									}}>
									{item.description}
								</p>
							</div>
							<div
								style={{
									display: "flex",
									alignItems: "center",
									justifyContent: "center",
									marginTop: 30,
								}}>
								<Button
									style={{
										width: 200,
										height: 50,
										border: "2px solid #FFFFFF",
									}}
									onClick={() => {
										setOpen(false);
									}}>
									<Typography
										style={{ color: "#FFFFFF", fontSize: 16 }}>
										{" "}
										Nope
									</Typography>
								</Button>
								<Button
									style={{
										width: 250,
										height: 50,
										marginLeft: 20,
										border: "2px solid #FFFFFF",
									}}
									onClick={() => burnItemPress(item)}>
									<Typography
										style={{ color: "#FFFFFF", fontSize: 16 }}>
										{" "}
										Burn!
									</Typography>
								</Button>
							</div>
						</div>
					) : (
						<>
							<img
								src={item?.image}
								style={{
									display: "block",
									margin: "auto",
									maxWidth: 350,
									padding: 50,
								}}
							/>
							<Typography
								style={{
									color: "yellow",
									textAlign: "center",
									fontSize: 50,
									marginTop: 20,
								}}>
								{item.name}
							</Typography>
							<div
								style={{
									width: "100%",
									display: "flex",
									alignItems: "center",
									justifyContent: "center",
								}}>
								<StyledTextField
									label="Tx Hash"
									variant="outlined"
									focused
									value={txData.hash}
									inputProps={{
										style: {
											fontFamily: "Roboto san-serif",
											color: "#FFFFFF",
										},
									}}
									style={{ marginTop: 40 }}
									InputLabelProps={{
										style: { color: "#fff" },
									}}
								/>
								<StyledTextField
									label="Sig"
									variant="outlined"
									focused
									value={txData.sig}
									inputProps={{
										style: {
											fontFamily: "Roboto san-serif",
											color: "#FFFFFF",
										},
									}}
									style={{ marginTop: 40 }}
									InputLabelProps={{
										style: { color: "#fff" },
									}}
								/>
							</div>

							<div
								style={{
									display: "flex",
									alignItems: "center",
									justifyContent: "center",
									marginTop: 30,
								}}>
								<Button
									style={{
										width: 250,
										height: 50,
										border: "2px solid #FFFFFF",
									}}
									onClick={() => {
										setTxData(null);
										setOpen(false);
									}}>
									<Typography
										style={{ color: "#FFFFFF", fontSize: 16 }}>
										{" "}
										Close
									</Typography>
								</Button>
							</div>
						</>
					)}
				</Box>
			</Modal>
		);
	};

	const showModal = (card) => {
		setItem(card);
		setOpen(true);
	};

	return (
		<ThemeProvider theme={theme}>
			<Header>
				<div>
					<div
						style={{
							display: "flex",
							maxWidth: "1536px",
							alignItems: "center",
							justifyContent: "center",
							margin: "auto",
						}}>
						<div style={{ display: "block", marginLeft: "auto" }}>
							<Button
								style={{
									marginTop: 50,
									width: 200,
									height: 40,
									border: "2px solid #FFFFFF",
								}}
								onClick={() => refreshPress()}>
								<Typography
									style={{ color: "#FFFFFF", fontSize: 16 }}>
									{" "}
									Refresh NFTs
								</Typography>
							</Button>
							<Button
								style={{
									marginTop: 50,
									width: 250,
									height: 40,
									marginLeft: 20,
									border: "2px solid #FFFFFF",
								}}
								onClick={() => connectWalletPress()}>
								<Typography
									style={{ color: "#FFFFFF", fontSize: 16 }}>
									{" "}
									{wallet ? truncate(addr, 15) : "CONNECT WALLET"}
								</Typography>
							</Button>
						</div>
					</div>
					<div style={{ maxWidth: "1700px" }}>
						<div style={{ width: "20%" }}>
							{agentArtwork && (
								<img
									src={agentArtwork.image}
									alt={agentArtwork.name}
									style={{
										height: "90vh",
										position: "absolute",
										left: 0,
										bottom: 0,
										display: "inline-block",
										overflow: "hidden",
									}}
								/>
							)}
						</div>
						<div
							style={{
								marginLeft: "auto",
								width: "60%",
							}}>
							{cards && cards.length > 0 ? (
								<ItemGrid />
							) : (
								<Typography
									style={{
										textAlign: "center",
										fontSize: 40,
										marginTop: 200,
										width: 500,
										marginLeft: 100,
									}}>
									No Game NFTs available
								</Typography>
							)}
							{item ? <ItemModal item={item} /> : null}
						</div>
					</div>
				</div>
			</Header>
		</ThemeProvider>
	);
}
const truncate = (str, n) => {
	return str.length > n ? str.substr(0, n - 1) + "..." : str;
};
export default App;
