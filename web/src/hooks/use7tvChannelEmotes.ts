import { useQuery } from "react-query";
import { QueryDefaults } from "../store";
import { StvChannelEmotesResponse } from "../types/7tv";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

export function use7tvChannelEmotes(channelId: string): Array<ThirdPartyEmote> {
	const { isLoading, error, data } = useQuery(["7tv:channel", { channelId: channelId }], () => {
		if (channelId === "") {
			return Promise.resolve(<StvChannelEmotesResponse>{});
		}

		return fetch(`https://7tv.io/v3/users/twitch/${channelId}`).then(res => {
			if (res.ok) {
				return res.json() as Promise<StvChannelEmotesResponse>;
			}

			return Promise.reject(res.statusText);
		});
	}, QueryDefaults);

	if (isLoading) {
		return [];
	}

	if (error) {
		console.error(error);
		return [];
	}

	const emotes = [];

	for (const channelEmote of data?.emote_set?.emotes ?? []) {
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