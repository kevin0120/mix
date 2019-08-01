import * as colors from '@material-ui/core/colors';

const primaryColor = colors.cyan[700];
const secondaryColor = colors.amber[700];
const warningColor = colors.orange[800];
const dangerColor = colors.red[700];
const successColor = colors.lightGreen[500];
const infoColor = '#00acc1';
const roseColor = '#e88e63';
const grayColor = '#999999';
const customerColor = '#36D7B7';



export default{
  type: 'light',
  primary: {
    main: primaryColor
  },
  secondary: {
    main: secondaryColor
  },
  danger: {
    main: dangerColor
  },
  warning: {
    main: warningColor
  },
  success: {
    main: successColor
  },
  info: {
    main: infoColor
  },
  rose: {
    main: roseColor
  },
  gray: {
    main: grayColor
  },
  customer:{
    main:customerColor
  },
}
