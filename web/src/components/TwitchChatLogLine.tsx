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

const TwitchChatLogLineContainer = styled.li`
	align-items: flex-start;
    margin-bottom: 1px;
    padding: 5px 20px;

    .timestamp {
        color: var(--text-dark);
        user-select: none;
        font-family: monospace;
        white-space: nowrap;
        margin-right: 5px;
        line-height: 1.1rem;
    }

    .user {
        display: inline-block;
        margin-right: 5px;
        user-select: none;
        font-weight: bold;
        line-height: 1.1rem;
    }

    .message {
        display: inline;
        line-height: 20px;

        a {
            word-wrap: break-word;
        }

        .emote {
            max-height: 28px;
            margin: 0 2px;
            margin-bottom: -10px;
            width: auto;
        }
    }
`;

export function TwitchChatLogLine({ message }: { message: LogMessage }) {
    const { state } = useContext(store);

    if (state.settings.showEmotes.value) {
        return <LogLineWithEmotes message={message} />;
    }

    return <TwitchChatLogLineContainer className="logLine">
        {state.settings.showTimestamp.value && <span className="timestamp">{dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}</span>}
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={[]} />
    </TwitchChatLogLineContainer>
}

function LogLineWithEmotes({ message }: { message: LogMessage }) {
    const thirdPartyEmotes = useThirdPartyEmotes(message.tags["room-id"])
    const { state } = useContext(store);

    return <TwitchChatLogLineContainer className="logLine">
        {state.settings.showTimestamp.value && <span className="timestamp">{dayjs(message.timestamp).format("YYYY-MM-DD HH:mm:ss")}</span>}
        {state.settings.showName.value && <User displayName={message.displayName} color={message.tags["color"]} />}
        <Message message={message} thirdPartyEmotes={thirdPartyEmotes} />
    </TwitchChatLogLineContainer>
}