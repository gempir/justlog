import React, { useContext } from "react";
import styled from "styled-components";
import { store } from "../store";

const LogContainerDiv = styled.div`
    margin: 2rem;
    padding: 0.5rem;
    background: var(--bg);
    color: white;
`;

export function LogContainer() {
    const {state} = useContext(store);


    return <LogContainerDiv>
        {state.currentChannel}
        {state.currentUsername}
    </LogContainerDiv>
}