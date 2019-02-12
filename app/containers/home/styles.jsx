/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

const styles = theme => ({
  wrap: {
    backgroundColor: '#f2f2f2',
    // height: 'calc(100% - 64px)',
    height:'100%',
    overflowX: 'hidden'
  },
  iconWrap: {
    position: 'absolute',
    top: 0,
    left: '50%',
    transform: 'translate(-50%, -50%)',
    width: 80,
    height: 80,
    backgroundColor: '#fff',
    borderRadius: '99%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center'
  },
  root: {
    display: 'flex',
    flexWrap: 'wrap',
    minWidth: 300,
    width: '100%',
    height: '100%',
    backgroundColor: '#232c39',
    backgroundImage:
      'linear-gradient(45deg, rgba(150, 216, 255, 0.5) 10%, rgba(150, 160, 187, 0.7))'
  },
  container: {
    padding: '40px 10px',
    textAlign: 'center',
    marginTop: 0
  },
  card: {
    width: '15%',
    minWidth: '280px',
    height: '320px',
    borderRadius: '4px'
  },
  cardActionArea: {
    height: '100%'
  },
  cardGridItem: {
    paddingRight: '60px'
  },
  cardContent: {
    position: 'relative',
    padding: '50px 10px'
  },
  fabLeft: {
    fontSize: 20,
    position: 'absolute',
    bottom: theme.spacing.unit * 10,
    left: theme.spacing.unit * 2
  },
  fabRight: {
    position: 'absolute',
    bottom: theme.spacing.unit * 10,
    right: theme.spacing.unit * 2,
    fontSize: 20
  },
  fabMoveUp: {
    transform: 'translate3d(0, -46px, 0)',
    transition: theme.transitions.create('transform', {
      duration: theme.transitions.duration.enteringScreen,
      easing: theme.transitions.easing.easeOut
    })
  },
  fabMoveDown: {
    transform: 'translate3d(0, 0, 0)',
    transition: theme.transitions.create('transform', {
      duration: theme.transitions.duration.leavingScreen,
      easing: theme.transitions.easing.sharp
    })
  },
  media: {
    background: '#ddd',
    paddingTop: '50%', // 16:9
    width: '100%'
  },
  title: {
    fontSize: '34px',
    color: '#fff',
    marginBottom: '20px'
  },
  subTitle: {
    fontSize: '16px',
    color: '#fff'
  },
  icon: {
    fontSize: '50px'
  },
  btnWrap: {
    position: 'relative',
    height: 200,
    [theme.breakpoints.down('xs')]: {
      width: '100% !important', // Overrides inline-style
      height: 100
    },
    '&:hover, &$focusVisible': {
      zIndex: 1,
      opacity: 0.65,
      '& $titleMarked': {
        opacity: 0
      },
      '& $navTitle': {
        border: '4px solid currentColor'
      }
    }
  },
  focusVisible: {},
  extendedIcon: {
    marginRight: theme.spacing.unit
  },
  itemButton: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: theme.palette.common.white
  },
  navTitle: {
    position: 'relative',
    padding: `${theme.spacing.unit * 2}px ${theme.spacing.unit * 4}px ${theme
      .spacing.unit + 6}px`
  },
  titleMarked: {
    height: 3,
    width: 18,
    backgroundColor: theme.palette.common.white,
    position: 'absolute',
    bottom: -2,
    left: 'calc(50% - 9px)',
    transition: theme.transitions.create('opacity')
  }
});

export default styles;
