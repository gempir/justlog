import React, { Component } from 'react';
import { connect } from "react-redux";
import twitchEmotes from "../emotes/twitch";
import reactStringReplace from "react-string-replace";

class LogView extends Component {

	render() {
		const oldLogs = [];
		
		for (let month = this.props.month - 1; month >= 1; month--) {
			oldLogs.push(month)
		}

		console.log(oldLogs);

		return (
			<div className={"log-view"}>
				{this.props.messages.reverse().map((value, key) =>
					<div key={key} className="line" onClick={() => this.setState({})}>
						<span id={value.timestamp} className="timestamp">{this.formatDate(value.timestamp)}</span>{this.renderMessage(value.text)}
					</div>
				)}
				{this.props.loading && <div>loading</div>}
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

	formatDate = (timestamp) => {
		return new Date(timestamp).toUTCString();
	}

	buildTwitchEmote = (id) => {
		return `https://static-cdn.jtvnw.net/emoticons/v1/${id}/1.0`;
	}
}

const mapStateToProps = (state) => {
	return {
		month: state.month,
		messages: state.logs.messages,
		loading: state.loading
	};
};

export default connect(mapStateToProps)(LogView);