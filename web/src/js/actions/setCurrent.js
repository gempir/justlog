export default (channel, username) => (dispatch) => {
    dispatch({
        type: 'SET_CURRENT',
        currentChannel: channel,
        currentUsername: username
    });
}