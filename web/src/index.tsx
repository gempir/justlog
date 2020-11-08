import React from 'react';
import { useContext } from 'react';
import ReactDOM from 'react-dom';
import { ReactQueryCacheProvider } from 'react-query';
import { Page } from './components/Page';
import { StateProvider, store } from './store';
import { unstable_createMuiStrictModeTheme as createMuiTheme } from '@material-ui/core';
import { ThemeProvider } from '@material-ui/core/styles';

const pageTheme = createMuiTheme({
	palette: {
		type: 'dark'
	},
});

function App() {
	const { state } = useContext(store);

	return <ReactQueryCacheProvider queryCache={state.queryCache}>
		<Page />
	</ReactQueryCacheProvider>
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
