import { CardMedia, CardContent, Box, Typography } from "@mui/material";
import React from "react";
import Card from "@mui/material/Card";
import useStyles from "./styles";

const GameItem = ({ imgPath, onCardClick }) => {
	console.log(imgPath);
	const classes = useStyles();
	return (
		<Card onClick={onCardClick} className={classes.cardDiv}>
			<CardContent className={classes.cardContent}>
				<img src={imgPath} className={classes.imageDiv} />
			</CardContent>
		</Card>
	);
};

export default GameItem;
