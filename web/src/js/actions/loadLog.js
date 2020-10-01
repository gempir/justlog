import setLogs from "./setLogs";
import setLoading from "./setLoading";
import Log from "../model/Log";

export default function (channel, username, year, month) {
    return function (dispatch, getState) {
        return new Promise((resolve, reject) => {
            channel = channel || getState().channel;
            username = username || getState().username;
            channel = channel.toLowerCase();
            username = username.toLowerCase();

            let channelPath = "channel";
            if (channel.startsWith("id:")) {
                channelPath = "id";
            }
            let usernamePath = "user";
            if (username.startsWith("id:")) {
                usernamePath = "userid";
            }

            dispatch(setLoading(true));

            const url = `${getState().apiBaseUrl}/${channelPath}/${channel.replace("id:", "")}/${usernamePath}/${username.replace("id:", "")}/${year}/${month}?reverse`;
            fetch(url, { headers: { "Content-Type": "application/json" } }).then((response) => {
                if (response.status >= 200 && response.status < 300) {
                    return response
                } else {
                    var error = new Error(response.statusText)
                    error.response = response
                    dispatch(setLoading(false));
                    reject(error);
                }
            }).then((response) => {
                return response.json()
            }).then((json) => {
                for (let value of json.messages) {
                    value.timestamp = Date.parse(value.timestamp)
                }

                const logs = {...getState().logs};

                logs[`${year}-${month}`] = new Log(year, month, json.messages, true);

                dispatch(setLogs(logs));
                dispatch(setLoading(false));
                resolve();
            }).catch((error) => {
                dispatch(setLoading(false));
                reject(error);
            });
        });
    };
}