import { Box, Button, Modal, TextField, Typography } from "@mui/material";
import { getTokenURI, verifyBurn } from "../../library/api";

import BasicWeaponsGrid from "../../components/GameGrid";
import { BoltLoader } from "react-awesome-loaders";
import React from "react";
import SkinsCollectionPage from "../SkinCollectionPage";
import { retrieveJSON } from "../../library/ipfs";
import { weaponItems } from "../../items";
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

const style = {
	bgcolor: "background.paper",
	p: 2,
	px: 4,
	pb: 3,
};

const HomePage = () => {
	const [weaponSelect, setWeaponSelect] = React.useState("");
	const [allGameItems, setAllGameItems] = React.useState(weaponItems);
	const [transferModal, setTransferModal] = React.useState(false);
	const [transferDetails, setTransferDetails] = React.useState({
		txhash: "",
		signature: "",
	});

	const [loader, setLoader] = React.useState(false);
	const [burnTxs, setBurnTxs] = React.useState([]);
	const bringBackItemPress = async (txhash, signature) => {
		setLoader(true);
		if (burnTxs.indexOf(txhash) !== -1) {
			alert("burn already redeemed");
		} else {
			const data = await verifyBurn(txhash, signature);

			if (data.valid) {
				const uri = await getTokenURI(data.ids[0]);
				const newItem = await retrieveJSON(uri);
				newItem.id = allGameItems.length + 1;
				setAllGameItems([newItem, ...allGameItems]);
				setBurnTxs([...burnTxs, txhash]);
			} else {
				alert("signature invalid");
			}
		}
		setTransferDetails({
			txhash: "",
			signature: "",
		});
		setLoader(false);
		setTransferModal(false);
	};
	return (
		<div
			style={{
				display: "flex",
				maxWidth: "1536px",
				alignItems: "center",
				justifyContent: "center",
				margin: "auto",
			}}>
			{weaponSelect.length !== 0 ? (
				<SkinsCollectionPage
					skinName={weaponSelect}
					setSkinName={(val) => setWeaponSelect(val)}
					allGameItems={allGameItems}
					setAllGameItems={(val) => setAllGameItems(val)}
				/>
			) : (
				<>
					<div style={{ width: "5%" }}>
						<Typography
							style={{
								marginTop: 550,
								color: "transparent",
								transform: "rotate(270deg)",
								fontSize: 100,
								WebkitTextStrokeWidth: "2px",
								WebkitTextStrokeColor: "#FFFFFF",
							}}>
							Inventory
						</Typography>
					</div>

					<div style={{ width: "80%" }}>
						<div style={{ flexGrow: 1 }}>
							<Button
								style={{
									marginTop: 50,
									width: 200,
									height: 40,
									border: "2px solid #FFFFFF",
									float: "right",
								}}
								onClick={() => setTransferModal(true)}>
								<Typography
									style={{ color: "#FFFFFF", fontSize: 20 }}>
									{" "}
									TRANSFER NFT{" "}
								</Typography>
							</Button>
						</div>
						<div style={{ marginTop: 30, marginLeft: 100 }}>
							<BasicWeaponsGrid
								onWeaponSelect={(val) => setWeaponSelect(val)}
							/>
						</div>
					</div>
				</>
			)}
			{transferModal && (
				<Modal
					open={transferModal}
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
						}}
						sx={style}>
						{loader ? (
							<div
								style={{
									display: "flex",
									alignItems: "center",
									justifyContent: "center",
									flexDirection: "column",
								}}>
								<BoltLoader
									boltColor={"#FFFFFF"}
									backgroundBlurColor={"#E0E7FF"}
								/>
								<Typography variant="h4" style={{ marginTop: 120 }}>
									Transferring game item!
								</Typography>
							</div>
						) : (
							<>
								<Typography
									variant="h4"
									align="center"
									style={{ marginTop: 20 }}>
									Transfer Game Items
								</Typography>
								<div
									style={{
										width: "100%",
										marginTop: 40,
										alignItems: "center",
										justifyContent: "center",
									}}>
									<div
										style={{
											width: "100%",
											display: "flex",
											alignItems: "center",
											justifyContent: "center",
											flexDirection: "column",
										}}>
										<StyledTextField
											label="Tx Hash"
											variant="outlined"
											focused
											value={transferDetails.txhash}
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
											onChange={(e) => {
												setTransferDetails({
													txhash: e.target.value,
													signature:
														transferDetails.signature,
												});
											}}
										/>
										<StyledTextField
											label="Tx Sig"
											variant="outlined"
											focused
											value={transferDetails.signature}
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
											onChange={(e) => {
												setTransferDetails({
													txhash: transferDetails.txhash,
													signature: e.target.value,
												});
											}}
										/>
									</div>
									<div
										style={{
											width: "100%",
											display: "flex",
											alignItems: "center",
											justifyContent: "center",
											marginTop: 20,
										}}>
										<Button
											style={{
												height: 40,
												border: "2px solid #FFFFFF",
											}}
											onClick={() => {
												setTransferModal(false);
												setTransferDetails({
													txhash: "",
													signature: "",
												});
											}}>
											<Typography
												style={{
													color: "#FFFFFF",
													fontSize: 14,
												}}>
												Nope
											</Typography>
										</Button>
										<Button
											style={{
												height: 40,
												marginLeft: 20,
												border: "2px solid #FFFFFF",
											}}
											onClick={() => {
												bringBackItemPress(
													transferDetails.txhash,
													transferDetails.signature
												);
											}}>
											<Typography
												style={{
													color: "#FFFFFF",
													fontSize: 14,
												}}>
												Transfer
											</Typography>
										</Button>
									</div>
								</div>
							</>
						)}
					</Box>
				</Modal>
			)}
		</div>
	);
};

export default HomePage;
