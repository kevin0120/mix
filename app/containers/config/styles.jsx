export default theme => ({
  root: {
    flexGrow: 1,
    zIndex: 1,
    overflow: 'hidden',
    position: 'relative',
    display: 'flex',
    color: '#333',
    overflowY: 'scroll',
    maxHeight: 'calc(100% - 64px)'
  },
  content: {
    flexGrow: 1,
    // backgroundColor: theme.palette.background.default,
    background: '#f2f2f2',
    padding: theme.spacing.unit * 2,
    minWidth: 0 // So the Typography noWrap works
  },
  toolbar: theme.mixins.toolbar,
  menuItem: {
    fontSize: 14
  },
  itemText: {
    fontSize: 14,
    padding: 0
  }
});
