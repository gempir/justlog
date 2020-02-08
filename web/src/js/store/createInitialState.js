export default () => {
    const urlParams = new URLSearchParams(window.location.search);

    const date = new Date();

    return {
        apiBaseUrl: process.env.apiBaseUrl,
        channels: [],
        logs: {
            messages: [],
        },
        loading: false,
        twitchEmotes: {},
        month: date.getMonth() + 1,
        year: date.getFullYear(),
        channel: urlParams.get("channel") || "",
        username: urlParams.get("username") || "",
    }
}