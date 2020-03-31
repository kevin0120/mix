import { createMuiTheme } from '@material-ui/core';
import cardHeader from './cardHeader';
import boxShadow from './boxShadow';
import palette from './palette';
import container from './container';
import title from './title';




const theme = createMuiTheme({
  palette,
  typography: {
    useNextVariants: true,
    // Use the system font instead of the default Roboto font.
    fontFamily: [
      'Noto Sans SC',
      'sans-serif'
    ].join(','),
    fontWeightRegular: '400',
    button: {
      fontSize: '20px'
    }
  },
  shape: {
    drawerWidth: '260px'
  },
  transitions: {
    all: {
      transition: 'all 0.33s cubic-bezier(0.685, 0.0473, 0.346, 1)'
    }
  },
  zIndex: {
    mobileStepper: 1000,
    drawer: 1100,
    appBar: 1200,
    modal: 1300,
    snackbar: 1400,
    tooltip: 1500
  },
  boxShadow,
  container,
  title,
  cardHeader
});

export default theme;
