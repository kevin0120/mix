import { warningColor } from '../../common/jss/material-react-pro';

export default theme => ({
  root: {
    // flexGrow: 1,
    zIndex: 1,
    overflow: 'hidden',
    position: 'relative',
    display: 'flex',
    color: '#333',
    height: '100%'
  },
  content: {
    flex: 1,
    // backgroundColor: theme.palette.background.default,
    background: '#f2f2f2',
    padding: theme.spacing(2),
    minWidth: 0, // So the Typography noWrap works
    height: '100%',
    paddingTop: '64px'
  },
  toolbar: theme.mixins.toolbar,
  menuItem: {
    fontSize: 14,
    height: 120,
    margin: '15px 5px'
  },
  menuItemSelected: {
    fontSize: 14,
    backgroundColor: warningColor,
    height: 120,
    color: '#FAFAFA',
    margin: '15px 5px'
  },
  itemText: {
    fontSize: 20,
    padding: 0
  },
  cardActionArea: {
    flex: 1,
    color: '#00352c',
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center'
  },
  cardActionAreaSelected: {
    flex: 1,
    color: 'white',
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center'
  }
});
