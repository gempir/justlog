import React, { useContext, useEffect } from "react";
import styled from "styled-components";
import { useAvailableLogs } from "../hooks/useAvailableLogs";
import { store } from "../store";
import { Log } from "./Log";

const LogContainerDiv = styled.div`
    color: white;
    padding: 2rem;
    padding-top: 0;
    width: 100%;
`;

export function LogContainer() {
    const { state } = useContext(store);

    useEffect(() => {
        const listener = function (e: KeyboardEvent) {
            if (e.ctrlKey && e.code === "KeyF") {
                e.preventDefault();
                if (state.activeSearchField) {
                    state.activeSearchField.focus();
                }
            }
        };

        window.addEventListener("keydown", listener)

        return () => window.removeEventListener("keydown", listener);
    }, [state.activeSearchField]);

    const availableLogs = useAvailableLogs(state.currentChannel, state.currentUsername);

    return <LogContainerDiv>
        {availableLogs.map((log, index) => <Log key={`${log.year}:${log.month}`} year={log.year} month={log.month} initialLoad={index === 0} />)}
    </LogContainerDiv>
}