export default theme => ({
  root: {
    zIndex: 1,
    height: '100%',
    overflow: 'hidden',
    position: 'relative',
    display: 'flex',
    color: '#333'
  },
  content: {
    flexGrow: 1,
    background: '#f2f2f2',
    padding: theme.spacing.unit * 3,
    minWidth: 0 // So the Typography noWrap works
  },
  toolbar: theme.mixins.toolbar,
  item: {
    height: '64px'
  },
  title: {
    fontSize: 20,
    marginBottom: 10,
    marginLeft: 20
  },
  wrap: {
    width: 872,
    margin: '10px auto 0',
    fontSize: 14
  }
});
