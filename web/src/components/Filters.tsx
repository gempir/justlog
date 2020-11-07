import { Button, TextField } from "@material-ui/core";
import React, { FormEvent, useContext } from "react";
import styled from "styled-components";
import { store } from "../store";

const FiltersContainer = styled.form`
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 15px;
    background: var(--bg-bright);
    border-bottom-left-radius: 3px;
    border-bottom-right-radius: 3px;
    width: 600px;
	margin: 0 auto;
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

    return <FiltersContainer onSubmit={handleSubmit} action="none">
        <TextField name="channel" label="channel" variant="filled" autoComplete="off" defaultValue={state.currentChannel} autoFocus={state.currentChannel === null} />
        <TextField name="username" label="username" variant="filled" autoComplete="off" defaultValue={state.currentUsername} autoFocus={state.currentChannel !== null && state.currentUsername === null} />
        <Button variant="contained" color="primary" size="large" type="submit">load</Button>
    </FiltersContainer>
}