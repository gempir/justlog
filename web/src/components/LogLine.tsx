import dayjs from "dayjs";
import React, { useContext } from "react";
import styled from "styled-components";
import { useThirdPartyEmotes } from "../hooks/useThirdPartyEmotes";
import { store } from "../store";
import { LogMessage } from "../types/log";
import { Message } from "./Message";
import { User } from "./User";
import utc from "dayjs/plugin/utc";
import timezone from "dayjs/plugin/timezone";

dayjs.extend(utc)
dayjs.extend(timezone)
dayjs.tz.guess()

const LogLineContainer = styled.li`
	display: flex;
	align-items: flex-start;
    margin-bottom: 1px;

    .timestamp {
        color: var(--text-dark);
        user-select: none;
        font-family: monospace;
        white-space: nowrap;
        line-height: 1.1rem;
    }

    .user {
        margin-left: 5px;
        user-select: none;
        font-weight: bold;
        line-height: 1.1rem;
    }

    .message {
        margin-left: 5px;
        line-height: 1.1rem;
    }
`;

export function LogLine({ message }: { message: LogMessage }) {
    const { state } = useContext(store);

    if (state.settings.showEmotes.value) {
        return <LogLineWithEmotes message={message} />;
    }

    return <LogLineContainer className="logLine">
        {state.settings.showTimestamp.value &&<span className="timestamp">{dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}</span>}
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={[]} />
    </LogLineContainer>
}

export function LogLineWithEmotes({ message }: { message: LogMessage }) {
    const thirdPartyEmotes = useThirdPartyEmotes(message.tags["room-id"])
    const { state } = useContext(store);

    return <LogLineContainer className="logLine">
        {state.settings.showTimestamp.value &&<span className="timestamp">{dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}</span>}
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={thirdPartyEmotes} />
    </LogLineContainer>
}