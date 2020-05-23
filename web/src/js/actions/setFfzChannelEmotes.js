export default (ffzChannelEmotes) => (dispatch) => {
    dispatch({
        type: 'SET_FFZ_CHANNEL_EMOTES',
        ffzChannelEmotes: ffzChannelEmotes
    });
}