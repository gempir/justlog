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
	align-items: center;

    .timestamp {
        color: var(--text-dark);
        user-select: none;
        font-family: monospace;
    }

    .user {
        margin-left: 5px;
        user-select: none;
        font-weight: bold;
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
        <span className="timestamp">{dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}</span>
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={[]} />
    </LogLineContainer>
}

export function LogLineWithEmotes({ message }: { message: LogMessage }) {
    const thirdPartyEmotes = useThirdPartyEmotes(message.tags["room-id"])
    const { state } = useContext(store);

    return <LogLineContainer>
        <span className="timestamp">{dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}</span>
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={thirdPartyEmotes} />
    </LogLineContainer>
}