import React, { useContext } from "react";
import Linkify from "react-linkify";
import styled from "styled-components";
import { store } from "../store";
import { LogMessage } from "../types/log";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";
import runes from "runes";

const MessageContainer = styled.div`

	a {
		margin: 0 2px;
		color: var(--theme2);
		text-decoration: none;

		&:hover, &:active, &:focus {
			color: var(--theme2-bright);
		}
	}
`;

const Emote = styled.img`
	max-height: 18px;
	margin: 0 2px;
	margin-bottom: -2px;
	width: auto;
`;

export function Message({ message, thirdPartyEmotes }: { message: LogMessage, thirdPartyEmotes: Array<ThirdPartyEmote> }): JSX.Element {
	const { state } = useContext(store);
	const renderMessage = [];

	let replaced;
	let buffer = "";

	let messageText = message.text;
	let renderMessagePrefix = "";
	if (message.tags['system-msg']) {
		messageText = messageText.replace(message.tags['system-msg'] + " ", "");

		renderMessagePrefix = `${message.tags['system-msg']} `;
	}

	const messageTextEmoji = runes(messageText);

	for (let x = 0; x <= messageTextEmoji.length; x++) {
		const c = messageTextEmoji[x];

		replaced = false;

		if (state.settings.showEmotes.value) {
			for (const emote of message.emotes) {
				if (emote.startIndex === x) {
					replaced = true;
					renderMessage.push(<Emote
						className="emote"
						key={x}
						alt={emote.code}
						title={emote.code}
						src={`https://static-cdn.jtvnw.net/emoticons/v2/${emote.id}/default/dark/1.0`}
					/>);
					x += emote.endIndex - emote.startIndex - 1;
					break;
				}
			}
		}

		if (!replaced) {
			if (c !== " " && x !== messageTextEmoji.length) {
				buffer += c;
				continue;
			}
			let emoteFound = false;

			for (const emote of thirdPartyEmotes) {
				if (buffer.trim() === emote.code) {
					renderMessage.push(<Emote
						className="emote"
						key={x}
						alt={emote.code}
						title={emote.code}
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

			if (c === " " && !state.settings.twitchChatMode.value) {
				renderMessage.push(<>&nbsp;</>);
			} else {
				renderMessage.push(c);
			}
		}
	}

	return <MessageContainer className="message">
		{renderMessagePrefix}{renderMessage}
	</MessageContainer>;
};
