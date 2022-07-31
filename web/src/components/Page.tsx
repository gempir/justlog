import React, { useContext } from "react";
import styled from "styled-components";
import { store } from "../store";
import { Filters } from "./Filters";
import { LogContainer } from "./LogContainer";
import { OptoutPanel } from "./Optout";

const PageContainer = styled.div`
	
`;

export function Page() {
	const {state} = useContext(store);

	return <PageContainer>
		<Filters />
		{state.showOptout && <OptoutPanel />}
		<LogContainer />
	</PageContainer>;
}