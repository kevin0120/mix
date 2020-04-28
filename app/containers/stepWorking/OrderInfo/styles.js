const OrderInfoStyles = (theme) => ({
  root: {
    margin: 3,
    maxWidth: '100%',
    flex: 1,
    overflow: 'auto',
    backgroundColor: theme.palette.common.white,
    padding: 2
  },
  panelRoot: {
    overflow: 'auto',
    maxWidth: '100%',
    margin: 2,
    display: 'flex',
    flexDirection: 'column',
    boxSizing: 'content-box'
  },
  titleContainer: {
    width: '100%',
    backgroundColor: theme.palette.primary[200],
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    color: theme.palette.gray[900],
    padding: 2
  },
  table: {
    flex: 1,
    width: 'auto',
    borderStyle: 'dashed',
    borderColor: theme.palette.grey[300],
    borderWidth: 2,
    borderRadius: 3,
    borderTop: 'none !important'
  },
  tableCell: {
    wordBreak: 'break-all',
    minWidth: '70px',
    textAlign: 'center'
  }
});

export default {
  OrderInfoStyles
};
