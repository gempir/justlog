import setTwitchEmotes from "./setTwitchEmotes";

export default function () {
    return function (dispatch, getState) {
        return new Promise((resolve, reject) => {
            
            fetch("https://twitchemotes.com/api_cache/v3/global.json").then((response) => {
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
                dispatch(setTwitchEmotes(json));
                
                resolve();
            }).catch(() => {
                reject();
            });
            
            
        });
    };
}
