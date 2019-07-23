const layout = {
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'row',
    backdropFilter: 'blur(2px)'
  },
  rightContainer: {
    flex: 2,
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    width: '30%'
    // padding:3,
  },
  leftContainer: {
    flex: 7,
    margin: 3,
    display:'flex',
    flexDirection:'column',
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
    flex: 1,
  },
  orderInfoContainer: {
    height: 60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center'
  }
};

const stepperContainer = {
  root: {
    padding:16
  },
  stepButton: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-start !important',
    height: 50,
    padding:'0 0',
    margin:'10px 0'
  }
};

const buttonsContainer={
  root: {
    height:60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center'
  }
};

const stepPageContainer={
  root:{
    display:'flex',
    flexDirection:'row',
    flex:1
  },
  left:{
    flex:3,
    height:'100%',
    display:'flex',
    flexDirection:'column',
  },
  right:{
    flex:1,
    height:'100%',
    display:'flex',
    flexDirection:'column',
  },
  image:{
    flex:1
  },
  description:{
    flex:1
  },
  result:{
    flex:2
  },
};

export default {
  layout,
  stepperContainer,
  buttonsContainer,
  stepPageContainer
};

