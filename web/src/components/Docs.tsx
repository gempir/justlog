import React, { useContext } from "react";
import styled from "styled-components";
import DescriptionIcon from '@mui/icons-material/Description';
import { IconButton } from "@mui/material";
import SwaggerUI from "swagger-ui-react"
import "swagger-ui-react/swagger-ui.css"
import ReactDOM from "react-dom";
import { store } from "../store";

const DocsWrapper = styled.div`

`;

export function Docs() {
    const { state, setShowSwagger } = useContext(store);

    const handleClick = () => {
        setShowSwagger(!state.showSwagger);
    }

    return <DocsWrapper>
        <IconButton aria-controls="docs" aria-haspopup="true" onClick={handleClick} size="small" color={state.showSwagger ? "primary" : "default"}>
            <DescriptionIcon />
        </IconButton>
        <Swagger show={state.showSwagger} />
    </DocsWrapper>;
}

const SwaggerWrapper = styled.div`
    position: absolute;
    top: 0;
    background: var(--bg);
    left: 0;
    right: 0;
    margin-top: 90px;
    z-index: 999;
    padding-bottom: 90px;

    .swagger-ui {
        background: var(--bg);
        
        .scheme-container {
            background: var(--bg-bright);
        }
    }
`;

interface SwaggerRequest {
    [k: string]: any;
}

function Swagger({ show }: { show: boolean }) {
    const { state } = useContext(store);
    const baseUrl = new URL(state.apiBaseUrl);

    const requestInterceptor = (req: SwaggerRequest): SwaggerRequest => {
        if (req.url.includes("swagger.json")) {
            return req;
        }

        const url = new URL(req.url);

        url.host = baseUrl.host;
        url.protocol = baseUrl.protocol;
        url.port = baseUrl.port;

        req.url = url.toString();

        return req;
    }

    return ReactDOM.createPortal(
        <SwaggerWrapper style={{ display: show ? "block" : "none" }}>
            <SwaggerUI url="/swagger.json" requestInterceptor={requestInterceptor} />
        </SwaggerWrapper>,
        document.body
    );
}