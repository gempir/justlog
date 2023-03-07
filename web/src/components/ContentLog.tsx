import { InputAdornment, TextField } from "@mui/material";
import { Search } from "@mui/icons-material";
import React, { useContext, useState, CSSProperties, useRef, useEffect } from "react";
import styled from "styled-components";
import { useLog } from "../hooks/useLog";
import { store } from "../store";
import { LogLine } from "./LogLine";
import { FixedSizeList as List } from 'react-window';

const ContentLogContainer = styled.ul`
    padding: 0;
    margin: 0;
    position: relative;

    .search {
        position: absolute;
        top: -52px;
        width: 320px;
        left: 0;
    }

    .logLine {
        white-space: nowrap;
    }

    .list {
        scrollbar-color: dark;
    }
`;

export function ContentLog({ year, month }: { year: string, month: string }) {
    const { state, setState } = useContext(store);
    const [searchText, setSearchText] = useState("");

    const logs = useLog(state.currentChannel ?? "", state.currentUsername ?? "", year, month)
        .filter(log => log.text.toLowerCase().includes(searchText.toLowerCase()));

    const Row = ({ index, style }: { index: number, style: CSSProperties }) => (
        <div style={style}><LogLine key={logs[index].id ? logs[index].id : index} message={logs[index]} /></div>
    );

    const search = useRef<HTMLInputElement>(null);

    const handleMouseEnter = () => {
        setState({ ...state, activeSearchField: search.current })
    }

    useEffect(() => {
        setState({ ...state, activeSearchField: search.current })
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return <ContentLogContainer onMouseEnter={handleMouseEnter}>
        <TextField
            className="search"
            label="Search"
            inputRef={search}
            onChange={e => setSearchText(e.target.value)}
            size="small"
            InputProps={{
                startAdornment: (
                    <InputAdornment position="start">
                        <Search />
                    </InputAdornment>
                ),
            }}
        />
        <List
            className="list"
            height={600}
            itemCount={logs.length}
            itemSize={20}
            width={"100%"}
        >
            {Row}
        </List>
    </ContentLogContainer>
}