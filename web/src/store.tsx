import React, { createContext, useState } from "react";
import { QueryCache } from "react-query";
import { useLocalStorage } from "./hooks/useLocalStorage";

export interface Settings {

}

export interface State {
    settings: Settings,
    queryCache: QueryCache,
    apiBaseUrl: string,
    currentChannel: string | null,
    currentUsername: string | null,
}

export type Action = Record<string, unknown>;

const url = new URL(window.location.href);
const defaultContext = {
    state: {
        queryCache: new QueryCache(),
        apiBaseUrl: process.env.REACT_APP_API_BASE_URL,
        currentChannel: url.searchParams.get("channel"),
        currentUsername: url.searchParams.get("username"),
    } as State,
    setState: (state: State) => {},
    setCurrents: (currentChannel: string | null = null, currentUsername: string | null = null) => {},
    setSettings: (newSettings: Settings) => {},
};

const store = createContext(defaultContext);
const { Provider } = store;

const StateProvider = ({ children }: { children: JSX.Element }): JSX.Element => {

    const [settings, setSettingsStorage] = useLocalStorage("settings", defaultContext.state.settings);
    const [state, setState] = useState({ ...defaultContext.state, settings });

    const setSettings = (newSettings: Settings) => {
        setSettingsStorage(newSettings);
        setState({ ...state, settings: newSettings });
    }

    const setCurrents = (currentChannel: string | null = null, currentUsername: string | null = null) => {
        setState({ ...state, currentChannel, currentUsername });
        
        const url = new URL(window.location.href);
        if (currentChannel) {
            url.searchParams.set("channel", currentChannel);
        }
        if (currentUsername) {
            url.searchParams.set("username", currentUsername);
        }

        window.history.replaceState( {} , "justlog", url.toString());
    }

    return <Provider value={{ state, setState, setSettings, setCurrents }}>{children}</Provider>;
};

export { store, StateProvider };
