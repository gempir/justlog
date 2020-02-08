import React, { Component } from "react";
import { connect } from "react-redux";
import setCurrent from "./../actions/setCurrent";
import AutocompleteInput from "./AutocompleteInput";

class Filter extends Component {

    constructor(props) {
        super(props);

        const date = new Date();

        this.state = {
            year: date.getFullYear(),
            month: date.getMonth() + 1,
        }
    }

    componentDidMount() {
        if (this.props.currentChannel && this.props.currentUsername) {
            this.props.searchLogs(this.props.currentChannel, this.props.currentUsername, this.state.year, this.state.month);
        }
    }

    render() {
        const completions = this.props.channels
            .filter(channel => channel.name.includes(this.props.currentChannel));

        const autocompletions = [];
        for (const completion of completions) {
            autocompletions.push(<li key={completion.userID} onClick={() => this.setChannel(completion.name)}>{completion.name}</li>);
        }

        return (
            <form className="filter" autoComplete="off" onSubmit={this.onSubmit}>
                <AutocompleteInput placeholder="pajlada" onChange={this.onChannelChange} value={this.props.currentChannel} autocompletions={this.props.channels} />
                <input
                    type="text"
                    placeholder="gempir"
                    onChange={this.onUsernameChange}
                    value={this.props.currentUsername}
                />
                <button type="submit" className="show-logs">Show logs</button>
            </form>
        )
    }

    onChannelChange = (e) => {
        this.props.dispatch(setCurrent(e.target.value, this.props.currentUsername));
    }

    onUsernameChange = (e) => {
        this.props.dispatch(setCurrent(this.props.currentChannel, e.target.value));
    }

    onYearChange = (e) => {
        this.setState({ year: e.target.value });
    }

    onMonthChange = (e) => {
        this.setState({ month: e.target.value });
    }

    onSubmit = (e) => {
        e.preventDefault();

        const url = new URL(window.location.href);
        const params = new URLSearchParams(url.search);
        params.set('channel', this.props.currentChannel);
        params.set('username', this.props.currentUsername);
        window.location.search = params.toString();

        this.props.searchLogs(this.props.currentChannel, this.props.currentUsername, this.state.year, this.state.month);
    }
}

const mapStateToProps = (state) => ({
    channels: state.channels,
    currentChannel: state.currentChannel,
    currentUsername: state.currentUsername
});

export default connect(mapStateToProps)(Filter);