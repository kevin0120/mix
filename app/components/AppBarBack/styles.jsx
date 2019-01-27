export default theme => ({
  root: {
    zIndex: theme.zIndex.drawer + 1
  },
  appBarTool: {
    padding: 0
  },
  backButton: {
    width: 130,
    height: 60,
    borderRight: '1px solid #bdbdbd',
    borderRadius: 0,
    justifyContent: 'space-around'
  },
  leftIcon: {
    marginRight: theme.spacing.unit
  }
});
