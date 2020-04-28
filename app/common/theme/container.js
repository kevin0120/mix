export default {
  common: {
    paddingRight: '15px',
    paddingLeft: '15px',
    marginRight: 'auto',
    marginLeft: 'auto',
    '@media (max-width: 768px)': {
      width: '100%'
    },
    '@media (min-width: 768px)': {
      width: '750px'
    },
    '@media (min-width: 992px)': {
      width: '970px'
    },
    '@media (min-width: 1200px)': {
      width: '1170px'
    },
    '@media (min-width: 1600px)': {
      width: '1370px'
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
}
