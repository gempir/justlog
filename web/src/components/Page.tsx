import React from "react";
import styled from "styled-components";
import { Filters } from "./Filters";
import { LogContainer } from "./LogContainer";

const PageContainer = styled.div`
	
`;

export function Page() {
	return <PageContainer>
		<Filters />
		<LogContainer />
	</PageContainer>;
}