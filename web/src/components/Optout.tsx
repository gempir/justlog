import { IconButton, Button } from "@mui/material";
import { useContext, useState } from "react";
import styled from "styled-components";
import { store } from "../store";
import CancelIcon from '@mui/icons-material/Cancel';

const OptoutWrapper = styled.div`

`;

export function Optout() {
    const { state, setShowOptout } = useContext(store);

    const handleClick = () => {
        setShowOptout(!state.showOptout);
    }

    return <OptoutWrapper>
        <IconButton aria-controls="docs" aria-haspopup="true" onClick={handleClick} size="small" color={state.showOptout ? "primary" : "default"}>
            <CancelIcon />
        </IconButton>
    </OptoutWrapper>;
}

const OptoutPanelWrapper = styled.div`
    background: var(--bg-bright);
    color: var(--text);
    margin: 3rem;
    font-size: 1.5rem;
    padding: 2rem;

    code {
        background: var(--bg);
        padding: 1rem;
        border-radius: 3px;
    }

    .generator {
        margin-top: 2rem;
        display: flex;
        gap: 1rem;
        align-items: center;

        input {
            background: var(--bg);
            border: none;
            color: white;
            padding: 0.6rem;
            font-size: 1.5rem;
            text-align: center;
            border-radius: 3px;
        }
    }

    .small {
        font-size: 0.8rem;
        font-family: monospace;
    }
`;

export function OptoutPanel() {
    const { state } = useContext(store);
    const [code, setCode] = useState("");

    const generateCode = () => {
        fetch(state.apiBaseUrl + "/optout", { method: "POST" }).then(res => res.json()).then(setCode).catch(console.error);
    };

    return <OptoutPanelWrapper>
        <p>
            You can opt out from being logged. This will also disable access to your previously logged data.<br />
            This applies to all chats of that justlog instance.<br />
            Opting out is permanent, there is no reverse action. So think twice if you want to opt out.
        </p>
        <p>
            If you still want to optout generate a token here and paste the command into a logged chat.<br />
            You will receive a confirmation message from the bot "@username, opted you out".
        </p>
        <br />
        <div><code>!justlog optout {"<code>"}</code></div>
        <div className="generator">
            <input readOnly type="text" value={code} /><Button variant="contained" onClick={generateCode} color="primary" size="large">Generate Code</Button>
        </div>
        {code && <p className="small">
            This code is valid for 60 seconds
        </p>}
    </OptoutPanelWrapper>;
}