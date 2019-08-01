import { createMuiTheme } from '@material-ui/core';
import cardHeader from './cardHeader';
import boxShadow from './boxShadow';
import palette from './palette';


const title = {
  color: '#3C4858',
  textDecoration: 'none',
  fontWeight: '300',
  marginTop: '30px',
  marginBottom: '25px',
  minHeight: '32px',
  fontFamily: '\'Roboto\', \'Helvetica\', \'Arial\', sans-serif',
  '& small': {
    color: '#777',
    fontSize: '65%',
    fontWeight: '400',
    lineHeight: '1'
  }
};

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
  container: {
    common: {
      paddingRight: '15px',
      paddingLeft: '15px',
      marginRight: 'auto',
      marginLeft: 'auto',
      '@media (min-width: 768px)': {
        width: '750px'
      },
      '@media (min-width: 992px)': {
        width: '970px'
      },
      '@media (min-width: 1200px)': {
        width: '1170px'
      },
      '&:before,&:after': {
        display: 'table',
        content: '" "'
      },
      '&:after': {
        clear: 'both'
      }
    },
    fluid: {
      paddingRight: '15px',
      paddingLeft: '15px',
      marginRight: 'auto',
      marginLeft: 'auto',
      '&:before,&:after': {
        display: 'table',
        content: '" "'
      },
      '&:after': {
        clear: 'both'
      }
    }
  },
  title: {
    common: title,
    card: {
      ...title,
      marginTop: '0',
      marginBottom: '3px',
      minHeight: 'auto',
      '& a': {
        ...title,
        marginTop: '.625rem',
        marginBottom: '0.75rem',
        minHeight: 'auto'
      }
    }
  },
  cardHeader
});

export default theme;
