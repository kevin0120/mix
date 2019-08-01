const customInputStyle = theme=>({
  disabled: {
    '&:before': {
      borderColor: 'transparent !important'
    }
  },
  underline: {
    '&:hover:not($disabled):before,&:before': {
      borderColor: '#D2D2D2 !important',
      borderWidth: '1px !important'
    },
    '&:after': {
      borderColor: theme.palette.primary.main
    }
  },
  underlineError: {
    '&:after': {
      borderColor: theme.palette.danger.main
    }
  },
  underlineSuccess: {
    '&:after': {
      borderColor: theme.palette.success.main
    }
  },
  labelRoot: {
    fontFamily:theme.typography.fontFamily,
    color: '#AAAAAA !important',
    fontWeight: '400',
    fontSize: '14px',
    lineHeight: '1.42857',
    top: '10px',
    '& + $underline': {
      marginTop: '0px'
    }
  },
  labelRootError: {
    color: `${theme.palette.danger.main  } !important`
  },
  labelRootSuccess: {
    color: `${theme.palette.success.main  } !important`
  },
  feedback: {
    position: 'absolute',
    bottom: '3px',
    right: '0',
    zIndex: '2',
    display: 'block',
    width: '1em',
    height: '1em',
    textAlign: 'center',
    pointerEvents: 'none'
  },
  feedbackRight: {
    marginRight: '22px'
  },
  formControl: {
    margin: '0 0 17px 0',
    paddingTop: '27px',
    position: 'relative',
    '& svg,& .fab,& .far,& .fal,& .fas': {
      color: '#495057'
    }
  },
  whiteUnderline: {
    '&:hover:not($disabled):before,&:before': {
      backgroundColor: '#FFFFFF'
    },
    '&:after': {
      backgroundColor: '#FFFFFF'
    }
  },
  largeInput: {
    '&,&::placeholder': {
      fontSize: '50px !important'
    }
  },
  input: {
    color: '#495057',
    '&,&::placeholder': {
      fontSize: '14px',
      fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
      fontWeight: '400',
      lineHeight: '1.42857',
      opacity: '1'
    },
    '&::placeholder': {
      color: '#AAAAAA'
    }
  },
  whiteInput: {
    '&,&::placeholder': {
      color: '#FFFFFF',
      opacity: '1'
    }
  }
});

export default customInputStyle;

// WEBPACK FOOTER //
// ./src/assets/jss/material-dashboard-pro-react/components/customInputStyle.jsx
