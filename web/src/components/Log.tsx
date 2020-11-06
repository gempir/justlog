import React from "react";
import styled from "styled-components";

const LogContainer = styled.div`
    background: var(--bg);
    padding: 0.5rem;
    margin-top: 1rem;
`;

export function Log({ year, month, initialLoad = false }: { year: string, month: string, initialLoad?: boolean }) {


    return <LogContainer>
        {year} {month}
    </LogContainer>
}