// ##############################
// // // Badge component styles
// #############################

const badgeStyle = theme=>({
  badge: {
    borderRadius: '12px',
    padding: '5px 12px',
    textTransform: 'uppercase',
    fontSize: '10px',
    fontWeight: '700',
    lineHeight: '1',
    color: '#fff',
    textAlign: 'center',
    verticalAlign: 'baseline',
    display: 'inline-block'
  },
  primary: {
    backgroundColor: theme.palette.primary.main
  },
  warning: {
    backgroundColor: theme.palette.warning.main
  },
  danger: {
    backgroundColor: theme.palette.danger.main
  },
  success: {
    backgroundColor: theme.palette.success.main
  },
  info: {
    backgroundColor: theme.palette.info.main
  },
  rose: {
    backgroundColor: theme.palette.rose.main
  },
  gray: {
    backgroundColor: theme.palette.gray.main
  }
});

export default badgeStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/badgeStyle.jsx
