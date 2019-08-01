import { keyframes } from '@emotion/core';

const TOPHEIGHT = '150px';

const twinkling = keyframes`
  0% {
    opacity: .1;
  }
  100% {
    opacity: 1;
  }
`;

const ripple = keyframes`
  0% {transform:scale(0.75); }
  75% {transform:scale(1); opacity:1;}
  100% {transform:scale(1.25); opacity:0;}
`;

export default (theme) => ({
  sWorkingInitializing: {
    backgroundColor: theme.palette.gray.main
    // animation: `${ripple} 0.5s infinite cubic-bezier(1, 1, 1, 1)`,
  },
  sWorkingDoing: {
    backgroundColor: theme.palette.success.main
  },
  sWorkingReady: {
    backgroundColor: theme.palette.success.main,
    animation: `${twinkling} 1s infinite cubic-bezier(1, 1, 1, 1)`
  },
  sWorkingFail: {
    backgroundColor: theme.palette.warning.main
  },
  sWorkingError: {
    backgroundColor: theme.palette.danger.main
  },
  sWorkingTimeout: {
    backgroundColor: theme.palette.warning.main
  },
  sWorkingContinue: {
    backgroundColor: theme.palette.success.main
  },
  statusContainer: {
    padding: '10px 10px 10px 0'
  },
  progress: {
    boxSizing: 'border-box',
    height: '100%',
    backgroundColor: '#fff',
    border: '4px solid #00aaf2'
  },
  // progressRoot: {
  //   background: 'primary',//'linear-gradient(to right, #00E676 30%, #FFEB3B, #FF5722 90%)',
  //   borderRadius: 3,
  //   border: 0,
  //   color: 'white',
  //   height: '100%',
  //   padding: '0 30px',
  //   boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)',
  // },
  progressMote: {
    borderRadius: 3,
    border: 0,
    height: '100%'
    // boxShadow: '0 3px 5px 2px rgba(255, 105, 135, .3)',
  },
  paperStatus: {
    padding: '0',
    height: '100%',
    textAlign: 'center',
    // TOPHEIGHT - 上下padding 20px
    lineHeight: '130px',
    boxSizing: 'border-box',
    fontSize: '44px',
    fontWeight: 'bold',
    color: '#fff'
  },
  progressWrap: {
    height: '100%',
    position: 'relative'
  },
  progressText: {
    position: 'absolute',
    top: '0',
    left: '0',
    width: '98%',
    boxSizing: 'border-box',
    overflow: 'hidden',
    textAlign: 'center',
    height: TOPHEIGHT,
    lineHeight: TOPHEIGHT,
    fontSize: '60px',
    color: '#333',
    background: 'transparent',
    boxShadow: 'none',
    fontWeight: 'bold',
    whiteSpace: 'nowrap'
  }
});
