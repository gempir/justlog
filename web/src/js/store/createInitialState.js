export default () => {
    const urlParams = new URLSearchParams(window.location.search);

    return {
        apiBaseUrl: process.env.apiBaseUrl,
        channels: [],
        logs: {},
        loading: false,
        twitchEmotes: {},
        channel: urlParams.get("channel") || "",
        username: urlParams.get("username") || "",
    }
}