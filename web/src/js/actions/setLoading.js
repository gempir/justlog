export default (loading) => (dispatch) => {
    dispatch({
        type: 'SET_LOADING',
        loading: loading
    });
}