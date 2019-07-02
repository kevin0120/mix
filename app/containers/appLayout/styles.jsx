import pagesStyle from '../../common/jss/layouts/pagesStyle';
import imagesStyles from '../../common/jss/imagesStyles';
import popoverStyles from '../../common/jss/popoverStyles';
import { dangerColor, successColor } from '../../common/jss/material-react-pro';

import { keyframes } from 'react-emotion';

const twinkling = keyframes`
  0% {
    opacity: .1;
  }
  100% {
    opacity: 1;
  }
`;

export default theme => ({
  ...pagesStyle,
  ...imagesStyles,
  ...popoverStyles,
  layout: {
    height: '64px',
    flexGrow: 1
  },
  appBar: {
    height: '64px',
    width: '100%',
    bottom:0,
    position:'relative',
    background: '#353744',
    color: '#FFFFFF',
    padding: '0 30px',
    display:'flex',
    flexDirection:'row',
    justifyContent:'space-between',
    alignItems:'center'
  },
  BottomNavigation: {
    flex:1,
    maxWidth:'60%',
    background: 'transparent',
    display:'flex',
    justifyContent:'space-around',
  },
  menuBtnWrapAvatar: {
    display:'flex',
    width: '50px',
    height: '50px',
    borderStyle:'solid',
    borderRadius:'50%',
    borderWidth:'1px',
    overflow:'hidden'
  },
  menuUserName: {
    // width: '80px',
    height: '100%',
    lineHeight: '100%',
    display:'flex',
    justifyContent:'center',
    alignItems:'center',
    marginLeft:'5px',
    marginRight:'5px',
    '& p': {
      fontSize: '16px',
      lineHeight: '100%',
      textAlign: 'center',
      margin: '0'
    }
  },
  menuClock: {
    display:'flex',
    justifyContent:'center',
    alignItems:'center',
    fontSize: '20px',
    // width: '100px',
    marginLeft:'5px',
    marginRight:'5px',
    height: '100%',
    lineHeight: '100%',
  },
  BottomNavigationIcon: {
    color: '#FFFFFF',
    // marginLeft: '1vw',
    fontSize: '20px',
    fontWeight: 'bold',
    paddingLeft: 0,
    paddingRight:0,
    minWidth:'56px'
  },
  menuBtnWrapLeft: {
    marginRight: 'auto',
    marginLeft: 'auto',
    flex: 1
  },
  menuBtnWrapRight: {
    flex: 1,
    maxWidth: '30%',
    textAlign: 'right'
  },
  menuBtn: {
    width: 'auto'
  },
  menuStatusOK: {
    background: successColor,
    fontSize: '18px',
    '&,&:focus,&:hover': {
      background: successColor,
      fontSize: '18px',
    }
  },
  menuStatusFail: {
    background: dangerColor,
    fontSize: '18px',
    animation: `${twinkling} 2s infinite cubic-bezier(1, 1, 1, 1)`,
    '&,&:focus,&:hover': {
      background: dangerColor,
      fontSize: '18px',
      animation: `${twinkling} 2s infinite cubic-bezier(1, 1, 1, 1)`
    }
  },
  itemWrap: {
    margin: '0 8px'
  },
  navTitle: {
    fontSize: '24px'
  },
  menuButton: {},
  menuItem: {
    paddingLeft: '32px',
    paddingTop: '0px',
    paddingBottom: '0px',
    margin: '0',
    height: '50px',
    color: '#333',
    '&:hover': {
      fontWeight: 'bold'
    },
    '&:first-child': {
      marginTop: '5px'
    }
  },
  sideNav: {
    width: '250px',
    padding: 0
  }
});
