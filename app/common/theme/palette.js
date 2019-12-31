import * as colors from '@material-ui/core/colors';

export const primaryColor = colors.cyan[700];
export const secondaryColor = colors.amber[700];
export const warningColor = colors.orange[800];
export const dangerColor = colors.red[700];
export const successColor = colors.lightGreen[500];
export const infoColor = colors.cyan[700];
export const roseColor = '#e88e63';
export const grayColor = '#999999';
export const customerColor = '#36D7B7';

export default{
  type: 'light',
  primary: {
    ...colors.cyan,
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
    ...colors.cyan,
    main: infoColor
  },
  rose: {
    main: roseColor
  },
  gray: {
    ...colors.grey,
    main: grayColor,
  },
  customer:{
    main:customerColor
  },
}
