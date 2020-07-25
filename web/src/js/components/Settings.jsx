import React from "react";
import { connect } from "react-redux";
import setSettings from "../actions/setSettings";

const Settings = (props) => {

    const toggleSetting = (setting) => {
        props.dispatch(setSettings({ ...props.settings, [setting]: !props.settings[setting] }))
    }

    return <div className="settings">
        <div className="setting">
            <input type="checkbox" id="showEmotes" checked={props.settings.showEmotes} onChange={() => toggleSetting("showEmotes")} />
            <label htmlFor="showEmotes">show emotes</label>
        </div>
        <div className="setting">
            <input type="checkbox" id="showDisplayName" checked={props.settings.showDisplayName} onChange={() => toggleSetting("showDisplayName")} />
            <label htmlFor="showDisplayName">show displayName</label>
        </div>
    </div>
}

export default connect(state => ({
    settings: state.settings
}))(Settings);