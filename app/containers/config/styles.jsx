import customSelectStyle from '../../common/jss/customSelectStyle';

const config = theme => ({
  root: {
    // flexGrow: 1,
    zIndex: 1,
    overflow: 'hidden',
    position: 'relative',
    display: 'flex',
    color: '#333',
    height: '100%',
    width: '100%'
  },
  content: {
    flex: 1,
    // backgroundColor: theme.palette.background.default,
    background: '#f2f2f2',
    padding: theme.spacing(2),
    minWidth: 0, // So the Typography noWrap works
    height: '100%',
    paddingTop: '64px'
  },
  toolbar: theme.mixins.toolbar,
  menuItem: {
    fontSize: 14,
    height: 120,
    margin: '15px 5px'
  },
  menuItemSelected: {
    fontSize: 14,
    backgroundColor: theme.palette.warning.main,
    height: 120,
    color: '#FAFAFA',
    margin: '15px 5px'
  },
  itemText: {
    fontSize: 20,
    padding: 0
  },
  cardActionArea: {
    flex: 1,
    color: '#00352c',
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center'
  },
  cardActionAreaSelected: {
    flex: 1,
    color: 'white',
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center'
  }
});

const content = theme => ({
  ...customSelectStyle,
  section: {
    overflowY: 'auto',
    // flex: 1,
    height: '100%'
  },
  paperWrap: {
    fontSize: 14
  },
  content: {
    height: '100%'
  },
  sectionTitle: {
    fontSize: 14,
    marginTop: 0,
    marginBottom: 10,
    marginLeft: 20
  },
  leftIcon: {
    marginRight: theme.spacing(1)
  },
  button: {
    margin: '10px 20px'
  },
  modalRoot: {
    overflowY: 'auto',
    display: 'block'
  },
  inputItem: {
    padding: 20,
    '& label': {
      fontSize: 14
    },
    '& input': {
      fontSize: 14
    }
  },
  statusWrap: {
    '&>span': {
      verticalAlign: 'middle'
    },
    '&>span:first-child': {
      marginRight: 5
    }
  },
  statusCircle: {
    minWidth: 15,
    height: 15,
    borderRadius: '50%',
    display: 'inline-block'
  },
  success: {
    background: '#3cb87e'
  },
  successText: {
    color: '#3cb87e'
  },
  fail: {
    background: '#f04848'
  },
  failText: {
    color: '#f04848'
  },
  info: {
    background: '#bdbdbd'
  },
  infoText: {
    color: '#bdbdbd'
  },
  // Net
  inputLabel: {
    width: 160
  },
  input: {
    width: 300
  },
  // Io
  ioInputLabel: {
    width: 80
  },
  ioInput: {
    width: 250
  },
  ioFunctionSelect: {
    margin: '0 70px',
    width: 150,
    fontSize: 14,
    height: 32
  },
  // Connect
  sectionTitleInner: {
    marginTop: 20
  },
  // test
  testButton: {
    // padding: '18px 30px',
    margin: 'auto 40px'
  }
});
export default {
  config,
  content
};
