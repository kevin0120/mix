// ##############################
// // // Pages Header styles
// #############################

const pagesHeaderStyle = theme => ({
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
    ...theme.container.common,
    minHeight: '50px'
  },
  flex: {
    flex: 1
  },
  title: {
    fontWeight:theme.typography.fontWeightRegular,
    fontFamily:theme.typography.fontFamily,
    lineHeight: '30px',
    fontSize: '18px',
    borderRadius: '3px',
    textTransform: 'none',
    color: '#FFFFFF',
    '&:hover,&:focus': {
      background: 'transparent',
      color: '#FFFFFF'
    }
  },
  indicator: {
    fontWeight:theme.typography.fontWeightRegular,
    fontFamily:theme.typography.fontFamily,
    fontSize: '14px',
    borderRadius: '3px',
    textTransform: 'none',
    color: '#FFFFFF',
    '&:hover,&:focus': {
      background: 'transparent',
      color: '#FFFFFF'
    }
  },
  appResponsive: {
    top: '8px'
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
  list: {
    fontWeight:theme.typography.fontWeightRegular,
    fontFamily:theme.typography.fontFamily,
    fontSize: '14px',
    margin: 0,
    marginRight: '-15px',
    paddingLeft: '0',
    listStyle: 'none',
    color: '#FFFFFF',
    paddingTop: '0',
    paddingBottom: '0'
  },
  listItem: {
    // float: 'left',
    // position: 'relative',
    display: 'flex',
    flexDirection:'row',
    justifyContent:'center',
    alignItems:'center',
    width: 'auto',
    margin: '0',
    padding: '0',
    [theme.breakpoints.down('sm')]: {
      zIndex: '999',
      width: '100%',
      paddingRight: '15px'
    }
  },
  navLink: {
    // display: 'flex',
    // flexDirection:'row',
    // justifyContent:'center',
    // alignItems:'center',
    // color: '#FFFFFF',
    // margin: '0 5px',
    // paddingTop: '15px',
    // paddingBottom: '15px',
    fontWeight: '400',
    // fontSize: '14px',
    // textTransform: 'uppercase',
    // borderRadius: '3px',
    // padding: '10px 15px',
    // textDecoration: 'none',
    // '&:hover,&:focus': {
    //   color: '#FFFFFF',
    //   background: 'rgba(200, 200, 200, 0.2)'
    // }
  },
  listItemIcon: {
    position: 'relative',
    marginRight: '3px',
    width: '24px',
    height: '24px',
    color: 'inherit',
  },
  listItemText: {
    flex: 'none',
    padding: '0',
    margin: 0,
    whiteSpace: 'nowrap'
  },
  navLinkActive: {
    backgroundColor: 'rgba(255, 255, 255, 0.1)'
  },
  drawerPaper: {
    border: 'none',
    bottom: '0',
    transitionProperty: 'top, bottom, width',
    transitionDuration: '.2s, .2s, .35s',
    transitionTimingFunction: 'linear, linear, ease',
    ...theme.boxShadow.common,
    width: theme.shape.drawerWidth,
    ...theme.boxShadow.common,
    position: 'fixed',
    display: 'block',
    top: '0',
    height: '100vh',
    right: '0',
    left: 'auto',
    visibility: 'visible',
    overflowY: 'visible',
    borderTop: 'none',
    textAlign: 'left',
    paddingRight: '0px',
    paddingLeft: '0',
    ...theme.transitions.all,
    '&:before,&:after': {
      position: 'absolute',
      zIndex: '3',
      width: '100%',
      height: '100%',
      content: '""',
      display: 'block',
      top: '0'
    },
    '&:after': {
      background: '#000',
      opacity: '.8'
    }
  },
  sidebarButton: {
    '&,&:hover,&:focus': {
      color: '#FFFFFF'
    },
    top: '-2px'
  }
});

export default pagesHeaderStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/pagesHeaderStyle.jsx
