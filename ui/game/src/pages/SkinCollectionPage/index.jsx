import { Box, TextField } from "@mui/material";
import { Modal, Typography } from "@mui/material";
import React, { useEffect, useState } from "react";
import { mintRequest, signMintPayload } from "../../library/api";

import { ArtWork } from "./skinsPageArtwork";
import { BoltLoader } from "react-awesome-loaders";
import Button from "@mui/material/Button";
import KeyboardArrowLeft from "@mui/icons-material/KeyboardArrowLeft";
import KeyboardArrowRight from "@mui/icons-material/KeyboardArrowRight";
import MobileStepper from "@mui/material/MobileStepper";
import { upload } from "../../library/ipfs";
import { useTheme } from "@mui/material/styles";
import { withStyles } from "@mui/styles";

const MobileSlider = withStyles({
	root: {
		padding: 10,
		backgroundColor: "transparent",
	},
	dotActive: {
		backgroundColor: "#FFFFFF",
	},
})(MobileStepper);

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

const SkinsCollectionPage = ({
	skinName,
	setSkinName,
	allGameItems,
	setAllGameItems,
}) => {
	const theme = useTheme();
	const weaponName = skinName;
	const capitalize = (value) => {
		return value.charAt(0).toUpperCase() + value.slice(1);
	};
	const [mintModal, setMintModal] = useState(false);
	const [agentArtwork, setAgentArtwork] = useState();
	const [selectedItem, setSelectedItem] = useState();
	const [ethAddress, setEthAddress] = useState("");
	const gameItems = allGameItems.filter((skins) => {
		return skins?.type?.toLowerCase() === weaponName.toLowerCase();
	});
	const maxSteps = gameItems.length;
	const [loader, setLoader] = useState(false);

	const [activeStep, setActiveStep] = React.useState(0);

	const handleNext = () => {
		setActiveStep((prevActiveStep) => prevActiveStep + 1);
	};

	const handleBack = () => {
		setActiveStep((prevActiveStep) => prevActiveStep - 1);
	};

	useEffect(() => {
		const selectedArt = ArtWork[(Math.random() * ArtWork.length).toFixed()];
		setAgentArtwork(selectedArt);
	}, []);

	const mintItemPress = async (item, mintaddress) => {
		setLoader(true);
		const uri = await upload(JSON.stringify(item));
		const res = await signMintPayload("ipfs://" + uri, mintaddress).then(
			mintRequest
		);
		try {
			if (res.status === "completed") {
				const newCards = allGameItems.filter((card) => card.id !== item.id);
				setAllGameItems(newCards);
				setActiveStep(0);
				setLoader(false);
				setMintModal(false);
			}
		} catch (err) {
			console.error(err);
			setActiveStep(0);
			setLoader(false);
			setMintModal(false);
		}
	};

	return (
		<div
			style={{
				margin: "auto",
				padding: 40,
			}}>
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

			<div style={{ flexDirection: "column" }}>
				<Button
					style={{
						height: 40,
						marginLeft: "30%",
					}}
					onClick={() => {
						setSkinName("");
					}}>
					<Typography style={{ color: "#FFFFFF", fontSize: 14 }}>
						Back
					</Typography>
				</Button>
				{gameItems.length ? (
					<>
						<Typography
							variant="h2"
							style={{ marginTop: 40, textAlign: "center" }}>
							{capitalize(gameItems[activeStep].name)}
						</Typography>
						<Box
							style={{
								marginTop: 50,
								height: "fit-content",
								alignItems: "center",
								padding: 70,
								width: "90vw",
							}}>
							<img
								src={gameItems[activeStep].image}
								alt="game skin"
								style={{
									display: "block",
									width: "30%",
									margin: "auto",
								}}
							/>
							<Button
								variant="outlined"
								onClick={() => {
									setMintModal(true);
									setSelectedItem(gameItems[activeStep]);
								}}
								style={{
									display: "block",
									margin: "auto",
									marginTop: 60,
									color: "#FFFFFF",
									border: "2px solid #FFFFFF",
									fontSize: 18,
								}}>
								Mint
							</Button>
						</Box>
						<MobileSlider
							style={{
								justifyContent: "center",
								margin: "auto",
								marginTop: 20,
							}}
							variant="dots"
							steps={maxSteps}
							position="static"
							activeStep={activeStep}
							sx={{ maxWidth: 400, flexGrow: 1 }}
							nextButton={
								<Button
									size="small"
									onClick={handleNext}
									disabled={activeStep === maxSteps - 1}
									style={{
										color: "white",
										fontSize: 20,
										marginLeft: 20,
									}}>
									Next
									{theme.direction === "rtl" ? (
										<KeyboardArrowLeft />
									) : (
										<KeyboardArrowRight />
									)}
								</Button>
							}
							backButton={
								<Button
									size="small"
									onClick={handleBack}
									disabled={activeStep === 0}
									style={{
										color: "white",
										fontSize: 20,
										marginRight: 20,
									}}>
									{theme.direction === "rtl" ? (
										<KeyboardArrowRight />
									) : (
										<KeyboardArrowLeft />
									)}
									Back
								</Button>
							}
						/>
					</>
				) : (
					<Typography
						variant="h3"
						style={{ marginTop: 40, textAlign: "center" }}>
						No Items Available to Mint
					</Typography>
				)}
			</div>
			{mintModal && (
				<Modal
					open={mintModal}
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
									Minting your items
								</Typography>
							</div>
						) : (
							<>
								<Typography
									variant="h4"
									align="center"
									style={{ marginTop: 20 }}>
									Mint into NFT
								</Typography>
								<div
									style={{
										width: "100%",
										marginTop: 40,
										alignItems: "center",
										justifyContent: "center",
									}}>
									<img
										src={selectedItem.image}
										alt="skin"
										style={{
											height: 160,
											width: 400,
											display: "block",
											marginLeft: "auto",
											marginRight: "auto",
											objectFit: "contain",
										}}
									/>
									<Typography
										align="center"
										style={{
											marginTop: 30,
											fontSize: 25,
											color: "#FFFFFF",
											fontFamily: "Roboto, sans-serif",
										}}>
										<strong>
											Are you sure, you want to mint{" "}
											{selectedItem.name} ?
										</strong>
									</Typography>
									<Typography
										align="center"
										style={{
											marginTop: 30,
											fontSize: 20,
											color: "#FFFFFF",
											fontFamily: "Roboto, sans-serif",
										}}>
										Enter ETH address below:
									</Typography>
									<div
										style={{
											width: "100%",
											display: "flex",
											alignItems: "center",
											justifyContent: "center",
										}}>
										<StyledTextField
											label="ETH Address"
											variant="outlined"
											focused
											value={ethAddress}
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
												setEthAddress(e.target.value);
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
												setMintModal(false);
												setSelectedItem();
												setEthAddress("");
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
												mintItemPress(
													selectedItem,
													ethAddress
												);
											}}
											disabled={ethAddress.length === 0}>
											<Typography
												style={{
													color: "#FFFFFF",
													fontSize: 14,
												}}>
												Mint
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

export default SkinsCollectionPage;
