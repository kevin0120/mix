import React from 'react';

const Color = require('color');

export default class Circle extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      progressColor: Color(this.props.startColor)
    };
    this.setSizes(props.lineWidth);
  }

  componentWillReceiveProps(nextProps) {
    const mixRate = 1 - (nextProps.progress * 0.6 + 0.4);
    this.setState({
      progressColor: Color(nextProps.startColor).mix(
        Color(nextProps.endColor),
        mixRate
      )
    });
  }

  setSizes = lineWidth => {
    this.viewSize = 1000;
    this.radius = (this.viewSize - lineWidth * 3) / 2;
    this.diameter = Math.round(Math.PI * this.radius * 2);
    this.getOffset = (val = 0) =>
      Math.round(((100 - val * 100) / 100) * this.diameter);
  };

  render() {
    const {
      text,
      progress,
      size,
      bgColor,
      lineWidth,
      animate,
      animationDuration,
      roundedStroke,
      responsive,
      onAnimationEnd,
      textStyle,
      textColor
    } = this.props;
    const { viewSize, radius, diameter, getOffset } = this;
    const strokeDashoffset = getOffset(progress);
    const transition = animate
      ? `stroke-dashoffset ${animationDuration}s linear, stroke ${animationDuration}s linear`
      : null;
    const strokeLinecap = roundedStroke ? 'round' : 'butt';
    const svgSize = responsive ? '100%' : size;
    const dashLength = diameter;
    return (
      <svg
        width={svgSize}
        height={svgSize}
        viewBox={`${-(viewSize / 2 - radius)} ${-(
          viewSize / 2 -
          radius
        )} ${viewSize} ${viewSize}`}
      >
        <circle
          stroke={bgColor}
          cx={radius}
          cy={radius}
          r={radius}
          strokeWidth={lineWidth * 0.9}
          fill="none"
        />
        <circle
          transform={`rotate(-90 ${radius} ${radius})`}
          stroke={`rgb(${this.state.progressColor.red()},${this.state.progressColor.green()},${this.state.progressColor.blue()})`}
          cx={radius}
          cy={radius}
          r={radius}
          strokeDasharray={`${dashLength}`}
          strokeWidth={lineWidth * 1.1}
          strokeDashoffset={`${dashLength}`}
          strokeLinecap={strokeLinecap}
          fill="none"
          style={{
            strokeDashoffset,
            transition
          }}
          onTransitionEnd={onAnimationEnd}
        />

        <circle
          stroke="none"
          cx={radius}
          cy={0}
          r={lineWidth * 1.5}
          fill={`rgb(${this.state.progressColor.red()},${this.state.progressColor.green()},${this.state.progressColor.blue()})`}
          transform={`rotate(${progress * 360} ${radius} ${radius})`}
          style={{
            transition: `transform ${animationDuration}s linear, fill ${animationDuration}s linear`
          }}
        />
        <text
          style={textStyle}
          fill={textColor}
          x={radius}
          y={radius}
          textAnchor="middle"
          dominantBaseline="central"
        >
          {text}
        </text>
      </svg>
    );
  }
}
