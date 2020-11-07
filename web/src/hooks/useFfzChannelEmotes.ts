import { useQuery } from "react-query";
import { EmoteSet, FfzChannelEmotesResponse } from "../types/Ffz";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

export function useFfzChannelEmotes(channelId: string): Array<ThirdPartyEmote> {
	const { isLoading, error, data } = useQuery(`ffz:channel:${channelId}`, () => {
		if (channelId === "") {
			return Promise.resolve({sets: {}});
		}

		return fetch(`https://api.frankerfacez.com/v1/room/id/${channelId}`).then(res =>
			res.json() as Promise<FfzChannelEmotesResponse>
		);
	});

	if (isLoading || !data?.sets) {
		return [];
	}

	if (error) {
		console.error(error);
		return [];
	}

	const emotes = [];

	for (const set of Object.values(data.sets) as Array<EmoteSet>) {
		for (const channelEmote of set.emoticons) {
			emotes.push({
				id: String(channelEmote.id),
				code: channelEmote.name,
				urls: {
					small: channelEmote.urls["1"],
					medium: channelEmote.urls["2"],
					big: channelEmote.urls["4"],
				}
			});
		}
	}

	return emotes;
}