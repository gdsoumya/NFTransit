import { makeStyles } from "@mui/styles";
const useStyles = makeStyles((theme) => ({
	cardDiv: {
		width: "250px",
		height: "250px",
		marginLeft: 30,
		margin: theme.spacing(2),
		backgroundImage: `url(${"/assets/game-item.svg"})`,
		backgroundSize: "cover",
		backgroundPosition: "center",
		background: "transparent",
		"&:hover": {
			cursor: "pointer",
			boxShadow: "0px 0px 10px #FFFFFF",
			opacity: 1,
		},
	},
	imageDiv: {
		width: "230px",
		transition: "transform 0.15s ease-in-out",
		marginTop: 70,
	},
	cardContent: {
		display: "flex",
		justifyContent: "center",
		alignItems: "center",
		flexDirection: "column",
	},
}));

export default useStyles;
