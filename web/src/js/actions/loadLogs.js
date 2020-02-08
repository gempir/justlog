import setLogs from "./setLogs";
import setLoading from "./setLoading";
import Log from "../model/Log";

export default function (channel, username, year, month) {
    return function (dispatch, getState) {
        return new Promise((resolve, reject) => {
            channel = channel || getState().channel;
            username = username || getState().username;
            const date = new Date();
            year = year || date.getFullYear();
            month = month ||  date.getMonth() + 1;

            dispatch(setLoading(true));
                        
            let options = {
                headers: {
                    "Content-Type": "application/json"
                }
            }
    
            fetch(`${getState().apiBaseUrl}/channel/${channel}/user/${username}/${year}/${month}`, options).then((response) => {
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

                const newLogs = {...getState().logs};
                newLogs[`${year}-${month}`] = new Log(year, month, json.messages, true);

                dispatch(setLogs(newLogs));
                dispatch(setLoading(false));
                resolve();
            }).catch((error) => {
                dispatch(setLoading(false));
                reject(error);
            });
        });
    };
}
