import pagesStyle from "../../common/jss/layouts/pagesStyle.jsx";
import imagesStyles from "../../common/jss/imagesStyles.jsx";



export default theme => ({
  ...pagesStyle,
  ...imagesStyles,
  layout: {
    height: '100%',
    flexGrow: 1,
  },
  appBar: {
    height: '64px',
    width: '100%',
    top: 'auto',
    bottom: 0,
    background: '#353744',
  },
  BottomNavigation: {
    background: 'transparent',
    // width: '600px',
  },
  menuBtnWrapAvatar: {
    width: '50px',
    marginRight: '15px',
  },
  menuUserName: {
    width: '100px',
    height: '64px',
    '& p': {
      marginTop: '25px',
      fontSize: '16px',
      textAlign: 'left',
    },

  },
  BottomNavigationIcon: {
    color: '#FFFFFF',
    marginLeft: '20px',
  },
  menuBtnWrapLeft: {
    marginLeft:'20px',
    flex: 1,
    textAlign: 'left',
  },
  menuBtnWrapRight: {
    flex: 1,
    textAlign: 'right',
  },
  menuBtn: {
    width: 'auto',
  },
  menuStatusOK: {
    background: 'transparent',
    fontSize: '18px',
    marginRight: '20px'
  },
  menuStatusFail: {
    background: '#ff8000',
    fontSize: '18px',
    marginRight: '20px'
  },
  itemWrap: {
    margin: '0 8px',
  },
  topBar: {
    padding: '0 30px',
  },
  navTitle: {
    fontSize: '24px',
  },
  menuButton: {
  },
  menuItem: {
    paddingLeft: '32px',
    paddingTop: '0px',
    paddingBottom: '0px',
    margin: '0',
    height: '50px',
    color: '#333',
    '&:hover': {
      fontWeight: 'bold',
    },
    '&:first-child': {
      marginTop: '5px',
    },
  },
  sideNav: {
    width: '250px',
    padding: 0,
  },
});
