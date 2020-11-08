import React, { useContext } from "react";
import styled from "styled-components";
import { useThirdPartyEmotes } from "../hooks/useThirdPartyEmotes";
import { store } from "../store";
import { LogMessage } from "../types/log";
import { Message } from "./Message";
import { User } from "./User";

const LogLineContainer = styled.li`
	display: flex;
	align-items: center;

    .timestamp {
        color: var(--text-dark);
        user-select: none;
    }

    .user {
        margin-left: 5px;
        user-select: none;
    }

    .message {
        margin-left: 5px;
    }
`;

export function LogLine({ message }: { message: LogMessage }) {
    const { state } = useContext(store);

    if (state.settings.showEmotes.value) {
        return <LogLineWithEmotes message={message} />;
    }

    return <LogLineContainer>
        <span className="timestamp">{message.timestamp.toISOString()}</span>
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={[]} />
    </LogLineContainer>
}

export function LogLineWithEmotes({ message }: { message: LogMessage }) {
    const thirdPartyEmotes = useThirdPartyEmotes(message.tags["room-id"])
    const { state } = useContext(store);

    return <LogLineContainer>
        <span className="timestamp">{message.timestamp.toISOString()}</span>
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={thirdPartyEmotes} />
    </LogLineContainer>
}