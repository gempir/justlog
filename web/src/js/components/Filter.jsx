import React, { Component } from 'react';

export default class Filter extends Component {

    constructor(props) {
        super(props);

        const date = new Date();

        this.state = {
            channel: "",
            username: "",
            year: date.getFullYear(),
            month: date.getMonth() + 1,
        }
    }

    render() {
        return (
            <form className="filter" autoComplete="off" onSubmit={this.onSubmit}>
                <input
                    type="text"
                    placeholder="forsen"
                    autoComplete={"off"}
                    onChange={this.onChannelChange}
                />
                <input
                    type="text"
                    autoComplete={"off"}
                    onChange={this.onUsernameChange}
                    placeholder="gempir"
                />
                <div className="date">
                    <select onChange={this.onYearChange} value={this.state.year}>
                        <option value="2020">2020</option>
                        <option value="2019">2019</option>
                        <option value="2018">2018</option>
                    </select>
                    <select onChange={this.onMonthChange} value={this.state.month}>
                        <option value="1">1</option>
                        <option value="2">2</option>
                        <option value="3">3</option>
                        <option value="4">4</option>
                        <option value="5">5</option>
                        <option value="6">6</option>
                        <option value="7">7</option>
                        <option value="8">8</option>
                        <option value="9">9</option>
                        <option value="10">10</option>
                        <option value="11">11</option>
                        <option value="12">12</option>
                    </select>
                </div>
                <button type="submit" className="show-logs">Show logs</button>
            </form>
        )
    }

    onChannelChange = (e) => {
        this.setState({ channel: e.target.value });
    }

    onUsernameChange = (e) => {
        this.setState({ username: e.target.value });
    }

    onYearChange = (e) => {
        this.setState({ year: e.target.value });
    }

    onMonthChange = (e) => {
        this.setState({ month: e.target.value });
    }

    onSubmit = (e) => {
        e.preventDefault();
        this.props.searchLogs(this.state.channel, this.state.username, this.state.year, this.state.month);
    }
}