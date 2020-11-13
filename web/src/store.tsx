import React, { createContext, useState } from "react";
import { QueryCache } from "react-query";
import { useLocalStorage } from "./hooks/useLocalStorage";

export interface Settings {
    showEmotes: Setting,
    showName: Setting,
    showTimestamp: Setting,
    performanceMode: Setting,
}

export interface Setting {
    displayName: string,
    value: boolean,
}

export interface State {
    settings: Settings,
    queryCache: QueryCache,
    apiBaseUrl: string,
    currentChannel: string | null,
    currentUsername: string | null,
    error: boolean,
    activeSearchField: HTMLInputElement | null,
}

export type Action = Record<string, unknown>;

const url = new URL(window.location.href);
const defaultContext = {
    state: {
        queryCache: new QueryCache(),
        apiBaseUrl: process.env.REACT_APP_API_BASE_URL ?? window.location.protocol + "//" + window.location.host,
        settings: {
            showEmotes: {
                displayName: "Show Emotes",
                value: true,
            },
            showName: {
                displayName: "Show Name",
                value: true,
            },
            showTimestamp: {
                displayName: "Show Timestamp",
                value: true,
            },
            performanceMode: {
                displayName: "Performace Mode",
                value: true,
            }
        },
        currentChannel: url.searchParams.get("channel"),
        currentUsername: url.searchParams.get("username"),
        error: false,
    } as State,
    setState: (state: State) => { },
    setCurrents: (currentChannel: string | null = null, currentUsername: string | null = null) => { },
    setSettings: (newSettings: Settings) => { },
};

const store = createContext(defaultContext);
const { Provider } = store;

const StateProvider = ({ children }: { children: JSX.Element }): JSX.Element => {

    const [settings, setSettingsStorage] = useLocalStorage("justlog:settings", defaultContext.state.settings);
    const [state, setState] = useState({ ...defaultContext.state, settings });

    const setSettings = (newSettings: Settings) => {
        setSettingsStorage(newSettings);
        setState({ ...state, settings: newSettings });
    }

    const setCurrents = (currentChannel: string | null = null, currentUsername: string | null = null) => {
        currentChannel = currentChannel?.toLowerCase().trim() ?? null;
        currentUsername = currentUsername?.toLowerCase().trim() ?? null;

        setState({ ...state, currentChannel, currentUsername, error: false });

        const url = new URL(window.location.href);
        if (currentChannel) {
            url.searchParams.set("channel", currentChannel);
        }
        if (currentUsername) {
            url.searchParams.set("username", currentUsername);
        }

        window.history.replaceState({}, "justlog", url.toString());
    }

    return <Provider value={{ state, setState, setSettings, setCurrents }}>{children}</Provider>;
};

export { store, StateProvider };

export const QueryDefaults = {
	staleTime: 5 * 60  * 1000,
};