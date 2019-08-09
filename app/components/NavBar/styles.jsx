import { keyframes } from '@emotion/core';
import imagesStyles from '../../common/jss/imagesStyles';
import popoverStyles from '../../common/jss/popoverStyles';

const twinkling = keyframes`
  0% {
    opacity: .1;
  }
  100% {
    opacity: 1;
  }
`;

const pages = {
  BottomNavigation: {
    flex: 1,
    maxWidth: '60%',
    background: 'transparent',
    display: 'flex',
    justifyContent: 'space-around'
  },
  BottomNavigationIcon: {
    color: '#FFFFFF',
    fontSize: '20px',
    fontWeight: 'bold',
    paddingLeft: 0,
    paddingRight: 0,
    minWidth: '56px'
  }
};

const clock = {
  menuClock: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    fontSize: '20px',
    // width: '100px',
    marginLeft: '5px',
    marginRight: '5px',
    height: '100%',
    lineHeight: '100%'
  },
  timeContent: {}
};

const avatar = {
  menuBtnWrapAvatar: {
    display: 'flex',
    width: '70px',
    height: '50px'
  }
};

const languageMenu = {
  ...imagesStyles,
  ...popoverStyles,
  menuItem: {
    // paddingLeft: '32px',
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
  }
};

const NavBarMenu = theme => ({
  ...popoverStyles,
  menuStatusOK: {
    background: theme.palette.success.main,
    fontSize: '18px',
    '&,&:focus,&:hover': {
      background: theme.palette.success.main,
      fontSize: '18px'
    }
  },
  menuStatusFail: {
    background: theme.palette.danger.main,
    fontSize: '18px',
    animation: `$twinkling 2s infinite cubic-bezier(1,1,1,1)`,
    '&,&:focus,&:hover': {
      background: theme.palette.danger.main,
      fontSize: '18px',
    }

  },
  '@keyframes twinkling': {
    '0%': {
      opacity: 0.1
    },
    '100%': {
      opacity: 1
    }
  }
});
const root = {
  appBar: {
    height: '64px',
    width: '100%',
    bottom: 0,
    position: 'relative',
    background: '#353744',
    // backgroundImage: 'url("../resources/imgs/texture.png")',
    // backgroundRepeat: 'repeat',
    backgroundColor: '#333',
    color: '#FFFFFF',
    padding: '0 30px',
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center'
  }
};

export default {
  root,
  pages,
  clock,
  avatar,
  languageMenu,
  NavBarMenu
};
