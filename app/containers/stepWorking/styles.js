const layout = {
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'row',
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',
    backgroundColor: '#7c7c7c'

  },
  rightContainer: {
    flex: 2,
    display: 'flex',
    flexDirection: 'column',
    height: '100%',
    width: '30%',
    // padding:3,
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',

  },
  leftContainer: {
    flex: 7,
    margin: 3,
    display: 'flex',
    flexDirection: 'column',
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',
    backgroundColor: '#eee'
  },
  timerContainer: {
    flex: 1,
    margin: 3,
    zIndex: 1,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',
    backgroundColor: '#eee'
  },
  stepperContainer: {
    margin: 3,
    flex: 4,
    overflowY: 'auto',
    display:'flex',
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',
    backgroundColor: '#eee',
    padding:0
  },
  contentContainer: {
    flex: 1,
    backdropFilter: 'blur(2px)',

  },
  orderInfoContainer: {
    height: 60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center',
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',
    backgroundColor: '#eee',
  }
};

const stepperContainer = {
  root: {
    padding: 16,
    backgroundColor:'transparent'
  },
  stepButton: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-start !important',
    height: 50,
    padding: '0 0',
    margin: '10px 0'
  }
};

const buttonsContainer = {
  root: {
    height: 60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center',
  }
};

const stepPageContainer = {
  root: {
    // display: 'flex',
    // flexDirection: 'row',
    flex: 1,
    '& *':{
      backgroundColor:'transparent',
    }
  },

  left: {
    flex: 3,
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
  },
  right: {
    flex: 1,
    height: '100%',
  },
  image: {
    flex: 1
  },
  description: {
    flex: 1
  },
  Paper:{
    width:'100%',
    height:'100%',
    padding:6
  },
  result: {
    flex: 2
  }
};

export default {
  layout,
  stepperContainer,
  buttonsContainer,
  stepPageContainer
};

