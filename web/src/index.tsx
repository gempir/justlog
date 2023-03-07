import { StrictMode } from 'react';
import { useContext } from 'react';
import { createRoot } from 'react-dom/client';
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

const container = document.getElementById('root') as Element;
const root = createRoot(container);

root.render(
	<StrictMode>
		<StateProvider>
			<ThemeProvider theme={pageTheme}>
				<App />
			</ThemeProvider>
		</StateProvider>
	</StrictMode>
	);
