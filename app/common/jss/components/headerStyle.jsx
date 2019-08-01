// ##############################
// // // Header styles
// #############################

const headerStyle = theme => ({
  appBar: {
    backgroundColor: 'transparent',
    boxShadow: 'none',
    borderBottom: '0',
    marginBottom: '0',
    position: 'absolute',
    width: '100%',
    paddingTop: '10px',
    zIndex: '1029',
    color: '#555555',
    border: '0',
    borderRadius: '3px',
    padding: '10px 0',
    transition: 'all 150ms ease 0s',
    minHeight: '50px',
    display: 'block'
  },
  container: {
    ...theme.container.fluid,
    minHeight: '50px'
  },
  flex: {
    flex: 1
  },
  title: {
    fontFamily:theme.typography.fontFamily,
    fontWeight:theme.typography.fontWeightRegular,
    lineHeight: '30px',
    fontSize: '18px',
    borderRadius: '3px',
    textTransform: 'none',
    color: 'inherit',
    paddingTop: '0.625rem',
    paddingBottom: '0.625rem',
    margin: '0 !important',
    '&:hover,&:focus': {
      background: 'transparent'
    }
  },
  primary: {
    backgroundColor: theme.palette.primary.main,
    color: '#FFFFFF',
    ...theme.boxShadow.default
  },
  info: {
    backgroundColor: theme.palette.info.main,
    color: '#FFFFFF',
    ...theme.boxShadow.default
  },
  success: {
    backgroundColor: theme.palette.success.main,
    color: '#FFFFFF',
    ...theme.boxShadow.default
  },
  warning: {
    backgroundColor: theme.palette.warning.main,
    color: '#FFFFFF',
    ...theme.boxShadow.default
  },
  danger: {
    backgroundColor: theme.palette.danger.main,
    color: '#FFFFFF',
    ...theme.boxShadow.default
  },
  sidebarMinimize: {
    float: 'left',
    padding: '0 0 0 15px',
    display: 'block',
    color: '#555555'
  },
  sidebarMinimizeRTL: {
    padding: '0 15px 0 0 !important'
  },
  sidebarMiniIcon: {
    width: '20px',
    height: '17px'
  }
});

export default headerStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/headerStyle.jsx
