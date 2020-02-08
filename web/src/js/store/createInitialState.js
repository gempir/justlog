import Log from "../model/Log";

export default () => {
    const urlParams = new URLSearchParams(window.location.search);

    const date = new Date();
    const month =  date.getMonth() + 1;
    const year = date.getFullYear();

    const logs = {};

    for (let prevMonth = month; prevMonth >= 1; prevMonth--) {
        logs[`${year}-${prevMonth}`] = new Log(year, prevMonth, [], false);
    }

    for (let prevMonth = 12; prevMonth >= 1; prevMonth--) {
        logs[`${"2019"}-${prevMonth}`] = new Log("2019", prevMonth, [], false);
    }

    return {
        apiBaseUrl: process.env.apiBaseUrl,
        channels: [],
        logs: logs,
        loading: false,
        twitchEmotes: {},
        channel: urlParams.get("channel") || "",
        username: urlParams.get("username") || "",
    }
}