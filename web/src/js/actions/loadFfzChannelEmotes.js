import setFfzChannelEmotes from "./setFfzChannelEmotes";

export default function (channelid) {
    return function (dispatch, getState) {
        return new Promise((resolve, reject) => {
            fetch("https://api.frankerfacez.com/v1/room/id/" + channelid).then((response) => {
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
                dispatch(setFfzChannelEmotes(json))

                resolve();
            }).catch(err => {
                reject(err);
            });
        });
    };
}
