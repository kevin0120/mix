const layout = (theme) => ({
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    backgroundImage: 'url("../resources/imgs/texture.png")',
    backgroundRepeat: 'repeat',
    backgroundColor: theme.palette.grey[100]

  },
  main: {
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
    justifyContent: 'flex-start',
    textAlign: 'center',
    padding: `0 ${theme.spacing(2)}px`,
    alignItems: 'center',
    backgroundColor: theme.palette.common.white
  },
  orderStatus: {
    marginRight: theme.spacing(1)
  },
  statusTodo: {
    color: theme.palette.gray.main
  },
  statusWIP: {
    color: theme.palette.primary.main
  },
  statusDone: {
    color: theme.palette.success.main
  },
  statusCancel: {
    color: theme.palette.danger.main
  },
  statusPending: {
    color: theme.palette.warning.main
  }
});

const stepperContainer = theme => ({
  root: {
    padding: 16
  },
  stepButton: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'flex-start !important',
    height: 50,
    padding: '0 0',
    margin: '10px 0'
  },
  stepIconDoing: {
    animation: '$doing-icon-rotation 1s infinite linear',
    color: theme.palette.primary.main
  },
  fail: {
    color: theme.palette.error.main
  },
  '@keyframes doing-icon-rotation': {
    '0%': {
      transform: 'rotate(0deg)'
    },
    '100%': {
      transform: 'rotate(-360deg)'
    }
  },
  stepLabelRoot: {
    flex: 1,
    display: 'flex'
  },
  stepLabelContainer: {
    flex: 1,
    display: 'flex'
  },
  stepLabel: {
    flex: 1,
    display: 'flex',
    alignItems: 'center',
    flexDirection: 'row',
    justifyContent: 'space-between'
  },
  infoRow: {
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center'
  },
});

const buttonsContainer = {
  root: {
    height: 70,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: '0 10px',
    alignItems: 'center'
  },
  menuIcon: {},
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
    overflow:'hidden'
  },

  left: {
    flex: 3,
    display: 'flex',
    flexDirection: 'column',
    height:'100%',
    overflow:'hidden'
  },
  right: {
    flex: 1,
    height: '100%'
  },
  image: {
    flex: 1,
    overflow: 'hidden',
  },
  descriptionContainer: {
    flex: 1
  },
  Paper: {
    width: '100%',
    height: '100%',
    padding: 'auto 6px'
  },
  Description: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
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

