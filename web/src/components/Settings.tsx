import { IconButton, Menu, MenuItem } from "@mui/material";
import { Check, Clear, Settings as SettingsIcon } from "@mui/icons-material";
import React, { MouseEvent, useContext, useState } from "react";
import styled from "styled-components";
import { Setting, store } from "../store";

const SettingsContainer = styled.div`
	
`;

export function Settings() {
	const { state, setSettings } = useContext(store);
	const [anchorEl, setAnchorEl] = useState<HTMLButtonElement | null>(null);

	const handleClick = (event: MouseEvent<HTMLButtonElement>) => {
		setAnchorEl(event.currentTarget);
	};

	const handleClose = () => {
		setAnchorEl(null);
	};

	const toggleSetting = (name: string, setting: Setting) => {
		const newSetting = { ...setting, value: !setting.value };

		setSettings({ ...state.settings, [name]: newSetting });
	};

	const menuItems = [];

	for (const [name, setting] of Object.entries(state.settings)) {
		menuItems.push(
			<MenuItem key={name} onClick={() => toggleSetting(name, setting)}>
				{setting.value ? <Check /> : <Clear />}&nbsp;&nbsp;{setting.displayName} 
			</MenuItem>
		);
	}

	return (
		<SettingsContainer>
			<IconButton aria-controls="settings" aria-haspopup="true" onClick={handleClick} size="small">
				<SettingsIcon />
			</IconButton>
			<Menu
				id="settings"
				anchorEl={anchorEl}
				keepMounted
				open={Boolean(anchorEl)}
				onClose={handleClose}
			>
				{menuItems}
			</Menu>
		</SettingsContainer>
	);
}