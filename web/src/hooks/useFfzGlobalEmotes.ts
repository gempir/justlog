import { useQuery } from "react-query";
import { EmoteSet, FfzGlobalEmotesResponse } from "../types/Ffz";
import { ThirdPartyEmote } from "../types/ThirdPartyEmote";

export function useFfzGlobalEmotes(): Array<ThirdPartyEmote> {
	const { isLoading, error, data } = useQuery("ffz:global", () => {
		return fetch("https://api.frankerfacez.com/v1/set/global").then(res =>
			res.json() as Promise<FfzGlobalEmotesResponse>
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