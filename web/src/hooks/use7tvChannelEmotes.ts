import { useQuery } from "react-query";
import { QueryDefaults } from "../store";
import { StvChannelEmotesResponse } from "../types/7tv";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

export function use7tvChannelEmotes(channelId: string): Array<ThirdPartyEmote> {
	const { isLoading, error, data } = useQuery(["7tv:channel", { channelId: channelId }], () => {
		if (channelId === "") {
			return Promise.resolve([]);
		}

		return fetch(`https://api.7tv.app/v2/users/${channelId}/emotes`).then(res =>
			res.json() as Promise<StvChannelEmotesResponse>
		);
	}, QueryDefaults);

	if (isLoading) {
		return [];
	}

	if (error) {
		console.error(error);
		return [];
	}

	const emotes = [];

	for (const channelEmote of data ?? []) {
		emotes.push({
			id: channelEmote.id,
			code: channelEmote.name,
			urls: {
				small: `https://cdn.7tv.app/emote/${channelEmote.id}/1x`,
				medium: `https://cdn.7tv.app/emote/${channelEmote.id}/2x`,
				big: `https://cdn.7tv.app/emote/${channelEmote.id}/3x`,
			}
		});
	}

	return emotes;
}