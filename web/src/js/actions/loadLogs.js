import setLogs from "./setLogs";
import setLoading from "./setLoading";

export default function (channel, username, year, month) {
    return function (dispatch, getState) {
        return new Promise((resolve, reject) => {

            dispatch(setLoading(true));
                        
            let options = {
                headers: {
                    "Content-Type": "application/json"
                }
            }
    
            fetch(`/channel/${channel}/user/${username}/${year}/${month}`, options).then((response) => {
                if (response.status >= 200 && response.status < 300) {
                    return response
                } else {
                    var error = new Error(response.statusText)
                    error.response = response
                    throw error
                }
            }).then((response) => {
                return response.json()
            }).then((json) => {
    
                for (let value of json.messages) {
                    value.timestamp = Date.parse(value.timestamp)
                }

                dispatch(setLogs(json));
                dispatch(setLoading(false));
                resolve();
            }).catch((error) => {
                dispatch(setLoading(false));
                reject(error);
            });
        });
    };
}
