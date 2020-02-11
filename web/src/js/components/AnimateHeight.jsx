// MIT License

// Copyright (c) 2017 - today Stanko Tadić

// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import React from 'react';

const ANIMATION_STATE_CLASSES = {
    animating: 'rah-animating',
    animatingUp: 'rah-animating--up',
    animatingDown: 'rah-animating--down',
    animatingToHeightZero: 'rah-animating--to-height-zero',
    animatingToHeightAuto: 'rah-animating--to-height-auto',
    animatingToHeightSpecific: 'rah-animating--to-height-specific',
    static: 'rah-static',
    staticHeightZero: 'rah-static--height-zero',
    staticHeightAuto: 'rah-static--height-auto',
    staticHeightSpecific: 'rah-static--height-specific',
};

const PROPS_TO_OMIT = [
    'animateOpacity',
    'animationStateClasses',
    'applyInlineTransitions',
    'children',
    'contentClassName',
    'delay',
    'duration',
    'easing',
    'height',
    'onAnimationEnd',
    'onAnimationStart',
];

function omit(obj, ...keys) {
    if (!keys.length) {
        return obj;
    }

    const res = {};
    const objectKeys = Object.keys(obj);

    for (let i = 0; i < objectKeys.length; i++) {
        const key = objectKeys[i];

        if (keys.indexOf(key) === -1) {
            res[key] = obj[key];
        }
    }

    return res;
}

// Start animation helper using nested requestAnimationFrames
function startAnimationHelper(callback) {
    const requestAnimationFrameIDs = [];

    requestAnimationFrameIDs[0] = requestAnimationFrame(() => {
        requestAnimationFrameIDs[1] = requestAnimationFrame(() => {
            callback();
        });
    });

    return requestAnimationFrameIDs;
}

function cancelAnimationFrames(requestAnimationFrameIDs) {
    requestAnimationFrameIDs.forEach(id => cancelAnimationFrame(id));
}

function isNumber(n) {
    return !isNaN(parseFloat(n)) && isFinite(n);
}

function isPercentage(height) {
    // Percentage height
    return typeof height === 'string' &&
        height.search('%') === height.length - 1 &&
        isNumber(height.substr(0, height.length - 1));
}

function runCallback(callback, params) {
    if (callback && typeof callback === 'function') {
        callback(params);
    }
}

const AnimateHeight = class extends React.Component {
    constructor(props) {
        super(props);

        this.animationFrameIDs = [];

        let height = 'auto';
        let overflow = 'visible';

        if (isNumber(props.height)) {
            // If value is string "0" make sure we convert it to number 0
            height = props.height < 0 || props.height === '0' ? 0 : props.height;
            overflow = 'hidden';
        } else if (isPercentage(props.height)) {
            // If value is string "0%" make sure we convert it to number 0
            height = props.height === '0%' ? 0 : props.height;
            overflow = 'hidden';
        }

        this.animationStateClasses = { ...ANIMATION_STATE_CLASSES, ...props.animationStateClasses };

        const animationStateClasses = this.getStaticStateClasses(height);

        this.state = {
            animationStateClasses,
            height,
            overflow,
            shouldUseTransitions: false,
        };
    }

    componentDidMount() {
        const { height } = this.state;

        // Hide content if height is 0 (to prevent tabbing into it)
        // Check for contentElement is added cause this would fail in tests (react-test-renderer)
        // Read more here: https://github.com/Stanko/react-animate-height/issues/17
        if (this.contentElement && this.contentElement.style) {
            this.hideContent(height);
        }
    }

    componentDidUpdate(prevProps, prevState) {
        const {
            delay,
            duration,
            height,
            onAnimationEnd,
            onAnimationStart,
        } = this.props;

        // Check if 'height' prop has changed
        if (this.contentElement && height !== prevProps.height) {
            // Remove display: none from the content div
            // if it was hidden to prevent tabbing into it
            this.showContent(prevState.height);

            // Cache content height
            this.contentElement.style.overflow = 'hidden';
            const contentHeight = this.contentElement.offsetHeight;
            this.contentElement.style.overflow = '';

            // set total animation time
            const totalDuration = duration + delay;

            let newHeight = null;
            const timeoutState = {
                height: null, // it will be always set to either 'auto' or specific number
                overflow: 'hidden',
            };
            const isCurrentHeightAuto = prevState.height === 'auto';


            if (isNumber(height)) {
                // If value is string "0" make sure we convert it to number 0
                newHeight = height < 0 || height === '0' ? 0 : height;
                timeoutState.height = newHeight;
            } else if (isPercentage(height)) {
                // If value is string "0%" make sure we convert it to number 0
                newHeight = height === '0%' ? 0 : height;
                timeoutState.height = newHeight;
            } else {
                // If not, animate to content height
                // and then reset to auto
                newHeight = contentHeight; // TODO solve contentHeight = 0
                timeoutState.height = 'auto';
                timeoutState.overflow = null;
            }

            if (isCurrentHeightAuto) {
                // This is the height to be animated to
                timeoutState.height = newHeight;

                // If previous height was 'auto'
                // set starting height explicitly to be able to use transition
                newHeight = contentHeight;
            }



            // Animation classes
            let animationStateClasses = [this.animationStateClasses.animating];

            if (prevProps.height === 'auto' || height < prevProps.height) {
                animationStateClasses.push(this.animationStateClasses.animatingUp);
            }
            if (height === 'auto' || height > prevProps.height) {
                animationStateClasses.push(this.animationStateClasses.animatingDown);
            }
            if (timeoutState.height === 0) {
                animationStateClasses.push(this.animationStateClasses.animatingToHeightZero);
            }
            if (timeoutState.height === 'auto') {
                animationStateClasses.push(this.animationStateClasses.animatingToHeightAuto);
            }
            if (timeoutState.height > 0) {
                animationStateClasses.push(this.animationStateClasses.animatingToHeightSpecific);
            }

            animationStateClasses = animationStateClasses.join(" ");

            // Animation classes to be put after animation is complete
            const timeoutAnimationStateClasses = this.getStaticStateClasses(timeoutState.height);

            // Set starting height and animating classes
            // We are safe to call set state as it will not trigger infinite loop
            // because of the "height !== prevProps.height" check
            this.setState({ // eslint-disable-line react/no-did-update-set-state
                animationStateClasses,
                height: newHeight,
                overflow: 'hidden',
                // When animating from 'auto' we first need to set fixed height
                // that change should be animated
                shouldUseTransitions: !isCurrentHeightAuto,
            });

            // Clear timeouts
            clearTimeout(this.timeoutID);
            clearTimeout(this.animationClassesTimeoutID);

            if (isCurrentHeightAuto) {
                // When animating from 'auto' we use a short timeout to start animation
                // after setting fixed height above
                timeoutState.shouldUseTransitions = true;

                cancelAnimationFrames(this.animationFrameIDs);
                this.animationFrameIDs = startAnimationHelper(() => {
                    this.setState(timeoutState);

                    // ANIMATION STARTS, run a callback if it exists
                    runCallback(onAnimationStart, { newHeight: timeoutState.height });
                });

                // Set static classes and remove transitions when animation ends
                this.animationClassesTimeoutID = setTimeout(() => {
                    this.setState({
                        animationStateClasses: timeoutAnimationStateClasses,
                        shouldUseTransitions: false,
                    });

                    // ANIMATION ENDS
                    // Hide content if height is 0 (to prevent tabbing into it)
                    this.hideContent(timeoutState.height);
                    // Run a callback if it exists
                    runCallback(onAnimationEnd, { newHeight: timeoutState.height });
                }, totalDuration);
            } else {
                // ANIMATION STARTS, run a callback if it exists
                runCallback(onAnimationStart, { newHeight });

                // Set end height, classes and remove transitions when animation is complete
                this.timeoutID = setTimeout(() => {
                    timeoutState.animationStateClasses = timeoutAnimationStateClasses;
                    timeoutState.shouldUseTransitions = false;

                    this.setState(timeoutState);

                    // ANIMATION ENDS
                    // If height is auto, don't hide the content
                    // (case when element is empty, therefore height is 0)
                    if (height !== 'auto') {
                        // Hide content if height is 0 (to prevent tabbing into it)
                        this.hideContent(newHeight); // TODO solve newHeight = 0
                    }
                    // Run a callback if it exists
                    runCallback(onAnimationEnd, { newHeight });
                }, totalDuration);
            }
        }
    }

    componentWillUnmount() {
        cancelAnimationFrames(this.animationFrameIDs);

        clearTimeout(this.timeoutID);
        clearTimeout(this.animationClassesTimeoutID);

        this.timeoutID = null;
        this.animationClassesTimeoutID = null;
        this.animationStateClasses = null;
    }

    showContent(height) {
        if (height === 0) {
            this.contentElement.style.display = '';
        }
    }

    hideContent(newHeight) {
        if (newHeight === 0) {
            this.contentElement.style.display = 'none';
        }
    }

    getStaticStateClasses(height) {
        const classes = [this.animationStateClasses.static];

        if (height === 0) {
            classes.push(this.animationStateClasses.staticHeightZero);
        }
        if (height > 0) {
            classes.push(this.animationStateClasses.staticHeightSpecific);
        }
        if (height === 'auto') {
            classes.push(this.animationStateClasses.staticHeightAuto);
        }

        return classes.join(" ");
    }

    render() {
        const {
            animateOpacity,
            applyInlineTransitions,
            children,
            className,
            contentClassName,
            duration,
            easing,
            delay,
            style,
        } = this.props;
        const {
            height,
            overflow,
            animationStateClasses,
            shouldUseTransitions,
        } = this.state;


        const componentStyle = {
            ...style,
            height,
            overflow: overflow || style.overflow,
        };

        if (shouldUseTransitions && applyInlineTransitions) {
            componentStyle.transition = `height ${duration}ms ${easing} ${delay}ms`;

            // Include transition passed through styles
            if (style.transition) {
                componentStyle.transition = `${style.transition}, ${componentStyle.transition}`;
            }

            // Add webkit vendor prefix still used by opera, blackberry...
            componentStyle.WebkitTransition = componentStyle.transition;
        }

        const contentStyle = {};

        if (animateOpacity) {
            contentStyle.transition = `opacity ${duration}ms ${easing} ${delay}ms`;
            // Add webkit vendor prefix still used by opera, blackberry...
            contentStyle.WebkitTransition = contentStyle.transition;

            if (height === 0) {
                contentStyle.opacity = 0;
            }
        }

        let componentClasses = [animationStateClasses];
        if (className) {
            componentClasses.push(className);
        }
        componentClasses = componentClasses.join(" ");

        return (
            <div
                {...omit(this.props, ...PROPS_TO_OMIT)}
                aria-hidden={height === 0}
                className={componentClasses}
                style={componentStyle}
            >
                <div
                    className={contentClassName}
                    style={contentStyle}
                    ref={el => this.contentElement = el}
                >
                    {children}
                </div>
            </div>
        );
    }
};

AnimateHeight.defaultProps = {
    animateOpacity: false,
    animationStateClasses: ANIMATION_STATE_CLASSES,
    applyInlineTransitions: true,
    duration: 250,
    delay: 0,
    easing: 'ease',
    style: {},
};

export default AnimateHeight;