import React from 'react';
import PropTypes from 'prop-types';
import Paper from '@material-ui/core/Paper';
import { withStyles } from '@material-ui/core/styles';
import Fade from '@material-ui/core/Fade';

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

import Card from '../Card/Card';
import imagesStyles from '../../common/jss/imagesStyles';

const ripple = keyframes`
  0% {transform:scale(0.5); }
  75% {transform:scale(1.0); opacity:1;}
  100% {transform:scale(1.75); opacity:0;}
`;

const lodash = require('lodash');

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
    display: 'block',
    width: `${circleRadius * 2}px`,
    height: `${circleRadius * 2}px`,
    borderRadius: '50%',
    textAlign: 'center',
    lineHeight: `${circleRadius * 2}px`,
    fontSize: `${(circleRadius - 10) * 2}px`,
    overflow: 'hidden',
    background: '#dbdbdb'
  },
  circleSmallStatus: {
    display: 'block',
    width: `${circleRadius}px`,
    height: `${circleRadius}px`,
    borderRadius: '50%',
    textAlign: 'center',
    lineHeight: `${circleRadius}px`,
    fontSize: `${(circleRadius - 10)}px`,
    overflow: 'hidden',
    background: '#dbdbdb'
  },
  imgInfo: {
    margin: '0',
    position: 'absolute',
    color: '#333'
  },
  imgSmallBlock: {
    position: 'absolute',
    width: '30%',
    height: '30%',
    bottom: 0,
    left: 0,
    marginBottom: '5px'
  },
  imgBlock: {
    ...imagesStyles.imgCard,
    width: '100%',
    height: '100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
    position: 'relative'
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
  failActive: {
    background: dangerColor,
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
    this.focus = {
      do: false,
      transform: {
        x: 0,
        y: 0
      },
      scale: 1
    };

  }

  componentWillReceiveProps(nextProps) {
    const { operations, enableFocus } = this.props;
    const { operations: nextOperations } = nextProps;
    if (!enableFocus) {
      this.focus.do = false;
      return;
    }
    if (nextOperations.results.length === 0 || nextOperations.operationStatus === OPERATION_STATUS.READY) {
      // 当接受到的结果为空,没有拧紧点, 或者工单进入ready阶段(代表着上一个作业结束)
      this.doFocus({
        do: false,
        transform: {
          x: 0,
          y: 0
        },
        scale: 1
      });

    } else if (nextOperations.operationStatus === OPERATION_STATUS.PREDOING) {
      // do nothing
    } else {
      this.doFocus({
        do: true,
        transform: {
          x: (50 - nextOperations.results[nextProps.operations.activeResultIndex].offset_x) * 2,
          y: (50 - nextOperations.results[nextProps.operations.activeResultIndex].offset_y) * 2
        },
        scale: 2
      });
    }
  }


  doFocus(focus) {
    this.focus = focus;
    // this.imageTransform = `translate(${transform.x || 0}%,${transform.y || 0}%) scale(${scale},${scale})`;
  }

  statusDisplay = (small = false) => {
    const { classes, operations } = this.props;
    // let idx = 0;
    return operations.results.map((item, i) => {
      // const display = operations.activeResultIndex >= idx;
      const cR = small ? circleRadius/2 : circleRadius ;
      const postionStyle = {
        top: `calc(${item.offset_y}% - ${cR}px)`,
        left: `calc(${item.offset_x}% - ${cR}px)`
      };
      const circleStatus = small ? classes.circleSmallStatus : classes.circleStatus;
      // idx += 1;
      let status = 'waiting';
      if (item.result === OPERATION_RESULT.NOK) {
        status = 'error';
      } else if(item.result === OPERATION_RESULT.OK){
        status = 'success';
      }
      if (operations.results[operations.activeResultIndex].group_sequence === item.group_sequence) {
        switch (operations.operationStatus) {
          // case OPERATION_STATUS.TIMEOUT:
          case OPERATION_STATUS.DOING: {
            status = 'waitingActive';
            if (item.result === OPERATION_RESULT.NOK) {
              status = 'failActive';
            }
            break;
          }
          case OPERATION_STATUS.FAIL:
            status = 'error';
            break;
          default:
            break;
        }
      }
      //   if (operations.operationStatus === OPERATION_STATUS.DOING) {
      //     status = 'waitingActive';
      //     if  (item.result === OPERATION_RESULT.NOK) {
      //       status = 'failActive';
      //     }else {
      //       status = 'success';
      //     }
      //   } else if(operations.operationStatus === OPERATION_STATUS.FAIL) {
      //     status = 'error';
      //   }
      // }

      // if (item.result === OPERATION_RESULT.OK) {
      //   status = 'success';
      // } else if (item.result === OPERATION_RESULT.NOK) {
      //   status = 'failActive';
      // }

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
  };

  render() {
    const { classes, operations, enableFocus } = this.props;

    const smallImgDisplay = operations.operationStatus !== 'Ready' && operations.operationStatus !== 'PreDoing';


    return (
      <div className={classes.picWrap}>
        <Image
          src={operations.workSheet}
          alt=""
          focus={this.focus}
        >
          {this.statusDisplay(false)}
        </Image>
        {
          enableFocus ? <Fade in={smallImgDisplay}
                              {...(smallImgDisplay ? { timeout: 1000 } : {})}>
            <Card plain raised className={classes.imgSmallBlock}>
              <Image
                src={operations.workSheet}
                alt=""
              >
                {this.statusDisplay(true)}
              </Image>
            </Card>
          </Fade> : null
        }
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
