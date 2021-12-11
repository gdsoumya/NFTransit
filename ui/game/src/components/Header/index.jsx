import React from "react";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import Container from "@mui/material/Container";

const Header = ({ children }) => {
	return (
		<div>
			<AppBar style={{ background: "#000000" }} position="static">
				<Container maxWidth="xl">
					<Toolbar disableGutters>
						{/* <img
							src="https://www.cisco.com/c/en/us/about/case-studies-customer-success-stories/riot-games/jcr:content/Grid/category_atl/layout-category-atl/blade_1895754487/bladeContents/halves/H-Half-1/emphasis_copy/emphasisContents/image/image.img.png/1616176513071.png"
							alt="riot games"
							style={{ height: 30 }}
						/> */}
						<Typography
							variant="h6"
							noWrap
							component="div"
							sx={{ mr: 2, display: { xs: "none", md: "flex" } }}
							style={{ marginLeft: 30 }}>
							VALEGENDS NFT
						</Typography>
						<Typography
							variant="h6"
							noWrap
							component="div"
							sx={{ mr: 2, display: { xs: "none", md: "flex" } }}
							style={{
								marginLeft: 30,
								fontFamily: "Roboto, sans-serif",
								fontSize: 14,
							}}>
							What is NFT?
						</Typography>
						<Typography
							variant="h6"
							noWrap
							component="div"
							sx={{ mr: 2, display: { xs: "none", md: "flex" } }}
							style={{
								marginLeft: 30,
								fontFamily: "Roboto, sans-serif",
								fontSize: 14,
							}}>
							Transfer NFTs
						</Typography>
						<Typography
							variant="h6"
							noWrap
							component="div"
							sx={{ mr: 2, display: { xs: "none", md: "flex" } }}
							style={{
								marginLeft: 30,
								fontFamily: "Roboto, sans-serif",
								fontSize: 14,
							}}>
							Contact Us
						</Typography>
					</Toolbar>
				</Container>
			</AppBar>
			{children}
		</div>
	);
};

export default Header;
