import React, { useContext } from "react";
import styled from "styled-components";
import { useLog } from "../hooks/useLog";
import { store } from "../store";
import { TwitchChatLogLine } from "./TwitchChatLogLine";

const ContentLogContainer = styled.ul`
    list-style: none;
    padding: 0;
    margin: 0;
    width: 340px;
`;

export function TwitchChatContentLog({ year, month }: { year: string, month: string }) {
    const { state } = useContext(store);

    const logs = useLog(state.currentChannel ?? "", state.currentUsername ?? "", year, month)

    return <ContentLogContainer>
        {logs.map((log, index) => <TwitchChatLogLine key={log.id ? log.id : index} message={log} />)}
    </ContentLogContainer>
}