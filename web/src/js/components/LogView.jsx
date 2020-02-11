import React, { Component } from 'react';
import { connect } from "react-redux";
import twitchEmotes from "../emotes/twitch";
import reactStringReplace from "react-string-replace";
import loadLogs from '../actions/loadLogs';
import LoadingSpinner from "./LoadingSpinner";
import AnimateHeight from "./AnimateHeight";

class LogView extends Component {

	state = {
		loading: false,
		height: 0,
		buttonText: "load"
	};

	componentDidMount() {
		if (this.props.log.messages.length > 0) {
			setTimeout(() => this.setState({height: 'auto'}), 10);
		}
	}

	componentDidUpdate(prevProps) {
		if (prevProps.log.messages.length !== this.props.log.messages.length) {
			this.setState({
				height: 'auto'
			});
		}
	}

	render() {
		if (this.props.log.loaded === false) {
			return <div className="log-view not-loaded" onClick={this.loadLog}>
				<span>{this.props.log.getTitle()}</span>
				<button>{this.state.loading ? <LoadingSpinner /> : this.state.buttonText}</button>
			</div>;
		}

		return (
			<div className={"log-view"}>
				<AnimateHeight duration={500} easing={"ease-in-out"} height={this.state.height} animateOpacity>
				{this.props.log.messages.reverse().map((value, key) =>
					<div key={key} className="line">
						<span id={value.timestamp} className="timestamp">{this.formatDate(value.timestamp)}</span>{this.renderMessage(value.text)}
					</div>
				)}
				</AnimateHeight>
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
		this.setState({loading: true});
		this.props.dispatch(loadLogs(null, null, this.props.log.year, this.props.log.month)).then(() => this.setState({loading: false})).catch(() => {
			this.setState({loading: false, buttonText: "not found"});
		});
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