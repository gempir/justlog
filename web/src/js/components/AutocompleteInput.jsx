import React, { Component } from "react";

export default class AutocompleteInput extends Component {
    render() {            
        return <div className="AutocompleteInput">
            <input
                type="text"
                placeholder={this.props.placeholder}
                onChange={this.props.onChange}
                value={this.props.value}
            />
            <ul>
                {this.props.autocompletions
                .filter(channel => channel.name.includes(this.props.value))
                .map(channel => <li key={channel.userID} onClick={() => this.setChannel(channel.name)}>{channel.name}</li>)}
            </ul>
        </div>
    }
}