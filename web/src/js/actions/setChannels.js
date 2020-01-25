export default (channels) => (dispatch) => {
    dispatch({
        type: 'SET_CHANNELS',
        channels: channels
    });
}