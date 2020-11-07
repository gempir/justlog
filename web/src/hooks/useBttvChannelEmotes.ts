import { useQuery } from "react-query";
import { BttvChannelEmotesResponse } from "../types/Bttv";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

export function useBttvChannelEmotes(channelId: string): Array<ThirdPartyEmote> {
	const { isLoading, error, data } = useQuery(`bttv:channel:${channelId}`, () => {
		if (channelId === "") {
			return Promise.resolve({sharedEmotes: [], channelEmotes: []});
		}

		return fetch(`https://api.betterttv.net/3/cached/users/twitch/${channelId}`).then(res =>
			res.json() as Promise<BttvChannelEmotesResponse>
		);
	});

	if (isLoading) {
		return [];
	}

	if (error) {
		console.error(error);
		return [];
	}

	const emotes = [];

	for (const channelEmote of [...data?.channelEmotes ?? [], ...data?.sharedEmotes ?? []]) {
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