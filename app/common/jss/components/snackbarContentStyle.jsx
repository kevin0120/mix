// ##############################
// // // SnackbarContent styles
// #############################

const snackbarContentStyle = theme=>({
  root: {
    fontFamily:theme.typography.fontFamily,
    fontWeight:theme.typography.fontWeightRegular,
    flexWrap: 'unset',
    position: 'relative',
    padding: '20px 15px',
    lineHeight: '20px',
    marginBottom: '20px',
    fontSize: '14px',
    backgroundColor: 'white',
    color: '#555555',
    borderRadius: '3px',
    boxShadow:
      '0 12px 20px -10px rgba(255, 255, 255, 0.28), 0 4px 20px 0px rgba(0, 0, 0, 0.12), 0 7px 8px -5px rgba(255, 255, 255, 0.2)'
  },
  top20: {
    top: '20px'
  },
  top40: {
    top: '40px'
  },
  info: {
    backgroundColor: '#00d3ee',
    color: '#ffffff',
    ...theme.boxShadow.info
  },
  success: {
    backgroundColor: '#5cb860',
    color: '#ffffff',
    ...theme.boxShadow.success
  },
  warning: {
    backgroundColor: '#ffa21a',
    color: '#ffffff',
    ...theme.boxShadow.warning
  },
  maintenance: {
    backgroundColor: '#ffa21a',
    color: '#ffffff',
    ...theme.boxShadow.warning
  },
  danger: {
    backgroundColor: '#f55a4e',
    color: '#ffffff',
    ...theme.boxShadow.danger
  },
  error: {
    backgroundColor: '#ffa21a',
    color: '#ffffff',
    ...theme.boxShadow.danger
  },
  primary: {
    backgroundColor: '#af2cc5',
    color: '#ffffff',
    ...theme.boxShadow.primary
  },
  rose: {
    backgroundColor: '#eb3573',
    color: '#ffffff',
    ...theme.boxShadow.rose
  },
  message: {
    padding: '0',
    display: 'block',
    maxWidth: '89%'
  },
  close: {
    width: '11px',
    height: '11px'
  },
  iconButton: {
    width: '24px',
    height: '24px',
    padding: '0'
  },
  icon: {
    width: '38px',
    height: '38px',
    display: 'block',
    left: '15px',
    position: 'absolute',
    marginTop: '-39px',
    fontSize: '20px',
    backgroundColor: '#FFFFFF',
    padding: '6px',
    borderRadius: '50%',
    maxWidth: '38px',
    boxShadow:
      '0 10px 30px -12px rgba(0, 0, 0, 0.42), 0 4px 25px 0px rgba(0, 0, 0, 0.12), 0 8px 10px -5px rgba(0, 0, 0, 0.2)'
  },
  infoIcon: {
    color: '#00d3ee'
  },
  successIcon: {
    color: '#5cb860'
  },
  warningIcon: {
    color: '#ffa21a'
  },
  errorIcon: {
    color: '#ffa21a'
  },
  dangerIcon: {
    color: '#f55a4e'
  },
  primaryIcon: {
    color: '#af2cc5'
  },
  roseIcon: {
    color: '#eb3573'
  },
  iconMessage: {
    paddingLeft: '50px',
    display: 'block'
  }
});

export default snackbarContentStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/snackbarContentStyle.jsx
