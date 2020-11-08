import { useContext } from "react";
import { useQuery } from "react-query";
import { store } from "../store";

export interface Channel {
    userID: string,
    name: string
}

export function useChannels(): Array<Channel> {
    const { state } = useContext(store);

    const { data } = useQuery<Array<Channel>>(`channels`, () => {

        const queryUrl = new URL(`${state.apiBaseUrl}/channels`);

        return fetch(queryUrl.toString()).then((response) => {
            if (response.ok) {
                return response;
            }

            throw Error(response.statusText);
        }).then(response => response.json())
            .then((data: { channels: Array<Channel> }) => data.channels);
    }, { refetchOnWindowFocus: false, refetchOnReconnect: false });

    return data ?? [];
}