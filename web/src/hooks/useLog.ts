import { useContext } from "react";
import { useQuery } from "react-query";
import { isUserId } from "../services/isUserId";
import { store } from "../store";
import { LogMessage, UserLogResponse } from "../types/log";



export function useLog(channel: string, username: string, year: string, month: string): Array<LogMessage> {
    const { state } = useContext(store);

    const { data } = useQuery<Array<LogMessage>>(`${channel}:${username}:${year}:${month}`, () => {
        if (channel && username) {
            const channelIsId = isUserId(channel);
            const usernameIsId = isUserId(username);

            const queryUrl = new URL(`${state.apiBaseUrl}/channel${channelIsId ? "id" : ""}/${channel}/user${usernameIsId ? "id" : ""}/${username}/2020/11?reverse&json`);

            return fetch(queryUrl.toString()).then((response) => {
                if (response.ok) {
                    return response;
                }

                throw Error(response.statusText);
            }).then(response => response.json()).then((data: UserLogResponse) => {
                const messages: Array<LogMessage> = [];

                for (const msg of data.messages) {
                    messages.push({ ...msg, timestamp: new Date(msg.timestamp) })
                }

                return messages;
            });
        }

        return [];
    });

    return data ?? [];
}