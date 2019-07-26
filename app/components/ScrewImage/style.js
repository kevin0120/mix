import { dangerColor, successColor, warningColor } from '../../common/jss/material-react-pro';

const circleRadius = 30;

export default {
  image: {
    container: {
      width: '100%',
      height: '100%',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      overflow: 'hidden',
      position: 'relative'
    },
    image: {
      maxHeight: '100%',
      maxWidth: '100%',
      textAlign: 'center',
      objectFit: 'contain',
      resize: 'both',
      flex: 1
    }
  },
  point: {
    root: {
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
    active: {
      animation: '$activeRipple 1s infinite cubic-bezier(1, 1, 1, 1)'
    },
    '@keyframes activeRipple': {
      '0%': {
        transform: 'scale(0.5)'
      },
      '75%': {
        transform: 'scale(1.0)',
        opacity: 1
      },
      '100%': {
        transform: 'scale(1.75)',
        opacity: 0
      }
    },
    waiting: {
      background: warningColor
    },
    success: {
      background: successColor
    },
    error: {
      background: dangerColor
    }
  }
};
