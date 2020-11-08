export interface BttvChannelEmotesResponse {
    // id: string;
    // bots: string[];
    channelEmotes: Emote[];
    sharedEmotes: Emote[];
}

export type BttvGlobalEmotesResponse = Emote[];

export interface Emote {
    id: string;
    code: string;
    // imageType: ImageType;
    // userId?: string;
    // user?: User;
}

export enum ImageType {
    GIF = "gif",
    PNG = "png",
}

export interface User {
    id: string;
    name: string;
    displayName: string;
    providerId: string;
}