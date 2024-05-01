export interface UserLogResponse {
    messages: RawLogMessage[],
}

export interface LogMessage extends Omit<RawLogMessage, "timestamp"> {
    timestamp: Date,
    emotes: Array<Emote>,
}

export interface RawLogMessage {
    text: string,
    systemText: string,
    username: string,
    displayName: string,
    channel: string,
    timestamp: string,
    type: number,
    raw: string,
    id: string,
    tags: Record<string, string>,
}

export interface Emote {
    startIndex: number,
    endIndex: number,
    code: string,
    id: string,
}