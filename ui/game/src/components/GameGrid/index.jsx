import * as React from "react";
import { styled } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Grid from "@mui/material/Grid";
import GameItem from "../GameItem";
import { useNavigate } from "react-router-dom";

const Item = styled(Paper)(({ theme }) => ({
	...theme.typography.body2,
	padding: theme.spacing(1),
	textAlign: "center",
	fontSize: 18,
	color: theme.palette.text.secondary,
}));

const BasicWeaponsGrid = ({ onWeaponSelect }) => {
	const history = useNavigate();
	return (
		<Box sx={{ flexGrow: 1, marginTop: 15 }}>
			<Grid container spacing={2}>
				<Grid item xs={3}>
					<Item>SIDEARMS</Item>
					<GameItem
						imgPath="https://gateway.pinata.cloud/ipfs/QmYQdDyqQqS6XhWRcsqGRejVkFzEHxACiRiPhYAXLQNJGj/Knife/Glitchpop_Dagger.png"
						onCardClick={() => {
							onWeaponSelect("knife");
							// history(`/weapons/`);
						}}
						cardName={"Knife"}
					/>
					<GameItem
						imgPath="https://gateway.pinata.cloud/ipfs/QmYQdDyqQqS6XhWRcsqGRejVkFzEHxACiRiPhYAXLQNJGj/Sheriff/Reaver_Sheriff.png"
						onCardClick={() => {
							onWeaponSelect("sheriff");
						}}
						cardName={"Sheriff"}
					/>
				</Grid>
				<Grid item xs={3}>
					<Item>SMGS</Item>
					<GameItem
						imgPath="https://gateway.pinata.cloud/ipfs/QmYQdDyqQqS6XhWRcsqGRejVkFzEHxACiRiPhYAXLQNJGj/Spectre/Forsaken_Spectre.png"
						onCardClick={() => {
							onWeaponSelect("spectre");
						}}
						cardName={"Spectre"}
					/>
				</Grid>
				<Grid item xs={3}>
					<Item>RIFLES</Item>
					<GameItem
						imgPath="https://gateway.pinata.cloud/ipfs/QmYQdDyqQqS6XhWRcsqGRejVkFzEHxACiRiPhYAXLQNJGj/Phantom/Oni_Phantom.png"
						onCardClick={() => {
							onWeaponSelect("phantom");
						}}
						cardName={"Phantom"}
					/>
					<GameItem
						imgPath="https://gateway.pinata.cloud/ipfs/QmYQdDyqQqS6XhWRcsqGRejVkFzEHxACiRiPhYAXLQNJGj/Vandal/Elderflame_Vandal.png"
						onCardClick={() => {
							onWeaponSelect("vandal");
						}}
						cardName={"Vandal"}
					/>
					<GameItem
						imgPath="https://gateway.pinata.cloud/ipfs/QmYQdDyqQqS6XhWRcsqGRejVkFzEHxACiRiPhYAXLQNJGj/Guardian/Prime_Guardian.png"
						onCardClick={() => {
							onWeaponSelect("guardian");
						}}
						cardName={"Guardian"}
					/>
				</Grid>
				<Grid item xs={3}>
					<Item>SNIPER</Item>
					<GameItem
						imgPath="https://gateway.pinata.cloud/ipfs/QmYQdDyqQqS6XhWRcsqGRejVkFzEHxACiRiPhYAXLQNJGj/Operator/Ion_Operator.png"
						onCardClick={() => {
							onWeaponSelect("operator");
						}}
						cardName={"Operator"}
					/>
				</Grid>
			</Grid>
		</Box>
	);
};

export default BasicWeaponsGrid;
