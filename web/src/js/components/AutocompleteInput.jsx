import React, { Component } from "react";

export default class AutocompleteInput extends Component {

    state = {
        focused: false,
    };

    input;

    render() {
        return <div className="AutocompleteInput">
            <input
                type="text"
                ref={el => this.input = el}
                placeholder={this.props.placeholder}
                onChange={e => this.props.onChange(e.target.value)}
                onFocus={() => this.setState({ focused: true })}
                onBlur={() => this.setState({ focused: false })}
                value={this.props.value}
            />
            {this.state.focused && <ul>
                {this.props.autocompletions
                    .filter(completion => completion.includes(this.props.value))
                    .map(completion =>
                        <li key={completion} onClick={() => this.handleClick(completion)} onMouseDown={e => e.preventDefault()}>
                            {completion}
                        </li>
                    )}
            </ul>}
        </div>
    }

    handleClick = (completion) => {
        this.props.onChange(completion);
        this.input.blur();
        this.props.onAutocompletionClick(completion);
    }
}