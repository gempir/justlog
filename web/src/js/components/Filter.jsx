import React, { Component } from "react";
import {connect} from "react-redux";
import setCurrent from "./../actions/setCurrent";

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
        return (
            <form className="filter" autoComplete="off" onSubmit={this.onSubmit}>
                <div className="channel-wrapper">
                <input
                    type="text"
                    placeholder="pajlada"
                    onChange={this.onChannelChange}
                    value={this.props.currentChannel}
                />
                <ul className="channel-autocomplete">
                    {this.props.channels
                    .filter(channel => channel.name.includes(this.props.currentChannel))
                    .map(channel => <li key={channel.userID} onClick={() => this.setChannel(channel.name)}>{channel.name}</li>)}
                </ul>
                </div>
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

    setChannel = channel => this.props.dispatch(setCurrent(channel, this.props.currentUsername));

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
        this.props.searchLogs(this.props.currentChannel, this.props.currentUsername, this.state.year, this.state.month);
    }
}

const mapStateToProps = (state) => ({
    channels: state.channels,
    currentChannel: state.currentChannel,
    currentUsername: state.currentUsername
});

export default connect(mapStateToProps)(Filter);