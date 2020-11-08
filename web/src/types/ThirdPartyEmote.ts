export interface ThirdPartyEmote {
    code: string,
    id: string,
    urls: EmoteUrls,
}

export interface EmoteUrls {
    small: string,
    medium: string,
    big: string,
}