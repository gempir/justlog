import React, { Component } from 'react';
import { connect } from "react-redux";
import twitchEmotes from "../emotes/twitch";
import reactStringReplace from "react-string-replace";
import loadLogs from '../actions/loadLogs';

class LogView extends Component {

	render() {
		if (this.props.log.loaded === false) {
			return <div className="log-view not-loaded" onClick={this.loadLog}>
				<span>{this.props.log.getTitle()}</span>
				<button>load</button>
			</div>;
		}

		return (
			<div className={"log-view"}>
				{this.props.log.messages.reverse().map((value, key) =>
					<div key={key} className="line" onClick={() => this.setState({})}>
						<span id={value.timestamp} className="timestamp">{this.formatDate(value.timestamp)}</span>{this.renderMessage(value.text)}
					</div>
				)}
			</div>
		);
	}

	renderMessage = (message) => {
		for (let emoteCode in twitchEmotes) {
			const regex = new RegExp(`(?:^|\ )(${emoteCode})(?:$|\ )`);

			message = reactStringReplace(message, regex, (match, i) => (
				<img key={i} src={this.buildTwitchEmote(twitchEmotes[emoteCode].id)} alt={match} />
			));
		}

		return (
			<p>
				{message}
			</p>
		);
	}

	loadLog = () => {
		this.props.dispatch(loadLogs(null, null, this.props.log.year, this.props.log.month));
	}

	formatDate = (timestamp) => {
		return new Date(timestamp).toUTCString();
	}

	buildTwitchEmote = (id) => {
		return `https://static-cdn.jtvnw.net/emoticons/v1/${id}/1.0`;
	}
}
const mapStateToProps = (state) => {
    return {

    };
};

export default connect(mapStateToProps)(LogView);