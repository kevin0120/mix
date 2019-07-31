import { createMuiTheme } from '@material-ui/core';
import * as colors from '@material-ui/core/colors';

const primaryColor = colors.cyan[700];
const secondaryColor = colors.amber[700];
const warningColor = colors.orange[800];
const dangerColor = colors.red[700];
const successColor = colors.lightGreen[500];
const infoColor = '#00acc1';
const roseColor = '#e88e63';
const grayColor = '#999999';

const theme = createMuiTheme({
  palette: {
    type: 'light',
    primary: {
      main: primaryColor
    },
    secondary: {
      main: secondaryColor
    },
    danger:{
      main:dangerColor,
    },
    warning:{
      main:warningColor
    },
    success:{
      main:successColor
    },
    info:{
      main:infoColor
    },
    rose:{
      main:roseColor
    },
    gray:{
      main:grayColor
    },
  },
  typography: {
    useNextVariants: true,
    // Use the system font instead of the default Roboto font.
    fontFamily: [
      'Noto Sans SC',
      'sans-serif'
    ].join(','),
    // fontWeightRegular: 'bold',
    button: {
      fontSize: '20px'
    }
  },
  overrides:{
    MuiButtonBase:{

    },
  },
  zIndex: {
    mobileStepper: 1000,
    drawer: 1100,
    appBar: 1200,
    modal: 1300,
    snackbar: 1400,
    tooltip: 1500
  }
});
export default theme;
