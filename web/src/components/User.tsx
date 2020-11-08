import React from "react";
import styled from "styled-components";



const UserContainer = styled.div.attrs(props => ({
	style: {
		color: props.color,
	}
}))`
	display: inline-block;
	margin-left: 5px;
	font-weight: bold;
`;

export function User({ displayName, color }: { displayName: string, color: string }): JSX.Element {
	
	const renderColor = color !== "" ? color : "grey";

	return <UserContainer color={renderColor} className="user">
		{displayName}:
	</UserContainer>;
}