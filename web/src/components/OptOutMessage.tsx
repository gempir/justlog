import styled from "styled-components";
import React from "react";

const OptOutContainer = styled.div`
    display: block;
    font-weight: bold;
    color: var(--danger);
    font-size: 2rem;
    text-align: center;
    padding: 2rem;
`;

export function OptOutMessage() {
    return <OptOutContainer>User or channel has opted out</OptOutContainer>
}