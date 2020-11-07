import { useQuery } from "react-query";
import { BttvGlobalEmotesResponse } from "../types/Bttv";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

export function useBttvGlobalEmotes(): Array<ThirdPartyEmote> {
	const { isLoading, error, data } = useQuery("bttv:global", () => {
		return fetch("https://api.betterttv.net/3/cached/emotes/global").then(res =>
			res.json() as Promise<BttvGlobalEmotesResponse>
		);
	});

	if (isLoading || !data) {
		return [];
	}

	if (error) {
		console.error(error);
		return [];
	}

	const emotes = [];

	for (const channelEmote of data) {
		emotes.push({
			id: channelEmote.id,
			code: channelEmote.code,
			urls: {
				small: `https://cdn.betterttv.net/emote/${channelEmote.id}/1x`,
				medium: `https://cdn.betterttv.net/emote/${channelEmote.id}/2x`,
				big: `https://cdn.betterttv.net/emote/${channelEmote.id}/3x`,
			}
		});
	}

	return emotes;
}