import customSelectStyle from '../../common/jss/customSelectStyle';

export default theme => ({
  ...customSelectStyle,
  section: {
    overflowY: 'scroll',
    // flex: 1,
    height: '100%',

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
    marginRight: theme.spacing.unit
  },
  button: {
    margin: '10px 20px'
  },
  modalRoot: {
    overflow: 'auto',
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
    width: 138
  },
  input: {
    width: 300
  },
  // Io
  ioInputLabel: {
    width: 80
  },
  ioInput: {
    width: 250,
  },
  ioFunctionSelect: {
    margin: '0 70px',
    width: 150,
    fontSize: 14,
    height:32
  },
  // Connect
  sectionTitleInner: {
    marginTop: 20
  },
  // test
  testButton: {
    margin: 'auto 40px'
  }
});
