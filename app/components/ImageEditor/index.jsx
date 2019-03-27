import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Draggable from 'react-draggable';
import { withStyles } from '@material-ui/core/styles';
import { cloneDeep, get, isEqual } from 'lodash';

import Toolbar from '@material-ui/core/Toolbar';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import Fade from '@material-ui/core/Fade';
import { GridLoader } from 'react-spinners';
import { css } from 'react-emotion';

import styles from './styles';
import Image from '../ImageStick/Image';
import DraggablePoint from './DraggablePoint';
import InsertPhotoIcon from '@material-ui/icons/InsertPhoto';
import QueueIcon from '@material-ui/icons/Queue';
import DeleteForeverIcon from '@material-ui/icons/DeleteForever';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';

const override = css`
    display: block;
    margin: auto;
    border-color: red;
`;

const circleRadius = 30;

/* eslint-disable react/prefer-stateless-function */
class ImageEditor extends React.Component {
  static Transition(props) {
    return <Fade {...props} timeout={500}/>;
  }

  constructor(props) {
    super(props);
    this.state = {
      points: this.props.points,
      editing: false
    };
  }

  componentWillReceiveProps(nextProps) {
    const { points } = this.props;
    if (JSON.stringify(nextProps.points) !== JSON.stringify(points)) {
      this.setState({
        points: nextProps.points,
        editing: false
      });
    }
  }

  pointDraggingStop = (newPoint, idx) => {
    const { points } = this.state;
    const newPoints = cloneDeep(points);
    newPoints.splice(idx, 1, newPoint);
    this.setState({
      points: newPoints
    });
  };

  // 编辑/取消
  switchEdit = () => {
    const { editing } = this.state;
    const { points } = this.props;
    if (editing) {
      this.setState({
        points
      });
    }
    this.setState({
      editing: !editing
    });
  };

  // 添加点位
  addPoint = () => {
    const { points } = this.state;
    const newPoints = cloneDeep(points);
    newPoints.push({
      sequence: newPoints.length + 1,
      x_offset: 0,
      y_offset: 0
    });
    this.setState({
      points: newPoints
    });
  };

  // 移除点位
  removePoint = () => {
    const { points } = this.state;
    const newPoints = cloneDeep(points);
    newPoints.pop();
    this.setState({
      points: newPoints
    });
  };

  // 保存更改
  saveChanges = () => {
    const { doEdit, img } = this.props;
    const { points } = this.state;
    doEdit({
      points,
      img
    });
    this.switchEdit();
  };

  buttonsIdle = [
    {
      icon: <InsertPhotoIcon style={{ fill: '#009688' }}/>,
      text: '编辑',
      onClick: this.switchEdit
    }
  ];

  buttonsEditing = [
    {
      icon: null,
      text: '取消',
      onClick: this.switchEdit
    },
    {
      icon: <QueueIcon style={{ fill: '#ff9800' }}/>,
      text: '增加点位',
      onClick: this.addPoint
    },
    {
      icon: <DeleteForeverIcon style={{ fill: '#1e88e5' }}/>,
      text: '删除点位',
      onClick: this.removePoint
    },
    {
      icon: <CheckCircleIcon style={{ fill: '#f44336' }}/>,
      text: '保存',
      onClick: this.saveChanges
    }
  ];

  renderButton = (button,key) => (
    <Button key={key} onClick={() => button.onClick()}>
      {button.icon}
      {button.text}
    </Button>
  );


  render() {
    const { classes, img } = this.props;
    const { points, editing } = this.state;


    return (
      <div className={classes.imageWrap}>
        <Toolbar>
          {editing ? this.buttonsEditing.map((button,id) => (
            this.renderButton(button,id)
          )) : this.buttonsIdle.map((button,id) => (
            this.renderButton(button,id)
          ))}
        </Toolbar>
        <div style={{ width: '100%', height: 'calc(100% - 64px)' }}>
          <Image src={img} alt="">
            {points.map((point, idx) => (
              <DraggablePoint
                key={idx}
                point={point}
                idx={idx}
                radius={circleRadius}
                onStop={this.pointDraggingStop}
                disabled={!editing}
              />
            ))}
          </Image>
        </div>
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
