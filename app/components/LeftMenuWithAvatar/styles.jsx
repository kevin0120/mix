const drawerWidth = 200;

export default theme => ({
  drawerPaper: {
    position: 'relative',
    width: drawerWidth
  },
  toolbar: theme.mixins.toolbar,
  row: {
    display: 'flex',
    margin: '10px 0'
  },
  avatar: {
    marginRight: 10,
    width: 50,
    height: 50
  },
  userInfo: {
    color: '#333',
    fontSize: 12,
    padding: 0
  },
  userText: {
    fontSize: 12,
    paddingLeft: '10px'
  },
  timeWrap: {
    padding: '10px 5px'
  },
  timeContent: {
    margin: '10px 10px 0px',
    fontSize: 20
  },
  baseInfo: {
    boxSizing: 'border-box',
    position: 'absolute',
    width: '100%',
    padding: '10px 20px 20px',
    background: '#fff',
    bottom: 0,
    left: 0
  }
});
