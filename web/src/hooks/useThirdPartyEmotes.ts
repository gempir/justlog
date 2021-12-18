import { ThirdPartyEmote } from "../types/ThirdPartyEmote";
import { use7tvChannelEmotes } from "./use7tvChannelEmotes";
import { use7tvGlobalEmotes } from "./use7tvGlobalEmotes";
import { useBttvChannelEmotes } from "./useBttvChannelEmotes";
import { useBttvGlobalEmotes } from "./useBttvGlobalEmotes";
import { useFfzChannelEmotes } from "./useFfzChannelEmotes";
import { useFfzGlobalEmotes } from "./useFfzGlobalEmotes";

export function useThirdPartyEmotes(channelId: string): Array<ThirdPartyEmote> {
	const thirdPartyEmotes: Array<ThirdPartyEmote> = [
		...useBttvChannelEmotes(channelId),
		...useFfzChannelEmotes(channelId),
		...use7tvChannelEmotes(channelId),
		...useBttvGlobalEmotes(),
		...useFfzGlobalEmotes(),
		...use7tvGlobalEmotes(),
	];

	return thirdPartyEmotes;
}