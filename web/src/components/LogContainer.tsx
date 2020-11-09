import React, { useContext } from "react";
import styled from "styled-components";
import { useAvailableLogs } from "../hooks/useAvailableLogs";
import { Txt } from "../icons/Txt";
import { store } from "../store";
import { Log } from "./Log";

const LogContainerDiv = styled.div`
    color: white;
    padding: 2rem;
    padding-top: 0;
    width: 100%;
    position: relative;

    .txt {
        position: absolute;
        top: 5px;
        right: 20px;
        opacity: 0.5;

        &:hover {
            opacity: 1;
        }
    }
`;

export function LogContainer() {
    const { state } = useContext(store);

    const availableLogs = useAvailableLogs(state.currentChannel, state.currentUsername);

    return <LogContainerDiv>
        {availableLogs.map((log, index) => <React.Fragment key={`${log.year}:${log.month}`}>
            <a className="txt" target="__blank" href={`${state.apiBaseUrl}/channel/${state.currentChannel}/user/${state.currentUsername}/${log.year}/${log.month}`} rel="noopener noreferrer"><Txt /></a>
            <Log year={log.year} month={log.month} initialLoad={index === 0} />
        </React.Fragment>)}
    </LogContainerDiv>
}