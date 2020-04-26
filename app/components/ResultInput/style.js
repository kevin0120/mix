

const styles = theme => ({
  root: {
    height: '100%',
    width: '100%',
  },
  radioGroup:{
    height: '100%',
    width: '100%',
    display:'flex',
    flexDirection:'row',
    justifyContent:'space-around',
    alignItems:'space-around'
  },
  inputLabel:{
    width: '100%',
    height: '100%',
    display:'flex',
    flexDirection:'column',
    justifyContent:'center',
    alignItems:'center'
  },
  inputContainer:{
    height: '100px',
    width: '100%',
    padding: '2px 10px',
    margin:'2px',
    textAlign: 'center',
    display:'flex',
    flexDirection:'column',
    justifyContent:'center',
    alignItems:'center'
  }
});

export default styles;
