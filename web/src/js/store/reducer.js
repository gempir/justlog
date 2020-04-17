export default (state, action) => {
    switch (action.type) {
        case "SET_CHANNELS":
            return { ...state, channels: action.channels };
        case "SET_LOADING":
            return { ...state, loading: action.loading };
        case "SET_LOGS":
            return { ...state, logs: action.logs };
        case "SET_CURRENT":
            return { ...state, channel: action.channel, username: action.username };
        case "SET_TWITCH_EMOTES":
            return { ...state, twitchEmotes: action.twitchEmotes };
        case "SET_BTTV_CHANNEL_EMOTES":
            return { ...state, bttvChannelEmotes: action.bttvChannelEmotes };
        default:
            return { ...state };
    }
};