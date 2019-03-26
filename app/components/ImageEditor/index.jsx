import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Draggable from 'react-draggable';
import { withStyles } from '@material-ui/core/styles';
import { cloneDeep, get, isEqual } from 'lodash';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import Fade from '@material-ui/core/Fade';
import { GridLoader } from 'react-spinners';
import { css } from 'react-emotion';

import styles from './styles';
import Image from '../ImageStick/Image';
import DraggablePoint from './DraggablePoint';
import InsertPhotoIcon from "@material-ui/icons/InsertPhoto";
import QueueIcon from "@material-ui/icons/Queue";
import DeleteForeverIcon from "@material-ui/icons/DeleteForever";
import CheckCircleIcon from "@material-ui/icons/CheckCircle";

const override = css`
    display: block;
    margin: auto;
    border-color: red;
`;

const circleRadius = 30;

/* eslint-disable react/prefer-stateless-function */
class ImageEditor extends React.Component {
  // static getDerivedStateFromProps(nextProps, prevState) {
  //   if (!isEqual(nextProps.data, prevState.info) &&
  //     Object.keys(nextProps.data).length > 0) {
  //     return {
  //       ...prevState,
  //       info: {
  //         ...nextProps.data,
  //       },
  //     };
  //   }
  //   return null;
  // }

  static Transition(props) {
    return <Fade {...props} timeout={500}/>;
  }

  constructor(props) {
    super(props);
    this.state = {
      points: this.props.points,
      loading: false
    };
    this.handleStop = this.handleStop.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { points } = this.props;
    if (JSON.stringify(nextProps.points) !== JSON.stringify(points)) {
      this.setState({
        points: nextProps.points
      });
    }
  }

  handleStop(newPoint, idx) {
    const { points } = this.state;
    console.log(points);
    const newPoints = cloneDeep(points);
    newPoints.splice(idx, 1, newPoint);
    this.setState({
      points: newPoints
    });
  }

  addPoint() {
    const { points } = this.state;
    const newPoints = cloneDeep(points);
    newPoints.push({
      sequence: newPoints.length,
      x_offset: 0,
      y_offset: 0
    });
    this.setState({
      points: newPoints
    });
  }

  removePoint() {
    const { points } = this.state;
    const newPoints = cloneDeep(points);
    newPoints.pop();
    this.setState({
      points: newPoints
    });
  }

  saveChanges() {
    const { doEdit } = this.props;
    const { points, img } = this.state;
    doEdit({
      points,
      img
    });
  }


  render() {
    const { classes, img } = this.props;
    const { points } = this.state;


    return (
      <div className={classes.imageWrap}>
        <Dialog
          fullScreen
          classes={{
            root: classes.loadModal
          }}
          open={this.state.loading}
          style={{ opacity: 0.7 }}
          TransitionComponent={this.Transition}
        >
          <GridLoader
            className={override}
            sizeUnit={'px'}
            size={50}
            color={'#36D7B7'}
            loading={this.state.loading}
          />
        </Dialog>
        <Image src={img} alt="">
          {points.map((point, idx) => (
            <DraggablePoint
              point={point}
              idx={idx}
              radius={30}
              onStop={this.handleStop}
            />
          ))}
        </Image>
      </div>
    );
  }
}

ImageEditor
  .propTypes = {
  classes: PropTypes.shape({}).isRequired,
  doEdit: PropTypes.func.isRequired
};

export default withStyles(styles)(ImageEditor);
