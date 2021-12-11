import { CardMedia, CardContent, Box, Typography } from "@mui/material";
import React from "react";
import Card from "@mui/material/Card";
import useStyles from "./styles";

const GameItem = ({ imgPath, onCardClick, cardName }) => {
	const classes = useStyles();
	return (
		<Card
			onClick={onCardClick}
			className={classes.cardDiv}
			style={{ position: "relative" }}>
			<Box>
				<CardContent className={classes.cardContent}>
					<img src={imgPath} className={classes.imageDiv} />
				</CardContent>
				<Box
					sx={{
						position: "absolute",
						bottom: 0,
						width: "100%",
						bgcolor: "rgba(0, 0, 0, 0.54)",
						padding: "5px",
					}}>
					<Typography
						style={{
							fontSize: 16,
							color: "white",
							marginLeft: 20,
						}}>
						{cardName}
					</Typography>
				</Box>
			</Box>
		</Card>
	);
};

export default GameItem;
