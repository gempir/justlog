import { createContext, useState } from "react";
import { QueryClient } from 'react-query';
import { useLocalStorage } from "./hooks/useLocalStorage";

export interface Settings {
    showEmotes: Setting,
    showName: Setting,
    showTimestamp: Setting,
    twitchChatMode: Setting,
    newOnBottom: Setting,
}

export enum LocalStorageSettings {
    showEmotes,
    showName,
    showTimestamp,
    twitchChatMode,
    newOnBottom,
}

export interface Setting {
    displayName: string,
    value: boolean,
}

export interface State {
    settings: Settings,
    queryClient: QueryClient,
    apiBaseUrl: string,
    currentChannel: string | null,
    currentUsername: string | null,
    error: boolean,
    activeSearchField: HTMLInputElement | null,
    showSwagger: boolean,
    showOptout: boolean,
}

export type Action = Record<string, unknown>;

const url = new URL(window.location.href);
const defaultContext = {
    state: {
        queryClient: new QueryClient(),
        apiBaseUrl: import.meta.env.VITE_API_BASE_URL ?? window.location.protocol + "//" + window.location.host,
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
            twitchChatMode: {
                displayName: "Twitch Chat Mode",
                value: false,
            },
            newOnBottom: {
                displayName: "Newest messages on bottom",
                value: false,
            },
        },
        currentChannel: url.searchParams.get("channel"),
        currentUsername: url.searchParams.get("username"),
        showSwagger: url.searchParams.has("swagger"),
        showOptout: url.searchParams.has("optout"),
        error: false,
    } as State,
    setState: (state: State) => { },
    setCurrents: (currentChannel: string | null = null, currentUsername: string | null = null) => { },
    setSettings: (newSettings: Settings) => { },
    setShowSwagger: (show: boolean) => { },
    setShowOptout: (show: boolean) => { },
};

const store = createContext(defaultContext);
const { Provider } = store;

const StateProvider = ({ children }: { children: JSX.Element }): JSX.Element => {

    const [settings, setSettingsStorage] = useLocalStorage("justlog:settings", defaultContext.state.settings);
    const [state, setState] = useState({ ...defaultContext.state, settings });

    const setShowSwagger = (show: boolean) => {
        const url = new URL(window.location.href);

        if (show) {
            url.searchParams.set("swagger", "")
            url.searchParams.delete("optout");
        } else {
            url.searchParams.delete("swagger");
        }

        window.history.replaceState({}, "justlog", url.toString());

        setState({ ...state, showSwagger: show, showOptout: false })
    }

    const setShowOptout = (show: boolean) => {
        const url = new URL(window.location.href);

        if (show) {
            url.searchParams.set("optout", "");
            url.searchParams.delete("swagger");
        } else {
            url.searchParams.delete("optout");
        }

        window.history.replaceState({}, "justlog", url.toString());

        setState({ ...state, showOptout: show, showSwagger: false })
    }

    const setSettings = (newSettings: Settings) => {
        for (const key of Object.keys(newSettings)) {
            if (typeof (defaultContext.state.settings as unknown as Record<string, Setting>)[key] === "undefined") {
                delete (newSettings as unknown as Record<string, Setting>)[key];
            }
        }

        state.queryClient.removeQueries("log");

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

    return <Provider value={{ state, setState, setSettings, setCurrents, setShowSwagger, setShowOptout }}>{children}</Provider>;
};

export { store, StateProvider };

export const QueryDefaults = {
    staleTime: 5 * 10 * 1000,
};
