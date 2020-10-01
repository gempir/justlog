import React, { Component } from 'react';
import { connect } from "react-redux";
import loadLog from '../actions/loadLog';
import LoadingSpinner from "./LoadingSpinner";
import AnimateHeight from "./AnimateHeight";
import { parse } from "irc-message";
import loadBttvChannelEmotes from "../actions/loadBttvChannelEmotes";
import loadFfzChannelEmotes from "../actions/loadFfzChannelEmotes";
import Txt from '../icons/Txt';

class LogView extends Component {

	state = {
		loading: false,
		height: 0,
		buttonText: "load",
	};

	loadedBttvEmotes = false;
	loadedFfzEmotes = false;

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
				<a className="txt" target="__blank" href={`/channel/${this.props.channel}/user/${this.props.username}/${this.props.log.year}/${this.props.log.month}`}><Txt /></a>
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
		let message = this.htmlencode(value.text);

		if (this.props.settings.showEmotes) {
			message = [...message];

			if (this.loadedBttvEmotes !== msgObj.tags["room-id"]) {
				this.props.dispatch(loadBttvChannelEmotes(msgObj.tags["room-id"]));
				this.loadedBttvEmotes = msgObj.tags["room-id"];
			}

			if (this.loadedFfzEmotes !== msgObj.tags["room-id"]) {
				this.props.dispatch(loadFfzChannelEmotes(msgObj.tags["room-id"]));
				this.loadedFfzEmotes = msgObj.tags["room-id"];
			}

			const replacements = {};

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
						replacements[Number(start)] = { length: Number(end) - Number(start), emoteId };
					}
				}
			}

			message.forEach((char, key) => {
				if (typeof replacements[key] !== "undefined") {
					const emote = `<img src="${this.buildTwitchEmote(replacements[key].emoteId)}" alt="${message.slice(key, replacements[key].end)}" />`;
					message[key] = emote;
					for (let i = key + 1; i < (key + replacements[key].length + 1); i++) {
						message[i] = "";
					}
				}
			});

			message = message.join("");

			const replacedEmoteCodes = {};

			if (this.props.bttvChannelEmotes) {
				for (const emote of [...this.props.bttvChannelEmotes.channelEmotes, ...this.props.bttvChannelEmotes.sharedEmotes]) {
					if (replacedEmoteCodes[emote.code]) {
						continue;
					}
					replacedEmoteCodes[emote.code] = true;

					const regex = new RegExp(`\\b(${emote.code})\\b`, "g");
					message = message.replace(regex, `<img src="${this.buildBttvEmote(emote.id)}" alt="${emote.code}" />`);
				}
			}

			if (this.props.ffzChannelEmotes) {
				for (const emote of Object.values(this.props.ffzChannelEmotes.sets).map(set => set.emoticons).flat()) {
					if (replacedEmoteCodes[emote.name]) {
						continue;
					}
					replacedEmoteCodes[emote.name] = true;

					const regex = new RegExp(`\\b(${emote.name})\\b`, "g");
					message = message.replace(regex, `<img src="${emote.urls[1]}" alt="${emote.name}" />`);
				}
			}

			if (this.props.bttvEmotes) {
				for (const emote of this.props.bttvEmotes) {
					if (replacedEmoteCodes[emote.code]) {
						continue;
					}
					replacedEmoteCodes[emote.code] = true;

					const regex = new RegExp(`\\b(${emote.code})\\b`, "g");
					message = message.replace(regex, `<img src="${this.buildBttvEmote(emote.id)}" alt="${emote.code}" />`);
				}
			}
		}

		if (this.props.settings.showDisplayName) {
			message = `${value.displayName}: ${message}`;
		}

		return (
			<p dangerouslySetInnerHTML={{ __html: message }}>
			</p>
		);
	}

	loadLog = () => {
		this.setState({ loading: true });
		this.props.dispatch(loadLog(null, null, this.props.log.year, this.props.log.month)).then(() => this.setState({ loading: false })).catch(err => {
			console.error(err);
			this.setState({ loading: false, buttonText: err.message });
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

	/**
	 * HTML entities encode
	 *
	 * @param {string} str Input text
	 * @return {string} Filtered text
	 */
	htmlencode = (str) => {
		var div = document.createElement('div');
		div.appendChild(document.createTextNode(str));
		return div.innerHTML;
	}
}
const mapStateToProps = (state) => {
	return {
		bttvChannelEmotes: state.bttvChannelEmotes,
		bttvEmotes: state.bttvEmotes,
		ffzChannelEmotes: state.ffzChannelEmotes,
		settings: state.settings,
	};
};

export default connect(mapStateToProps)(LogView);