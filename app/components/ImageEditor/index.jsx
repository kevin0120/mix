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
import { css } from '@emotion/core';


// icons
import InsertPhotoIcon from '@material-ui/icons/InsertPhoto';
import QueueIcon from '@material-ui/icons/Queue';
import DeleteForeverIcon from '@material-ui/icons/DeleteForever';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import CancelIcon from '@material-ui/icons/Cancel';
import EditIcon from '@material-ui/icons/Edit';
import DraggablePoint from './DraggablePoint';
import Image from '../ImageStick/Image';
import styles from './styles';

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
      editing: false,
      img: this.props.img
    };
  }

  componentWillReceiveProps(nextProps) {
    const { points, img, loading } = this.props;
    if (
      JSON.stringify(nextProps.points) !== JSON.stringify(points) ||
      nextProps.img !== img ||
      nextProps.loading !== loading
      && nextProps.loading === false
    ) {
      this.setState({
        points: nextProps.points,
        editing: false,
        img: nextProps.img
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
    const { points, img } = this.props;
    if (editing) {
      this.setState({
        points,
        img
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
    const { doEdit } = this.props;
    const { points, img, editing } = this.state;
    doEdit({
      points,
      img
    });
    // this.switchEdit();

    this.setState({
      editing: !editing
    });
  };

  setImage = (img) => {
    this.setState({
      img
    });
  };

  imgChange = (e) => {
    const fileObj = e.target.files[0];
    if (fileObj && /^image/.test(fileObj.type)) {
      const reader = new FileReader();
      reader.readAsDataURL(fileObj);
      reader.onload = () => {
        this.setImage(reader.result);
      };
    } else {
      console.log('ERR! 选择的不是图片');
    }
  };

  buttonsIdle = [
    {
      icon: <EditIcon style={{ fill: '#03a9f4' }}/>,
      text: '编辑',
      onClick: this.switchEdit
    }
  ];

  buttonsEditing = [
    {
      icon: <CancelIcon style={{ fill: '#f44336' }}/>,
      text: '取消',
      onClick: this.switchEdit
    },
    {
      text: '选择图片',
      icon: <InsertPhotoIcon style={{ fill: '#009688' }}/>,
      onClick: this.selectImage
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
      icon: <CheckCircleIcon style={{ fill: '#4caf50' }}/>,
      text: '保存',
      onClick: this.saveChanges
    }
  ];

  renderButton = (button, key) => {
    if (button.text === '选择图片') {
      return (
        <React.Fragment key={key}>
          <input
            accept="image/*"
            id="contained-button-file"
            type="file"
            style={{ display: 'none' }}
            onChange={this.imgChange}
          />
          <label htmlFor="contained-button-file">
            <Button
              key={key}
              // variant="contained"
              component="span"
            >
              {button.icon}
              {button.text}
            </Button>
          </label>
        </React.Fragment>
      );
    }
    return (
      <Button
        key={key}
        onClick={() => button.onClick()}
      >
        {button.icon}
        {button.text}
      </Button>
    );
  };


  render() {
    const { classes, loading } = this.props;
    const { points, editing, img } = this.state;


    return (
      <div className={classes.imageWrap}>
        <Toolbar>
          {editing ? this.buttonsEditing.map((button, id) => (
            this.renderButton(button, id)
          )) : this.buttonsIdle.map((button, id) => (
            this.renderButton(button, id)
          ))}
        </Toolbar>
        <div style={{ width: '100%', height: 'calc(100% - 64px)' }}>
          <Image src={img} alt="">
            {points.map((point, idx) => (
              <DraggablePoint
                key={point.sequence}
                point={point}
                idx={idx}
                radius={circleRadius}
                onStop={this.pointDraggingStop}
                disabled={!editing}
              />
            ))}
          </Image>
        </div>
        <Dialog fullScreen
                classes={{
                  root: classes.loadModal
                }}
                open={loading}
                style={{ opacity: 0.7 }}
                TransitionComponent={this.Transition}
        >
          <GridLoader
            className={override}
            sizeUnit="px"
            size={50}
            color="#36D7B7"
            loading={loading}
          />
        </Dialog>
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
