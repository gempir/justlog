export default class Log {

    constructor(year, month, messages = [], loaded = false) {
        const monthList = [ "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December" ];
        this.title = `${monthList[month - 1]} ${year} `;
        this.year = year;
        this.month = month;
        this.messages = messages;
        this.loaded = loaded;
    }

    getTitle() {
        return this.title;
    }

    getMessages() {
        return this.messages;
    }

    getLoaded() {
        return this.loaded;
    }
}