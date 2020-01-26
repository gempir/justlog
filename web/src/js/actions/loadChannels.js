import setChannels from "./setChannels";

export default function () {
    return function (dispatch, getState) {
        return new Promise((resolve, reject) => {
            fetch(`${getState().apiBaseUrl}/channels`).then((response) => {
                if (response.status >= 200 && response.status < 300) {
                    return response
                } else {
                    var error = new Error(response.statusText)
                    error.response = response
                    throw error
                }
            }).then((response) => {
                return response.json();
            }).then((json) => {
                dispatch(setChannels(json.channels));

                resolve();
            }).catch(() => {
                reject();
            });
        });
    };
}
