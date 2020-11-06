export interface UserLogResponse {
    messages: LogMessage[];
}

export interface LogMessage {
    text: string;
    username: string;
    displayName: string;
    channel: string;
    timestamp: Date;
    type: number;
    raw: string;
}