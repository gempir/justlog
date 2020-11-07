import { ThirdPartyEmote } from "../types/ThirdPartyEmote";
import { useBttvChannelEmotes } from "./useBttvChannelEmotes";
import { useBttvGlobalEmotes } from "./useBttvGlobalEmotes";
import { useFfzChannelEmotes } from "./useFfzChannelEmotes";
import { useFfzGlobalEmotes } from "./useFfzGlobalEmotes";

export function useThirdPartyEmotes(channelId: string): Array<ThirdPartyEmote> {
	const thirdPartyEmotes: Array<ThirdPartyEmote> = [
		...useBttvChannelEmotes(channelId),
		...useFfzChannelEmotes(channelId),
		...useBttvGlobalEmotes(),
		...useFfzGlobalEmotes(),
	];

	return thirdPartyEmotes;
}