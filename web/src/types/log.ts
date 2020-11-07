export interface UserLogResponse {
    messages: RawLogMessage[],
}

export interface LogMessage extends Omit<RawLogMessage, "timestamp"> {
    timestamp: Date,
}

export interface RawLogMessage {
    text: string,
    username: string,
    displayName: string,
    channel: string,
    timestamp: string,
    type: number,
    raw: string,
    id: string,
    tags: Record<string, string>,
}