import React, { useContext } from "react";
import styled from "styled-components";
import { store } from "../store";
import { LogMessage } from "../types/log";

const LogLineContainer = styled.li`

    .timestamp {
        color: var(--text-dark);
        user-select: none;
    }

    .name {
        margin-left: 5px;
        user-select: none;
    }

    .text {
        margin-left: 5px;
    }
`;

export function LogLine({ message }: { message: LogMessage }) {
    const {state} = useContext(store);

    return <LogLineContainer>
        <span className="timestamp">{message.timestamp.toLocaleString()}</span>
        {state.settings.showName.value && <span className="name">{message.displayName}:</span>}
        <span className="text">{message.text}</span>
    </LogLineContainer>
}