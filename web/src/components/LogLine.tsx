import React from "react";
import styled from "styled-components";
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

    return <LogLineContainer>
        <span className="timestamp">{message.timestamp.toLocaleString()}</span>
        <span className="name">{message.displayName}:</span>
        <span className="text">{message.text}</span>
    </LogLineContainer>
}