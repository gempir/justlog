export default (bttvEmotes) => (dispatch) => {
    dispatch({
        type: 'SET_BTTV_EMOTES',
        bttvEmotes: bttvEmotes
    });
}