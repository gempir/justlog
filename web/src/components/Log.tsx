import { Button } from "@material-ui/core";
import React, { useContext, useState } from "react";
import styled from "styled-components";
import { useLog } from "../hooks/useLog";
import { store } from "../store";

const LogContainer = styled.div`
    background: var(--bg);
    padding: 0.5rem;
    margin-top: 1rem;
`;

export function Log({ year, month, initialLoad = false }: { year: string, month: string, initialLoad?: boolean }) {
    const [load, setLoad] = useState(initialLoad);

    if (!load) {
        return <LoadableLog year={year} month={month} onLoad={() => setLoad(true)} />
    }

    return <LogContainer>
        <ContentLog year={year} month={month} />
    </LogContainer>
}

const ContentLogContainer = styled.div`
`;

function ContentLog({ year, month }: { year: string, month: string }) {
    const {state} = useContext(store);

    const logs = useLog(state.currentChannel ?? "", state.currentUsername ?? "", year, month)

    return <ContentLogContainer>
        {JSON.stringify(logs)}
    </ContentLogContainer>
}

const LoadableLogContainer = styled.div`

`;

function LoadableLog({ year, month, onLoad }: { year: string, month: string, onLoad: () => void }) {
    return <LoadableLogContainer>
        <Button variant="contained" color="primary" size="large" onClick={onLoad}>load</Button>
    </LoadableLogContainer>
}