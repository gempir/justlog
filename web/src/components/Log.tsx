import { Button } from "@material-ui/core";
import React, { useContext, useState } from "react";
import styled from "styled-components";
import { useLog } from "../hooks/useLog";
import { Txt } from "../icons/Txt";
import { store } from "../store";
import { LogLine } from "./LogLine";

const LogContainer = styled.div`
    position: relative;
    background: var(--bg-bright);
    border-radius: 3px;
    padding: 0.5rem;
    margin-top: 1rem;

    .txt {
        position: absolute;
        top: 5px;
        right: 5px;
        opacity: 0.5;

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

    return <LogContainer>
        <a className="txt" target="__blank" href={`${state.apiBaseUrl}/channel/${state.currentChannel}/user/${state.currentUsername}/${year}/${month}?reverse`} rel="noopener noreferrer"><Txt /></a>
        <ContentLog year={year} month={month} />
    </LogContainer>
}

const ContentLogContainer = styled.ul`
    list-style: none;
    padding: 0;
    margin: 0;
`;

function ContentLog({ year, month }: { year: string, month: string }) {
    const { state } = useContext(store);

    const logs = useLog(state.currentChannel ?? "", state.currentUsername ?? "", year, month)

    return <ContentLogContainer>
        {logs.map((log, index) => <LogLine key={log.id ? log.id : index} message={log} />)}
    </ContentLogContainer>
}

const LoadableLogContainer = styled.div`

`;

function LoadableLog({ year, month, onLoad }: { year: string, month: string, onLoad: () => void }) {
    return <LoadableLogContainer>
        <Button variant="contained" color="primary" size="large" onClick={onLoad}>load {year}/{month}</Button>
    </LoadableLogContainer>
}