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
export default {
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
};
