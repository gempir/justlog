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

            dispatch(setLoading(true));

            let channelPath = "channel";
            if (channel.startsWith("id:")) {
                channelPath = "id";
            }
            let usernamePath = "user";
            if (username.startsWith("id:")) {
                usernamePath = "userid";
            }

            const logs = {}

            fetchAvailableLogs(getState().apiBaseUrl, channel, username).then(data => {
                data.availableLogs.map(log => logs[`${log.year}-${log.month}`] = new Log(log.year, log.month, [], false));

                if (Object.keys(logs).length === 0) {
                    dispatch(setLoading(false));
                    reject(new Error("not found"));
                    return;
                }
                
                year = year || data.availableLogs[0].year;
                month = month || data.availableLogs[0].month;

                const url = `${getState().apiBaseUrl}/${channelPath}/${channel.replace("id:", "")}/${usernamePath}/${username.replace("id:", "")}/${year}/${month}?reverse`;
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

                    logs[`${year}-${month}`] = new Log(year, month, json.messages, true);

                    dispatch(setLogs(logs));
                    dispatch(setLoading(false));
                    resolve();
                }).catch((error) => {
                    dispatch(setLoading(false));
                    reject(error);
                });
            }).catch((error) => {
                dispatch(setLoading(false));
                reject(new Error("not found"));
            });
        });
    };
}

function fetchAvailableLogs(baseUrl, channel, username) {
    let channelQuery = "channel=" + channel;
    if (channel.startsWith("id:")) {
        channelQuery = "channelid" + channel.replace("id:", "")
    }
    let userQuery = "user=" + username;
    if (username.startsWith("id:")) {
        userQuery = "userid=" + username.replace("id:", "")
    }

    return fetch(`${baseUrl}/list?${channelQuery}&${userQuery}`, { headers: { "Content-Type": "application/json" } }).then((response) => {
        if (response.status >= 200 && response.status < 300) {
            return response
        } else {
            var error = new Error(response.statusText)
            error.response = response
            throw error
        }
    }).then((response) => {
        return response.json()
    });
}