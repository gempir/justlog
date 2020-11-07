import React, { useContext } from "react";
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

    const availableLogs = useAvailableLogs(state.currentChannel, state.currentUsername);

    return <LogContainerDiv>
        {availableLogs.map(log => <Log key={`${log.year}:${log.month}`} year={log.year} month={log.month} />)}
    </LogContainerDiv>
}