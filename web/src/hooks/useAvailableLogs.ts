import { useContext } from "react";
import { useQuery } from "react-query";
import { isUserId } from "../services/isUserId";
import { store } from "../store";

export type AvailableLogs = Array<{ month: string, year: string }>;

export function useAvailableLogs(channel: string | null, username: string | null): AvailableLogs {
    const { state } = useContext(store);

    const { data } = useQuery<AvailableLogs>(`${channel}:${username}`, () => {
        if (channel && username) {
            const channelIsId = isUserId(channel);
            const usernameIsId = isUserId(username);
            
            const queryUrl = new URL(`${state.apiBaseUrl}/list`);
            queryUrl.searchParams.append(`channel${channelIsId ? "id" : ""}`, channel);
            queryUrl.searchParams.append(`user${usernameIsId ? "id" : ""}`, username);

            return fetch(queryUrl.toString()).then((response) => {
                if (response.ok) {
                    return response;
                }

                throw Error(response.statusText);
            }).then(response => response.json())
            .then((data: { availableLogs: AvailableLogs }) => data.availableLogs);
        }

        return [];
    });

    return data ?? [];
}