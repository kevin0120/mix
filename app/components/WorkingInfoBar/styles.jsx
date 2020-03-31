export default () => ({
  infoText: {
    textAlign: 'right',
    padding: 0,
    flex: 'initial',
    color: '#888',
    fontSize: 14,
    marginBottom: '0',
    display:'flex',
    flexDirection:'row',
    justifyContent:'space-between',
    alignItems:'center',
  },
  infoItem: {
    width: '100%',
    flex:1,
    // '&:first-child': {
    //   marginTop: '5px'
    // },
    // '&:last-child': {
    //   marginBottom: '5px'
    // },
    display:'flex',
    flexDirection:'row',
    justifyContent:'space-between',
    alignItems:'center',
    padding: '0 20px',
    // marginBottom: '10px',
    // marginTop: '10px'
  },
  itemIitle: {
    height:'80%',
    fontSize: 14
  }
});
