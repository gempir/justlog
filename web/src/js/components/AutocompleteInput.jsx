import React, { Component } from "react";

export default class AutocompleteInput extends Component {

    state = {
        focused: false,
        previewValue: null,
        selectedIndex: -1,
    };

    input;

    render() {
        return <div className="AutocompleteInput">
            <input
                type="text"
                ref={el => this.input = el}
                placeholder={this.props.placeholder}
                onChange={this.handleChange}
                onFocus={() => this.setState({ focused: true })}
                onBlur={this.handleBlur}
                value={this.state.previewValue ?? this.props.value}
                onKeyDown={this.handleKeyDown}
            />
            {this.state.focused && <ul>
                {this.getAutocompletions()
                    .map((completion, index) =>
                        <li className={index === this.state.selectedIndex ? "selected" : ""} key={completion} onClick={() => this.handleClick(completion)} onMouseDown={e => e.preventDefault()}>
                            {completion}
                        </li>
                    )}
            </ul>}
        </div>
    }

    getAutocompletions = () => {
        return this.props.autocompletions
            .filter(completion => completion.includes(this.props.value) && this.props.value !== "")
            .sort();
    }

    handleBlur = () => {
        if (this.state.selectedIndex !== -1) {
            this.props.onChange(this.getAutocompletions()[this.state.selectedIndex]);
        }

        this.setState({
            focused: false, previewValue: null, selectedIndex: -1
        });
    }

    handleChange = (e) => {
        this.props.onChange(e.target.value);
    }

    handleKeyDown = (e) => {
        if (["Backspace"].includes(e.key)) {
            if (this.state.selectedIndex !== -1 && this.getAutocompletions().length > 0) {
                this.props.onChange(this.getAutocompletions()[this.state.selectedIndex]);
                this.setState({ previewValue: null, selectedIndex: -1 });
            }
            return;
        }

        if (!["ArrowDown", "ArrowUp", "Enter"].includes(e.key)) {
            this.setState({ previewValue: null });
            return;
        }
        e.preventDefault();

        if (e.key === "ArrowDown") {
            if (this.state.selectedIndex === this.props.autocompletions.length - 1) {
                return;
            }

            const newIndex = this.state.selectedIndex + 1;

            this.setState({
                selectedIndex: newIndex,
                previewValue: this.getAutocompletions()[newIndex]
            });
        } else if (e.key === "ArrowUp") {
            if (this.state.selectedIndex === -1) {
                return;
            }

            const newIndex = this.state.selectedIndex - 1;

            this.setState({
                selectedIndex: newIndex,
                previewValue: this.getAutocompletions()[newIndex]
            });
        } else if (e.key === "Enter") {
            if (this.state.selectedIndex === -1) {
                this.props.onSubmit();
                return;
            }

            this.props.onChange(this.state.previewValue);
            this.props.onSubmit();
            this.setState({
                selectedIndex: -1,
                previewValue: null,
            });
        }
    };

    handleClick = (completion) => {
        this.props.onChange(completion);
        this.input.blur();
        this.props.onSubmit();
    };
}