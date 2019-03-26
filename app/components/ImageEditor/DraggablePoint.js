import React from 'react';
import { withStyles } from '@material-ui/core';
import Draggable from 'react-draggable';
import styles from './styles';

class DraggablePoint extends React.Component {
  draggableStop = (event, data) => {
    const { point, idx, onStop, imageSize, radius } = this.props;
    const newPoint = {
      ...point,
      y_offset: (data.lastY + radius) / imageSize.height * 100,
      x_offset: (data.lastX + radius) / imageSize.width * 100
    };
    onStop(newPoint, idx);
  };

  render() {
    const { imageSize, point, classes, radius } = this.props;
    const diameter = radius * 2;
    const position = {
      y: imageSize.height * point.y_offset / 100 - radius,
      x: imageSize.width * point.x_offset / 100 - radius
    };

    return (
      <Draggable
        defaultClassName={classes.circleWrap}
        key={point.sequence}
        axis="both"
        bounds="parent"
        handle=".dragCircle"
        position={position}
        onStop={this.draggableStop}
      >
        <div
          className={'dragCircle'}
          style={{
            width: diameter,
            height: diameter,
            lineHeight: `${diameter}px`,
            textAlign: 'center',
            borderRadius: '50%',
            background: 'red',
            color: '#fff',
            cursor: 'pointer',
            position: 'absolute'
          }}
        >
          {point.sequence}
        </div>
      </Draggable>
    );
  }
}

export default withStyles(styles)(DraggablePoint);
