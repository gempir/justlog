import React, { useContext } from "react";
import { store } from "../store";

export function Page() {
	const { state } = useContext(store);

	return <div>{state.apiBaseUrl}</div>;
}