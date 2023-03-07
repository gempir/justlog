import { useQuery } from "react-query";
import { QueryDefaults } from "../store";
import { StvGlobalEmotesResponse } from "../types/7tv";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

export function use7tvGlobalEmotes(): Array<ThirdPartyEmote> {
	const { isLoading, error, data } = useQuery("7tv:global", () => {
		return fetch("https://7tv.io/v3/emote-sets/global").then(res => {
			if (res.ok) {
				return res.json() as Promise<StvGlobalEmotesResponse>;
			}

			return Promise.reject(res.statusText);
		});
	}, QueryDefaults);

	if (isLoading || !data) {
		return [];
	}

	if (error) {
		console.error(error);
		return [];
	}

	const emotes = [];

	for (const channelEmote of data.emotes ?? []) {
		const webpEmotes = channelEmote.data.host.files.filter(i => i.format === 'WEBP');
		const emoteURL = channelEmote.data.host.url;
		emotes.push({
			id: channelEmote.id,
			code: channelEmote.name,
			urls: {
				small: `${emoteURL}/${webpEmotes[0].name}`,
				medium: `${emoteURL}/${webpEmotes[1].name}`,
				big: `${emoteURL}/${webpEmotes[2].name}`,
			}
		});
	}

	return emotes;
}