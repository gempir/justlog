export type StvGlobalEmotesResponse = StvGlobal
export type StvChannelEmotesResponse = StvChannel

interface StvGlobal {
    emotes: StvEmote[];
}

interface StvChannel {
    emote_set: StvEmoteSet | null;
}

interface StvEmoteSet {
    id: string;
    name: string;
    emotes: StvEmote[];
}


interface StvEmote {
    id: string;
    name: string;
    data: StvEmoteData;
}

interface StvEmoteData {
    id: string;
    name: string;
    listed: boolean;
    animated: boolean;
    host: StvEmoteHost;
}

interface StvEmoteHost {
    url: string;
    files: StvEmoteFile[];
}

interface StvEmoteFile {
    name: string;
    width: number;
    height: number;
    format: string;
}