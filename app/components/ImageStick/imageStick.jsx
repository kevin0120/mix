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

import popoverStyles from '../../common/jss/popoverStyles';

import {
  successColor,
  warningColor,
  dangerColor
} from '../../common/jss/material-react-pro';

import { keyframes } from 'react-emotion';
import classNames from 'classnames';
import { OPERATION_RESULT } from '../../reducers/operations';

const ripple = keyframes`
  0% {transform:scale(0.5); }
  75% {transform:scale(1.0); opacity:1;}
  100% {transform:scale(1.75); opacity:0;}
`;

const circleRadius = 30;

const withstyles = () => ({
  ...popoverStyles,
  picWrap: {
    position: 'relative',
    // marginTop: '10px',
    // marginLeft: '10px',
    height: '100%',
    // height: 'calc(100% - 50px)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden'
  },
  imgBlock: {
    maxHeight: '100%',
    maxWidth: '100%',
    // width:'100%',
    // height:'100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
    position: 'relative',
    backgroundSize: 'contain',
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center'
  },
  imgSheet: {
    position: 'relative',
    maxWidth: '100%',
    maxHeight: '100%',
    textAlign: 'center'
    // backgroundSize: 'contain',
    // backgroundPosition: 'center',
    // backgroundRepeat: 'no-repeat',
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
    lineHeight: '60px',
    fontSize: '40px',
    overflow: 'hidden',
    background: '#dbdbdb'
  },
  imgInfo: {
    margin: '20px',
    position: 'absolute',
    color: '#333'
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
  componentDidMount() {}

  componentDidUpdate() {}

  render() {
    const { classes, operations } = this.props;

    let idx = 0;

    const statusDisplay = operations.results.map((item, i) => {
      const display = operations.activeResultIndex >= idx;

      const postionStyle = {
        top: `calc(${item.offset_y}% - 30px)`,
        left: `calc(${item.offset_x}% - 30px)`
      };

      idx += 1;

      let status = 'waiting';
      if (operations.activeResultIndex === i) {
        status = 'waitingActive';
      }

      if (item.result === OPERATION_RESULT.OK) {
        status = 'success';
      } else if (item.result === OPERATION_RESULT.NOK) {
        status = 'error';
      }

      return (
        <div key={item.id} style={postionStyle} className={classes.imgInfo}>
          <span className={`${classes.circleStatus} ${classes[status]}`}>
            {item.sequence}
          </span>
          {display ? (
            <div className={classes.popover}>
              <div className={classes.popoverBody}>
                <p>角度: {item.wi || '-'}</p>
                <p>扭矩: {item.mi || '-'}</p>
                <p>时间: {item.ti || '-'}</p>
              </div>
            </div>
          ) : null}
        </div>
      );
    });
    return (
      <div elevation={4} className={classes.picWrap}>
        <div
          className={classes.imgBlock}
          style={{
            backgroundImage: `url(${operations.workSheet})`
          }}
        >
          <img
            src={operations.workSheet}
            className={classes.imgSheet}
            alt={''}
          />
          {statusDisplay}
        </div>
      </div>
    );
  }
}

ConnectedImageStick.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  operations: PropTypes.shape({}).isRequired
};

ConnectedImageStick.defaultProps = {};

export default withStyles(withstyles)(ConnectedImageStick);
