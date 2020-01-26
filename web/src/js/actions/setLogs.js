export default (logs) => (dispatch) => {
    dispatch({
        type: 'SET_LOGS',
        logs: logs
    });
}