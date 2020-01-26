export default () => {
    return {
        apiBaseUrl: process.env.apiBaseUrl,
        channels: [],
        logs: {
            messages: [],
        },
        loading: false,
        twitchEmotes: {},
    }
}