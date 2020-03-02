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
            month = month || date.getMonth() + 1;

            dispatch(setLoading(true));

            let channelPath = "channel";
            if (channel.toLowerCase().startsWith("id:")) {
                channelPath = "id";
            }
            let usernamePath = "user";
            if (username.toLowerCase().startsWith("id:")) {
                usernamePath = "id";
            }

            channel = channel.replace("id:", "")
            username = username.replace("id:", "")
            const url = `${getState().apiBaseUrl}/${channelPath}/${channel}/${usernamePath}/${username}/${year}/${month}`;

            fetch(url, { headers: { "Content-Type": "application/json" } }).then((response) => {
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
            
                const logs = {...getState().logs};
            
                for (let prevMonth = month; prevMonth >= 1; prevMonth--) {
                    logs[`${year}-${prevMonth}`] = new Log(year, prevMonth, [], false);
                }
            
                for (let prevMonth = 12; prevMonth >= 1; prevMonth--) {
                    logs[`${"2019"}-${prevMonth}`] = new Log("2019", prevMonth, [], false);
                }

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
