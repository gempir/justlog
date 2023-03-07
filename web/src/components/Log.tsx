import { Button } from "@mui/material";
import React, { useContext, useState } from "react";
import styled from "styled-components";
import { Txt } from "../icons/Txt";
import { getUserId, isUserId } from "../services/isUserId";
import { store } from "../store";
import { ContentLog } from "./ContentLog";
import { TwitchChatContentLog } from "./TwitchChatLogContainer";

const LogContainer = styled.div`
    position: relative;
    background: var(--bg-bright);
    border-radius: 3px;
    padding: 0.5rem;
    margin-top: 3rem;

    .txt {
        position: absolute;
        top: 5px;
        right: 15px;
        opacity: 0.9;
        cursor: pointer;
        z-index: 999;

        &:hover {
            opacity: 1;
        }
    }
`;

export function Log({ year, month, initialLoad = false }: { year: string, month: string, initialLoad?: boolean }) {
    const { state } = useContext(store);
    const [load, setLoad] = useState(initialLoad);

    if (!load) {
        return <LogContainer>
            <LoadableLog year={year} month={month} onLoad={() => setLoad(true)} />
        </LogContainer>
    }

    let txtHref = `${state.apiBaseUrl}`
    if (state.currentChannel && isUserId(state.currentChannel)) {
        txtHref += `/channelid/${getUserId(state.currentChannel)}`
    } else {
        txtHref += `/channel/${state.currentChannel}`
    }

    if (state.currentUsername && isUserId(state.currentUsername)) {
        txtHref += `/userid/${getUserId(state.currentUsername)}`
    } else {
        txtHref += `/user/${state.currentUsername}`
    }

    txtHref += `/${year}/${month}?reverse`;

    return <LogContainer>
        <a className="txt" target="__blank" href={txtHref} rel="noopener noreferrer"><Txt /></a>
        {!state.settings.twitchChatMode.value && <ContentLog year={year} month={month} />}
        {state.settings.twitchChatMode.value && <TwitchChatContentLog year={year} month={month} />}
    </LogContainer>
}

const LoadableLogContainer = styled.div`

`;

function LoadableLog({ year, month, onLoad }: { year: string, month: string, onLoad: () => void }) {
    return <LoadableLogContainer>
        <Button variant="contained" color="primary" size="large" onClick={onLoad}>load {year}/{month}</Button>
    </LoadableLogContainer>
}