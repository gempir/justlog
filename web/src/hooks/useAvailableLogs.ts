import { useContext, useEffect, useState } from "react";
import { store } from "../store";

export type AvailableLogs = Array<{ month: string, year: string }>;

export function useAvailableLogs(channel: string | null, username: string | null, channelIsId = false, usernameIsId = false): AvailableLogs {
    const { state } = useContext(store);
    const [availableLogs, setAvailableLogs] = useState<AvailableLogs>([]);

    useEffect(() => {
        if (channel && username) {
            const queryUrl = new URL(`${state.apiBaseUrl}/list`);
            queryUrl.searchParams.append(`channel${channelIsId ? "id" : ""}`, channel);
            queryUrl.searchParams.append(`user${usernameIsId ? "id" : ""}`, username);

            fetch(queryUrl.toString(), { headers: { "Content-Type": "application/json" } }).then((response) => {
                if (response.ok) {
                    return response;
                }

                throw Error(response.statusText);
            }).then(response => {
                return response.json()
            }).then((data: { availableLogs: AvailableLogs }) => {
                setAvailableLogs(data.availableLogs)
            }).catch(console.error);
        }
    }, [channel, username, channelIsId, usernameIsId, state.apiBaseUrl]);

    return availableLogs;
}