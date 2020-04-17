export default (bttvChannelEmotes) => (dispatch) => {
    dispatch({
        type: 'SET_BTTV_CHANNEL_EMOTES',
        bttvChannelEmotes: bttvChannelEmotes
    });
}