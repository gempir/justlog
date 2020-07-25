export default (settings) => (dispatch) => {
    localStorage.setItem('settings', JSON.stringify(settings));
    dispatch({
        type: 'SET_SETTINGS',
        settings: settings
    });
}