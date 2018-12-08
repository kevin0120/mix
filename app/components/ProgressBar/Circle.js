import React from 'react';
import PropTypes from 'prop-types';

const Color = require('color');


export default class Circle extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      progressColor: Color(props.startColor)
    };
    this.setSizes(props.lineWidth);
  }

  componentWillReceiveProps(nextProps) {
    const mixRate = 1 - (nextProps.progress * 0.6 + 0.4);
    this.setState({
      progressColor: Color(nextProps.startColor).mix(Color(nextProps.endColor), mixRate)
    });
  }

  setSizes = (lineWidth) => {
    this.viewSize = 1000;
    this.radius = (this.viewSize - lineWidth * 3) / 2;
    this.diameter = Math.round(Math.PI * this.radius * 2);
    this.getOffset = (val = 0) => Math.round((100 - val * 100) / 100 * this.diameter);
  };

  render() {
    const { text, progress, size, bgColor, lineWidth, animate, animationDuration, roundedStroke, responsive, textStyle, textColor } = this.props;
    const { viewSize, radius, diameter, getOffset } = this;
    const { progressColor } = this.state;
    const strokeDashoffset = getOffset(progress);
    const color = {
      red: Math.ceil(progressColor.red()),
      green: Math.ceil(progressColor.green()),
      blue: Math.ceil(progressColor.blue())
    };
    const transition = animate ? `stroke-dashoffset ${animationDuration}s linear, stroke ${animationDuration}s linear` : null;
    const strokeLinecap = roundedStroke ? 'round' : 'butt';
    const svgSize = responsive ? '100%' : size;
    const dashLength = diameter;
    return (
      <svg width={svgSize} height={svgSize}
           viewBox={`${-(viewSize / 2 - radius)} ${-(viewSize / 2 - radius)} ${viewSize} ${viewSize}`}>
        <circle stroke={bgColor} cx={radius} cy={radius} r={radius} strokeWidth={lineWidth * 0.9} fill="none"/>
        <circle
          transform={`rotate(-90 ${radius} ${radius})`}
          stroke={`rgb(${color.red},${color.green},${color.blue})`}
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
        />

        <circle
          stroke="none"
          cx={radius}
          cy={0}
          r={lineWidth * 1.5}
          fill={`rgb(${color.red},${color.green},${color.blue})`}
          transform={`rotate(${progress * 360} ${radius} ${radius})`}
          style={{
            transition: `transform ${animationDuration}s linear, fill ${animationDuration}s linear`
          }}
        />
        <text style={textStyle} fill={textColor} x={radius} y={radius} textAnchor="middle" dominantBaseline="central">
          {text}
        </text>
      </svg>
    );
  }
}

Circle.propTypes={
  text:PropTypes.string,
  size:PropTypes.number.isRequired,
  bgColor:PropTypes.string.isRequired,
  progress:PropTypes.number.isRequired,
  startColor:PropTypes.string.isRequired,
  endColor:PropTypes.string.isRequired,
  textColor:PropTypes.string.isRequired,
  animate:PropTypes.bool,
  animationDuration:PropTypes.number,
  roundedStroke:PropTypes.bool,
  responsive:PropTypes.bool,
  textStyle:PropTypes.shape({}),
  lineWidth:PropTypes.number,
};

Circle.defaultProps = {
  text:'',
  animate:true,
  animationDuration:0,
  roundedStroke:true,
  responsive:true,
  lineWidth:50,
  textStyle:{
    font: '30rem Helvetica, Arial, sans-serif'
  },
};
