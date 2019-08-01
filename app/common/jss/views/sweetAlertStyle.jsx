// ##############################
// // // SweetAlert view styles
// #############################

import buttonStyle from '../../jss/components/buttonStyle.jsx';

const sweetAlertStyle = theme=> ({
  cardTitle: {
    marginTop: '0',
    marginBottom: '3px',
    color: '#3C4858',
    fontSize: '18px'
  },
  center: {
    textAlign: 'center'
  },
  right: {
    textAlign: 'right'
  },
  left: {
    textAlign: 'left'
  },
  ...buttonStyle(theme)
});

export default sweetAlertStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/views/sweetAlertStyle.jsx
