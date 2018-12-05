import pagesStyle from '../../common/jss/layouts/pagesStyle';
import imagesStyles from '../../common/jss/imagesStyles';
import popoverStyles from '../../common/jss/popoverStyles';
import { dangerColor } from '../../common/jss/material-react-pro';

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
    height: '100%',
    flexGrow: 1
  },
  appBar: {
    height: '64px',
    width: '100%',
    top: 'auto',
    bottom: 0,
    background: '#353744'
  },
  BottomNavigation: {
    background: 'transparent',
    '& a:last-child': {
      marginRight: '20px'
    },
    // width: '600px',
  },
  menuBtnWrapAvatar: {
    width: '50px',
    // marginRight: '15px',
    height: '50px',
  },
  menuUserName: {
    width: '80px',
    height: '100%',
    lineHeight: '100%',
    '& p': {
      // marginTop: '25px',
      fontSize: '16px',
      lineHeight: '100%',
      textAlign: 'center',
      margin: '0',
    }
  },
  BottomNavigationIcon: {
    color: '#FFFFFF',
    marginLeft: '20px'
  },
  menuBtnWrapLeft: {
    marginRight: 'auto',
    marginLeft: 'auto',
    flex: 1,
    // textAlign: 'center'
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
    background: 'transparent',
    fontSize: '18px',
    marginRight: '20px'
  },
  menuStatusFail: {
    background: dangerColor,
    fontSize: '18px',
    marginRight: '20px',
    animation: `${twinkling} 2s infinite cubic-bezier(1, 1, 1, 1)`
  },
  itemWrap: {
    margin: '0 8px'
  },
  topBar: {
    padding: '0 30px'
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
