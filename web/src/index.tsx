import React from 'react';
import { useContext } from 'react';
import ReactDOM from 'react-dom';
import { ReactQueryCacheProvider } from 'react-query';
import { Page } from './components/Page';
import { StateProvider, store } from './store';

function App() {
	const { state } = useContext(store);

	console.log(process.env);

	return <ReactQueryCacheProvider queryCache={state.queryCache}>
		<Page />
	</ReactQueryCacheProvider>
}

ReactDOM.render(
	<React.StrictMode>
		<StateProvider>
			<App />
		</StateProvider>
	</React.StrictMode>,
	document.getElementById('root')
);
