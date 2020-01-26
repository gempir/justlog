import React, { Component } from 'react';
import { Autocomplete, TextField, Button, SelectField } from 'react-md';
import moment from "moment";

export default class Filter extends Component {

    constructor(props) {
        super(props);

        this.state = {
            channel: "",
            username: "",
            year: moment().year(),
            month: moment().format("M")
         }
    }

    render() {
		return (
            <form className="filter" autoComplete="off" onSubmit={this.onSubmit}>
                <Autocomplete
                    id="channel"
                    label="Channel"
                    placeholder="forsen"
                    onChange={this.onChannelChange}
                    onAutocomplete={this.onChannelChange}
                    data={this.props.channels.map(obj => obj.name)}
                />
                <TextField
                    id="username"
                    label="Username"
                    lineDirection="center"
                    onChange={this.onUsernameChange}
                    placeholder="gempir"
                />
                <SelectField
                    id="year"
                    label="Year"
                    defaultValue={this.state.year}
                    menuItems={[moment().year(), moment().subtract(1, "year").year(),  moment().subtract(2, "year").year()]}
                    onChange={this.onYearChange}
                    value={this.state.year}
                />
                <SelectField
                    id="month"
                    label="Month"
                    defaultValue={this.state.month}
                    menuItems={["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]}
                    onChange={this.onMonthChange}
                    value={this.state.month}
                />
                <Button flat primary swapTheming type="submit" className="show-logs">Show logs</Button>
            </form>
		)
    }

    onChannelChange = (value) => {
        this.setState({...this.state, channel: value});
    } 

    onUsernameChange = (value) => {
        this.setState({...this.state, username: value});
    }
    
    onYearChange = (value) => {
        this.setState({...this.state, year: value});
    } 

    onMonthChange = (value) => {
        this.setState({...this.state, month: moment().month(value).format("M")});
    }

    onSubmit = (e) => {
        e.preventDefault();
        this.props.searchLogs(this.state.channel, this.state.username, this.state.year, this.state.month);
    }
}