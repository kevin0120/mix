export default {
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'row',
    backgroundColor: 'rgb(74,84,89)',
    backdropFilter: 'blur(2px)',

  },
  background: {},
  rightContainer: {
    flex: 3,
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    width: '30%',
    // padding:3,
    backgroundColor: 'rgb(74,84,89)'

  },

  leftContainer: {
    flex: 7,
    margin: 3,
    backgroundColor: 'rgb(74,84,89)'

  },
  timerContainer: {
    flex: 1,
    margin: 3,
    zIndex: 1
  },
  stepperContainer: {
    margin: 3,
    flex: 4

  },
  stepButton: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-start',
    height: 50
  },
  stepper: {},
  contentContainer: {
    flex: 1
  },
  orderInfoContainer: {
    height: 60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between'
  },
  buttonsContainer: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    background: 'hsla(0,0%,100%,.3)',
    backdropFilter: 'blur(2px)',
  }
};
