import boxShadow from './boxShadow';

const warning = {
  background: 'linear-gradient(60deg, #ffa726, #fb8c00)',
  ...boxShadow.warning
};
const success = {
  background: 'linear-gradient(60deg, #66bb6a, #43a047)',
  ...boxShadow.success
};
const danger = {
  background: 'linear-gradient(60deg, #ef5350, #e53935)',
  ...boxShadow.danger
};
const info = {
  background: 'linear-gradient(60deg, #26c6da, #00acc1)',
  ...boxShadow.info
};
const primary = {
  background: 'linear-gradient(60deg, #ab47bc, #8e24aa)',
  ...boxShadow.primary
};
const rose = {
  background: 'linear-gradient(60deg, #ec407a, #d81b60)',
  ...boxShadow.rose
};
export default{
  warning,
  success,
  danger,
  info,
  primary,
  rose
}
