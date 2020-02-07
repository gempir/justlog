export default () => {
    const urlParams = new URLSearchParams(window.location.search);

    return {
        apiBaseUrl: process.env.apiBaseUrl,
        channels: [],
        logs: {
            messages: [],
        },
        loading: false,
        twitchEmotes: {},
        currentChannel: urlParams.get("channel") || "",
        currentUsername: urlParams.get("username") || "",
    }
}