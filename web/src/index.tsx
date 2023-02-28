import React from 'react';
import { useContext } from 'react';
import ReactDOM from 'react-dom';
import { QueryClientProvider } from 'react-query';
import { Page } from './components/Page';
import { StateProvider, store } from './store';
import { createTheme } from '@mui/material';
import { ThemeProvider } from '@mui/material/styles';

const pageTheme = createTheme({
	palette: {
		mode: 'dark'
	},
});

function App() {
	const { state } = useContext(store);

	return <QueryClientProvider client={state.queryClient}>
		<Page />
	</QueryClientProvider>
}

ReactDOM.render(
	<React.StrictMode>
		<StateProvider>
			<ThemeProvider theme={pageTheme}>
				<App />
			</ThemeProvider>
		</StateProvider>
	</React.StrictMode>,
	document.getElementById('root')
);
