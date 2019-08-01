// ##############################
// // // Footer styles
// #############################


const footerStyle = theme=>({
  block: {},
  left: {
    float: 'left!important',
    display: 'block'
  },
  right: {
    margin: '0',
    fontSize: '14px',
    float: 'right!important',
    padding: '15px'
  },
  footer: {
    bottom: '0',
    borderTop: '1px solid #e7e7e7',
    padding: '15px 0',
    fontFamily:theme.typography.fontFamily,
    fontWeight:theme.typography.fontWeight,
    fontSize:theme.typography.fontSize,
    zIndex: 4
  },
  container: {
    zIndex: 3,
    ...theme.container.common,
    position: 'relative'
  },
  containerFluid: {
    zIndex: 3,
    ...theme.container.fluid,
    position: 'relative'
  },
  a: {
    color: theme.palette.primary.main,
    textDecoration: 'none',
    backgroundColor: 'transparent'
  },
  list: {
    marginBottom: '0',
    padding: '0',
    marginTop: '0'
  },
  inlineBlock: {
    display: 'inline-block',
    padding: '0',
    width: 'auto'
  },
  whiteColor: {
    '&,&:hover,&:focus': {
      color: '#FFFFFF'
    }
  }
});
export default footerStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/footerStyle.jsx
