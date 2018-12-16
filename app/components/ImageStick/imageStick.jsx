import React from 'react';
import PropTypes from 'prop-types';
import Paper from '@material-ui/core/Paper';
import { withStyles } from '@material-ui/core/styles';

// import {
//   loadWorkOrderLogo,
//   fetchWorkOrder
// } from '../../actions/ongoingWorkOrder';

// import {
//   fetchRoutingWorkcenterbyCarType,
//   fetchRoutingWorkcenterbyJobId
// } from '../../actions/ongoingRoutingWorkcenter';

import { keyframes } from 'react-emotion';
import Image from './Image';
import popoverStyles from '../../common/jss/popoverStyles';

import {
  successColor,
  warningColor,
  dangerColor
} from '../../common/jss/material-react-pro';

import { OPERATION_STATUS, OPERATION_RESULT } from '../../reducers/operations';

import Card from "../Card/Card"
import imagesStyles from "../../common/jss/imagesStyles";

const ripple = keyframes`
  0% {transform:scale(0.5); }
  75% {transform:scale(1.0); opacity:1;}
  100% {transform:scale(1.75); opacity:0;}
`;

const circleRadius = 30;
const scaleRate = 2;

const imgStickStyles = () => ({
  ...popoverStyles,
  ...imagesStyles,
  picWrap: {
    position: 'relative',
    // marginTop: '10px',
    // marginLeft: '10px',
    height: '100%',
    width: '100%',
    // height: 'calc(100% - 50px)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
    padding: 'auto'
  },
  heightFirst: {
    height: '100%'
  },
  widthFirst: {
    width: '100%'
  },
  circleStatus: {
    display: 'inline-block',
    width: `${circleRadius * 2}px`,
    height: `${circleRadius * 2}px`,
    borderRadius: '99%',
    textAlign: 'center',
    lineHeight: `${circleRadius * 2}px`,
    fontSize: `${(circleRadius-10) * 2}px`,
    overflow: 'hidden',
    background: '#dbdbdb'
  },
  circleSmallStatus: {
    display: 'inline-block',
    width: `${circleRadius}px`,
    height: `${circleRadius}px`,
    borderRadius: '99%',
    textAlign: 'center',
    lineHeight: `${circleRadius}px`,
    fontSize: `${(circleRadius-10)}px`,
    overflow: 'hidden',
    background: '#dbdbdb'
  },
  imgInfo: {
    margin: '20px',
    position: 'absolute',
    color: '#333'
  },
  imgSmallBlock: {
    position: 'absolute',
    width:'30%',
    height:'30%',
    bottom:0,
    left:0,
    marginBottom: '5px'
  },
  imgStausInfo: {
    padding: '5px 10px',
    borderRadius: '4px',
    background: '#fff',
    boxShadow: '1px 1px 2px #dbdbdb inset',
    fontSize: '18px',
    marginTop: '0',

    '&>p': {
      lineHeight: 2
    }
  },
  waiting: {
    background: warningColor
  },
  waitingActive: {
    background: warningColor,
    animation: `${ripple} 1.0s infinite cubic-bezier(1, 1, 1, 1)`
  },
  success: {
    background: successColor
  },
  error: {
    background: dangerColor
  }
});

/* eslint-disable react/prefer-stateless-function */
class ConnectedImageStick extends React.Component {
  constructor(props) {
    super(props);
    this.imageTransform = 'none';
  }

  focused = false;

  doFocus({ transform, scale }) {
    this.imageTransform = `translate(calc(${transform.x || 0}% - ${transform.x && circleRadius * scaleRate
      }px),calc(${transform.y || 0}% - ${transform.y && circleRadius * scaleRate}px)) scale(${scale.x},${scale.y})`;
  }

  render() {
    const { classes, operations } = this.props;

    // let idx = 0;

    if (operations.results.length === 0) {
      this.doFocus({
        transform: {
          x: 0,
          y: 0
        },
        scale: {
          x: 1,
          y: 1
        }
      });
      this.focused=false;
    } else {
      this.doFocus({
        transform: {
          x: (50 - operations.results[operations.activeResultIndex].offset_x) * 2,
          y: (50 - operations.results[operations.activeResultIndex].offset_y) * 2
        },
        scale: {
          x: 2,
          y: 2
        }
      });
      this.focused=true;
    }
    const statusDisplay = (small=false) => operations.results.map((item, i) => {
      // const display = operations.activeResultIndex >= idx;

      const postionStyle = {
        top: `calc(${item.offset_y}% - ${this.focused ? circleRadius * scaleRate : circleRadius}px)`,
        left: `calc(${item.offset_x}% - ${this.focused ? circleRadius * scaleRate : circleRadius}px)`
      };

      const circleStatus = small? classes.circleStatus: classes.circleStatus;

      idx += 1;

      let status = 'waiting';
      if (operations.results[operations.activeResultIndex].group_sequence === item.group_sequence && operations.operationStatus === OPERATION_STATUS.DOING) {
        status = 'waitingActive';
      }

      if (item.result === OPERATION_RESULT.OK) {
        status = 'success';
      } else if (item.result === OPERATION_RESULT.NOK) {
        status = 'error';
      }

      return (
        <div key={item.id} style={postionStyle} className={classes.imgInfo}>
          <span className={`${circleStatus} ${classes[status]}`}>
            {item.sequence}
          </span>
          {/* {display ? ( */}
          {/* <div className={classes.popover}> */}
          {/* <div className={classes.popoverBody}> */}
          {/* <p>角度: {item.wi || '-'}</p> */}
          {/* <p>扭矩: {item.mi || '-'}</p> */}
          {/* <p>时间: {item.ti || '-'}</p> */}
          {/* </div> */}
          {/* </div> */}
          {/* ) : null} */}
        </div>
      );
    });



    return (
      <div elevation={4} className={classes.picWrap}>
        <Image
          src={operations.workSheet}
          alt=""
          style={{
            transform: this.imageTransform,
            transition: 'transform 1s'
          }}
        >
          {statusDisplay(false)}
        </Image>
        <Card plain raised className={classes.imgSmallBlock}>
          <Image className={classes.imgCard}
            src={operations.workSheet}
            alt=""
          >
            {statusDisplay(true)}
          </Image>
        </Card>
      </div>
    );
  }
}

ConnectedImageStick.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  operations: PropTypes.shape({}).isRequired
};

ConnectedImageStick.defaultProps = {};

export default withStyles(imgStickStyles)(ConnectedImageStick);
