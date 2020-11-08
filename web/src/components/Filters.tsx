import { Button, TextField } from "@material-ui/core";
import React, { FormEvent, useContext } from "react";
import styled from "styled-components";
import { store } from "../store";
import { Settings } from "./Settings";

const FiltersContainer = styled.form`
    display: inline-flex;
    align-items: center;
    padding: 15px;
    background: var(--bg-bright);
    border-bottom-left-radius: 3px;
    border-bottom-right-radius: 3px;
	margin: 0 auto;

    > * {
        margin-right: 15px !important;    

        &:last-child {
            margin-right: 0 !important;
        }
    }
`;

const FiltersWrapper = styled.div`
    text-align: center;
`;

export function Filters() {
    const { setCurrents, state } = useContext(store);

    const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (e.target instanceof HTMLFormElement) {
            const data = new FormData(e.target);

            setCurrents(data.get("channel") as string | null, data.get("username") as string | null);
        }
    };

    return <FiltersWrapper>
        <FiltersContainer onSubmit={handleSubmit} action="none">
            <TextField name="channel" label="channel" variant="filled" autoComplete="off" defaultValue={state.currentChannel} autoFocus={state.currentChannel === null} />
            <TextField error={state.error} name="username" label="username" variant="filled" autoComplete="off" defaultValue={state.currentUsername} autoFocus={state.currentChannel !== null && state.currentUsername === null} />
            <Button variant="contained" color="primary" size="large" type="submit">load</Button>
            <Settings />
        </FiltersContainer>
    </FiltersWrapper>
}