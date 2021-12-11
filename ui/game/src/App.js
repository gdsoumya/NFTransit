import { createTheme, ThemeProvider } from "@mui/material";
import Header from "./components/Header";
import HomePage from "./pages/HomePage";

function App() {
	const theme = createTheme({
		typography: {
			fontFamily: "VALORANT",
		},
	});

	return (
		<ThemeProvider theme={theme}>
			{/* Weapon info div */}
			<Header>
				<HomePage />
			</Header>
		</ThemeProvider>
	);
}

export default App;
