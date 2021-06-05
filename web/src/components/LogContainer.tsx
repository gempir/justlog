import React, { useContext, useEffect } from "react";
import styled from "styled-components";
import { OptOutError } from "../errors/OptOutError";
import { useAvailableLogs } from "../hooks/useAvailableLogs";
import { store } from "../store";
import { Log } from "./Log";
import { OptOutMessage } from "./OptOutMessage";

const LogContainerDiv = styled.div`
    color: white;
    padding: 2rem;
    padding-top: 0;
    width: 100%;
`;

export function LogContainer() {
    const { state } = useContext(store);

    const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
    const ctrlKey = isMac ? "metaKey" : "ctrlKey";

    useEffect(() => {
        const listener = function (e: KeyboardEvent) {
            if (e.key === 'f' && e[ctrlKey] && !state.settings.twitchChatMode.value) {
                e.preventDefault();
                if (state.activeSearchField) {
                    state.activeSearchField.focus();
                }
            }
        };

        window.addEventListener("keydown", listener)

        return () => window.removeEventListener("keydown", listener);
    }, [state.activeSearchField, state.settings.twitchChatMode.value, ctrlKey]);

    const [availableLogs, err] = useAvailableLogs(state.currentChannel, state.currentUsername);
    if (err instanceof OptOutError) {
        return <OptOutMessage />;
    }

    return <LogContainerDiv>
        {availableLogs.map((log, index) => <Log key={`${log.year}:${log.month}`} year={log.year} month={log.month} initialLoad={index === 0} />)}
    </LogContainerDiv>
}