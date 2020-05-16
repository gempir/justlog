import setBttvChannelEmotes from "./setBttvChannelEmotes";

export default function (channelid) {
    return function (dispatch, getState) {
        return new Promise((resolve, reject) => {
            fetch("https://api.betterttv.net/3/cached/users/twitch/" + channelid).then((response) => {
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
                dispatch(setBttvChannelEmotes(json))

                resolve();
            }).catch(err => {
                reject(err);
            });
        });
    };
}
