import React from "react";
import styled from "styled-components";
import { Filters } from "./Filters";

const PageContainer = styled.div`
	
`;

export function Page() {
	return <PageContainer>
		<Filters />
	</PageContainer>;
}