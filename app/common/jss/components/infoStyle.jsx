// ##############################
// // // Info component styles
// #############################

const infoStyle = theme=>({
  infoArea: {
    maxWidth: '360px',
    margin: '0 auto',
    padding: '70px 0 30px'
  },
  iconWrapper: {
    float: 'left',
    marginTop: '24px',
    marginRight: '10px'
  },
  primary: {
    color: theme.palette.primary.main
  },
  warning: {
    color: theme.palette.warning.main
  },
  danger: {
    color: theme.palette.danger.main
  },
  success: {
    color: theme.palette.success.main
  },
  info: {
    color: theme.palette.info.main
  },
  rose: {
    color: theme.palette.rose.main
  },
  gray: {
    color: theme.palette.gray.main
  },
  icon: {
    width: '36px',
    height: '36px'
  },
  descriptionWrapper: {
    color: theme.palette.gray.main,
    overflow: 'hidden'
  },
  title: {
    ...theme.title.common,
    margin: '1.75rem 0 0.875rem !important',
    minHeight: 'unset'
  },
  description: {
    color: theme.palette.gray.main,
    overflow: 'hidden',
    marginTop: '0px',
    '& p': {
      color: theme.palette.gray.main,
      fontSize: '14px'
    }
  },
  iconWrapperVertical: {
    float: 'none'
  },
  iconVertical: {
    width: '61px',
    height: '61px'
  }
});

export default infoStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-kit-pro-react/components/infoStyle.jsx
