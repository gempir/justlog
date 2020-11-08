import React from "react";
import Linkify from "react-linkify";
import styled from "styled-components";
import { LogMessage } from "../types/log";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

const MessageContainer = styled.div`
	display: inline-flex;
	align-items: center;
`;

const Emote = styled.img`
	/* transform: scale(0.5); */
	max-height: 14px;
	width: auto;
`;

export function Message({ message, thirdPartyEmotes }: { message: LogMessage, thirdPartyEmotes: Array<ThirdPartyEmote> }): JSX.Element {

	const renderMessage = [];

	let replaced;
	let buffer = "";

	for (let x = 0; x <= message.text.length; x++) {
		const c = message.text[x];

		replaced = false;
		for (const emote of message.emotes) {
			if (emote.startIndex === x) {
				replaced = true;
				renderMessage.push(<Emote
					key={x}
					alt={emote.code}
					src={`https://static-cdn.jtvnw.net/emoticons/v1/${emote.id}/1.0`}
				/>);
				x += emote.endIndex - emote.startIndex - 1;
				break;
			}
		}

		if (!replaced) {
			if (c !== " " && x !== message.text.length) {
				buffer += c;
				continue;
			}
			let emoteFound = false;

			for (const emote of thirdPartyEmotes) {
				if (buffer.trim() === emote.code) {
					renderMessage.push(<Emote
						key={x}
						alt={emote.code}
						src={emote.urls.small}
					/>);
					emoteFound = true;
					buffer = "";

					break;
				}
			}

			if (!emoteFound) {
				renderMessage.push(<Linkify key={x} componentDecorator={(decoratedHref, decoratedText, key) => (
					<a target="__blank" href={decoratedHref} key={key}>
						{decoratedText}
					</a>
				)}>{buffer}</Linkify>);
				buffer = "";
			}
			renderMessage.push(c);
		}
	}

	return <MessageContainer className="message">
		{renderMessage}
	</MessageContainer>;
};