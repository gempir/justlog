import React, { useContext } from "react";
import styled from "styled-components";
import { useAvailableLogs } from "../hooks/useAvailableLogs";
import { store } from "../store";
import { Log } from "./Log";

const LogContainerDiv = styled.div`
    margin: 2rem;
    color: white;
`;

export function LogContainer() {
    const { state } = useContext(store);

    const availableLogs = useAvailableLogs(state.currentChannel, state.currentUsername);

    return <LogContainerDiv>
        {availableLogs.map(log => <Log key={`${log.year}:${log.month}`} year={log.year} month={log.month} />)}
    </LogContainerDiv>
}