// ##############################
// // // Timeline component styles
// #############################

const timelineStyle = theme => ({
  timeline: {
    [theme.breakpoints.down('sm')]: {
      '&:before': {
        left: '5% !important'
      }
    },
    listStyle: 'none',
    padding: '20px 0 20px',
    position: 'relative',
    '&:before': {
      top: '50px',
      bottom: '0',
      position: 'absolute',
      content: '" "',
      width: '3px',
      backgroundColor: '#E5E5E5',
      left: '50%',
      marginLeft: '18px'
    }
  },
  timelineSimple: {
    padding: '0 0 10px',
    '&:before': {
      left: '5%'
    },
    marginTop: '0'
  },
  item: {
    marginBottom: '5px',
    position: 'relative',
    '&:before,&:after': {
      content: '" "',
      display: 'table'
    },
    '&:after': {
      clear: 'both'
    }
  },
  timelineBadge: {
    [theme.breakpoints.down('sm')]: {
      left: '5% !important'
    },
    color: '#FFFFFF',
    width: '50px',
    height: '50px',
    lineHeight: '51px',
    fontSize: '1.4em',
    textAlign: 'center',
    position: 'absolute',
    top: '10px',
    left: '50%',
    marginLeft: '-5px',
    zIndex: '100',
    borderTopRightRadius: '50%',
    borderTopLeftRadius: '50%',
    borderBottomRightRadius: '50%',
    borderBottomLeftRadius: '50%'
  },
  timelineSimpleBadge: {
    left: '5%'
  },
  info: {
    backgroundColor: theme.palette.info.main,
    boxShadow:
      '0 4px 20px 0px rgba(0, 0, 0, 0.14), 0 7px 10px -5px rgba(0, 188, 212, 0.4)'
  },
  success: {
    backgroundColor: theme.palette.success.main,
    boxShadow:
      '0 4px 20px 0px rgba(0, 0, 0, 0.14), 0 7px 10px -5px rgba(76, 175, 80, 0.4)'
  },
  danger: {
    backgroundColor: theme.palette.danger.main,
    boxShadow:
      '0 4px 20px 0px rgba(0, 0, 0, 0.14), 0 7px 10px -5px rgba(244, 67, 54, 0.4)'
  },
  warning: {
    backgroundColor: theme.palette.warning.main,
    boxShadow:
      '0 4px 20px 0px rgba(0, 0, 0, 0.14), 0 7px 10px -5px rgba(255, 152, 0, 0.4)'
  },
  primary: {
    backgroundColor: theme.palette.primary.main,
    boxShadow:
      '0 4px 20px 0px rgba(0, 0, 0, 0.14), 0 7px 10px -5px rgba(156, 39, 176, 0.4)'
  },
  badgeIcon: {
    width: '24px',
    height: '51px'
  },
  timelinePanel: {
    [theme.breakpoints.down('sm')]: {
      float: 'right !important',
      width: '50% !important',
      '&:before': {
        borderLeftWidth: '0 !important',
        borderRightWidth: '15px !important',
        left: '-15px !important',
        right: 'auto !important'
      },
      '&:after': {
        borderLeftWidth: '0 !important',
        borderRightWidth: '14px !important',
        left: '-14px !important',
        right: 'auto !important'
      }
    },
    width: '35%',
    float: 'left',
    padding: '5px',
    marginBottom: '0px',
    marginRight: '10px',
    position: 'relative',
    boxShadow: '0 1px 4px 0 rgba(0, 0, 0, 0.14)',
    borderRadius: '6px',
    color: 'rgba(0, 0, 0, 0.87)',
    background: '#fff',
    '&:before': {
      position: 'absolute',
      top: '20px',
      right: '-15px',
      display: 'inline-block',
      borderTop: '15px solid transparent',
      borderLeft: '15px solid #e4e4e4',
      borderRight: '0 solid #e4e4e4',
      borderBottom: '15px solid transparent',
      content: '" "'
    },
    '&:after': {
      position: 'absolute',
      top: '20px',
      right: '-14px',
      display: 'inline-block',
      borderTop: '14px solid transparent',
      borderLeft: '14px solid #FFFFFF',
      borderRight: '0 solid #FFFFFF',
      borderBottom: '14px solid transparent',
      content: '" "'
    }
  },
  timelineSimplePanel: {
    width: '75%'
  },
  timelinePanelInverted: {
    [theme.breakpoints.up('sm')]: {
      float: 'right',
      backgroundColor: '#FFFFFF',
      '&:before': {
        borderLeftWidth: '0',
        borderRightWidth: '15px',
        left: '-15px',
        right: 'auto'
      },
      '&:after': {
        borderLeftWidth: '0',
        borderRightWidth: '14px',
        left: '-14px',
        right: 'auto'
      }
    }
  },
  timelineHeading: {
    // float: 'right'
  },
  timelineTS: {
    // float: 'left',
    '& span': {
      color: theme.palette.gray.main,
      fontSize: '12px'
    }
  },
  timelineBody: {
    fontSize: '14px',
    lineHeight: '21px',
    '& p': {
      fontSize: '21px',
      marginBottom: '5px'
    }
  },
  timelineFooter: {
    zIndex: '1000',
    position: 'relative',
    float: 'left'
  },
  footerTitle: {
    color: '#333333',
    fontWeight: '400',
    margin: '10px 0px 0px'
  },
  footerLine: {
    marginTop: '10px',
    marginBottom: '5px'
  }
});

export default timelineStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/timelineStyle.jsx
