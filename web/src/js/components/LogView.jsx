import React, { Component } from 'react';
import { connect } from "react-redux";
import twitchEmotes from "../emotes/twitch";
import reactStringReplace from "react-string-replace";

class LogView extends Component {

	static LOAD_LIMIT = 100;

	state = {
		limitLoad: true,
	};

	render() {
		return (
			<div className={"log-view"}>
				{this.getLogs().map((value, key) =>
					<div key={key} className="line" onClick={() => this.setState({})}>
						<span id={value.timestamp} className="timestamp">{this.formatDate(value.timestamp)}</span>{this.renderMessage(value.text)}
					</div>
				)}
				{this.getLogs().length > 0 && this.state.limitLoad && <button className={"load-all"} raised primary onClick={() => this.setState({ ...this.state, limitLoad: false })}>Load all</button>}
				{this.props.loading && <div>loading</div>}
			</div>
		);
	}

	getLogs = () => {
		if (this.state.limitLoad) {
			return this.props.messages.slice(this.props.messages.length - LogView.LOAD_LIMIT, this.props.messages.length).reverse();
		} else {
			return this.props.messages.reverse();
		}
	};

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
		return new Date(timestamp).format("YYYY-MM-DD HH:mm:ss UTC");
	}

	buildTwitchEmote = (id) => {
		return `https://static-cdn.jtvnw.net/emoticons/v1/${id}/1.0`;
	}
}

const mapStateToProps = (state) => {
	return {
		messages: state.logs.messages,
		loading: state.loading
	};
};

export default connect(mapStateToProps)(LogView);