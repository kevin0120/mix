
const cardIconStyle = theme=>({
  cardIcon: {
    '&$warningCardHeader,&$successCardHeader,&$dangerCardHeader,&$infoCardHeader,&$primaryCardHeader,&$roseCardHeader': {
      borderRadius: '3px',
      backgroundColor: '#999',
      padding: '15px',
      marginTop: '-20px',
      marginRight: '15px',
      float: 'left'
    }
  },
  warningCardHeader:theme.cardHeader.warning,
  successCardHeader:theme.cardHeader.success,
  dangerCardHeader:theme.cardHeader.danger,
  infoCardHeader:theme.cardHeader.info,
  primaryCardHeader:theme.cardHeader.primary,
  roseCardHeader:theme.cardHeader.rose
});

export default cardIconStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/cardIconStyle.jsx
