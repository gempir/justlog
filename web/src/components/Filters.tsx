import { Button, TextField } from "@material-ui/core";
import React, { FormEvent } from "react";
import styled from "styled-components";

const FiltersContainer = styled.form`
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 15px;
    background: var(--bg);
    width: 600px;
	margin: 0 auto;
`;

export function Filters() {
    const handleSubmit = (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        if (e.target instanceof HTMLFormElement) {
            const data = new FormData(e.target);

            console.log(data.get("channel"));
            console.log(data.get("username"));
        }
    };

    return <FiltersContainer onSubmit={handleSubmit}>
        <TextField name="channel" label="channel" variant="filled" autoComplete="off" />
        <TextField name="username" label="username" variant="filled" autoComplete="off" />
        <Button variant="contained" color="primary" size="large" type="submit">load</Button>
    </FiltersContainer>
}