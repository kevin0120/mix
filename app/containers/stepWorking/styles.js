const layout = {
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'row',
    backdropFilter: 'blur(2px)'
  },
  rightContainer: {
    flex: 3,
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    width: '30%'
    // padding:3,
  },
  leftContainer: {
    flex: 7,
    margin: 3
  },
  timerContainer: {
    flex: 1,
    margin: 3,
    zIndex: 1,
    display:'flex',
    flexDirection:'row',
    alignItems:'center',
    justifyContent:'center',
  },
  stepperContainer: {
    margin: 3,
    flex: 4,
    overflowY:'scroll',
  },
  contentContainer: {
    flex: 1
  },
  orderInfoContainer: {
    height: 60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center'
  },
  buttonsContainer: {
    height:60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center'
  }
};

const stepper = {
  stepper: {},
  stepButton: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-start !important',
    height: 50,
    padding:0
  }
};

export default {
  layout,
  stepper
};

