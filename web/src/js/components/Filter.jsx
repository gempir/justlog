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

    username;

    componentDidMount() {
        if (this.props.currentChannel && this.props.currentUsername) {
            this.props.searchLogs(this.props.currentChannel, this.props.currentUsername, this.state.year, this.state.month);
        }
    }

    render() {
        return (
            <form className="filter" autoComplete="off" onSubmit={this.onSubmit}>
                <AutocompleteInput placeholder="pajlada" onChange={this.onChannelChange} value={this.props.currentChannel} onAutocompletionClick={() => this.username.focus()} autocompletions={this.props.channels.map(channel => channel.name)} />
                <input
                    ref={el => this.username = el}
                    type="text"
                    placeholder="gempir"
                    onChange={this.onUsernameChange}
                    value={this.props.currentUsername}
                />
                <button type="submit" className="show-logs">Show logs</button>
            </form>
        )
    }

    onChannelChange = (channel) => {
        this.props.dispatch(setCurrent(channel, this.props.currentUsername));
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