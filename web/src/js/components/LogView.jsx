import React, { Component } from 'react';
import { connect } from "react-redux";
import loadLogs from '../actions/loadLogs';
import LoadingSpinner from "./LoadingSpinner";
import AnimateHeight from "./AnimateHeight";
import { parse } from "irc-message";

class LogView extends Component {

	state = {
		loading: false,
		height: 0,
		buttonText: "load",
	};

	componentDidMount() {
		if (this.props.log.messages.length > 0) {
			setTimeout(() => this.setState({ height: 'auto' }), 10);
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
					{this.props.log.messages.map((value, key) =>
						<div key={key} className="line">
							<span id={value.timestamp} className="timestamp">{this.formatDate(value.timestamp)}</span>{this.renderMessage(value)}
						</div>
					)}
				</AnimateHeight>
			</div>
		);
	}

	renderMessage = (value) => {
		const msgObj = parse(value.raw);
		let message = value.text;

		const replacements = [];

		if (msgObj.tags.emotes && msgObj.tags.emotes !== true) {
			for (const emoteString of msgObj.tags.emotes.split("/")) {
				if (typeof emoteString !== "string") {
					continue;
				}
				const [emoteId, occurences] = emoteString.split(":");
				if (typeof occurences !== "string") {
					continue;
				}
				for (const occurence of occurences.split(",")) {
					const [start, end] = occurence.split("-");
					replacements.push({start: Number(start), end: Number(end) + 1, emoteId});
				}
			}
		}

		replacements.sort((a, b) => {
			if (a.start > b.start) {
				return -1;
			}
			if (a.start < b.start) {
				return 1;
			}
			return 0;
		})

		for (const replacement of replacements) {
			const emote = `<img src="${this.buildTwitchEmote(replacement.emoteId)}" alt="${replacement.emoteId}" />`;
			message = message.slice(0, replacement.start) + emote + message.slice(replacement.end);
		}

		if (this.props.bttvChannelEmotes) {
			for (const emote of this.props.bttvChannelEmotes.emotes) {
				const regex = new RegExp(`\\b(${emote.code})\\b`, "g");

				message = message.replace(regex, `<img src="${this.buildBttvEmote(emote.id)}" alt="${emote.id}" />`);
			}
		}

		if (this.props.bttvEmotes) {
			for (const emote of this.props.bttvEmotes) {
				const regex = new RegExp(`\\b(${emote.code})\\b`, "g");

				message = message.replace(regex, `<img src="${this.buildBttvEmote(emote.id)}" alt="${emote.id}" />`);
			}
		}

		return (
			<p dangerouslySetInnerHTML={{ __html: message }}>
			</p>
		);
	}

	loadLog = () => {
		this.setState({ loading: true });
		this.props.dispatch(loadLogs(null, null, this.props.log.year, this.props.log.month)).then(() => this.setState({ loading: false })).catch(() => {
			this.setState({ loading: false, buttonText: "not found" });
		});
	}

	formatDate = (timestamp) => {
		return new Date(timestamp).toUTCString();
	}

	buildTwitchEmote = (id) => {
		return `https://static-cdn.jtvnw.net/emoticons/v1/${id}/1.0`;
	}

	buildBttvEmote = (id) => {
		return `https://cdn.betterttv.net/emote/${id}/1x`;
	}
}
const mapStateToProps = (state) => {
	return {
		bttvChannelEmotes: state.bttvChannelEmotes,
		bttvEmotes: state.bttvEmotes,
	};
};

export default connect(mapStateToProps)(LogView);