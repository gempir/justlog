export default () => {
    const urlParams = new URLSearchParams(window.location.search);

    return {
        apiBaseUrl: process.env.apiBaseUrl,
        channels: [],
        bttvChannelEmotes: null,
        ffzChannelEmotes: null,
        bttvEmotes: null,
        logs: {},
        loading: false,
        twitchEmotes: {},
        channel: urlParams.get("channel") || "",
        username: urlParams.get("username") || "",
        settings: JSON.parse(localStorage.getItem('settings')) || { showEmotes: true, showDisplayName: false },
    }
}