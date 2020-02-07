export default () => {
    // const regex = /\/(\w+)\/(\w+)\??/;
    // console.log(window.location.pathname.match(regex));

    return {
        apiBaseUrl: process.env.apiBaseUrl,
        channels: [],
        logs: {
            messages: [],
        },
        loading: false,
        twitchEmotes: {},
        currentChannel: "",
        currentUsername: "",
    }
}