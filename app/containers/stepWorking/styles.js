const layout = (theme) => ({
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'row',
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',
    backgroundColor: theme.palette.grey[100]

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
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: theme.palette.common.white
  },
  timerContainer: {
    flex: 1,
    margin: 3,
    zIndex: 1,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: theme.palette.common.white
  },
  stepperContainer: {
    margin: 3,
    flex: 4,
    overflowY: 'auto',

    backgroundColor: theme.palette.common.white,
    padding: 0
  },
  contentContainer: {
    flex: 1
  },
  orderInfoContainer: {
    height: 60,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: `0 ${theme.spacing(2)}px`,
    alignItems: 'center',

    backgroundColor: theme.palette.common.white
  }
});

const stepperContainer = {
  root: {
    padding: 16,
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
    height: 70,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center'
  },
  menuIcon: {

  },
  menuButton: {},
  dialog: {},
  dialogContainer: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center'
  },
  bigButton: {
    width: '20vw',
    height: '15vh'
  }
};

const stepPageContainer = {
  root: {
    flex: 1,
  },

  left: {
    flex: 3,
    height: '100%',
    display: 'flex',
    flexDirection: 'column'
  },
  right: {
    flex: 1,
    height: '100%'
  },
  image: {
    flex: 1,
    overflow: 'hidden'
  },
  description: {
    flex: 1
  },
  Paper: {
    width: '100%',
    height: '100%',
    padding: 6
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

