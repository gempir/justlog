import React, { Component } from 'react';
import { createStore, applyMiddleware } from "redux";
import thunk from "redux-thunk";
import { Provider } from "react-redux";
import reducer from "./store/reducer";
import createInitialState from "./store/createInitialState";
import LogSearch from './components/LogSearch';


export default class App extends Component {

	constructor(props) {
		super(props);

		this.store = createStore(reducer, createInitialState(), applyMiddleware(thunk));
	}

	render() {
		return (
			<Provider store={this.store}>
				<LogSearch />
			</Provider>
		);
	}
}