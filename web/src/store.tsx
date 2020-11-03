import React, { createContext, useState } from "react";
import { QueryCache } from "react-query";
import { useLocalStorage } from "./hooks/useLocalStorage";

export interface Settings {

}

export interface State {
    settings: Settings,
    queryCache: QueryCache,
    apiBaseUrl: string,
}

export type Action = Record<string, unknown>;

const defaultContext = {
	state: {
        queryCache: new QueryCache(),
        apiBaseUrl: process.env.REACT_APP_API_BASE_URL 
	} as State,
	setState: (state: State) => {
		// do nothing
    },
    setSettings: (newSettings: Settings) => {
		// do nothing
	},
};

const store = createContext(defaultContext);
const { Provider } = store;

const StateProvider = ({ children }: { children: JSX.Element }): JSX.Element => {
	
	const [settings, setSettingsStorage] = useLocalStorage("settings", defaultContext.state.settings);
    const [state, setState] = useState({ ...defaultContext.state, settings});
    
    const setSettings = (newSettings: Settings) => {
        setSettingsStorage(newSettings);
        setState({...state, settings: newSettings});
    }

	return <Provider value={{ state, setState, setSettings}}>{children}</Provider>;
};

export { store, StateProvider };