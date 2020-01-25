export default (twitchEmotes) => (dispatch) => {
    dispatch({
        type: 'SET_TWITCH_EMOTES',
        twitchEmotes: twitchEmotes
    });
}