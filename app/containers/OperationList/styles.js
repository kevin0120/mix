const cardHeight='120px';

export default (theme) => ({
  root: {
    width: '100%',
    height: '100%',
    // padding: '20px',
    backgroundColor: '#232c39',
    backgroundImage: 'linear-gradient(45deg, rgba(150, 216, 255, 0.5) 10%, rgba(150, 160, 187, 0.7))'
  },
  container: {
    width: '100%',
    height: '100%',
    margin: 0
  },
  listTitle: {
    color: '#fdfdfd'
  },
  orderCardContainer: {
    height: cardHeight
  },
  orderCard: {
    width: '100%',
    height: cardHeight,
    display:'flex',
    flexDirection:'row',
    justifyContent:'flex-start'
  },
  bgOdd: {
    backgroundColor: 'rgb(135,135,135)'
  },
  bgEven: {
    backgroundColor: 'rgb(143,143,143)'
  },
  image: {
    height:cardHeight,
    width:'30%'
  },
  info:{
    display:'flex',
    flexDirection:'column',
    padding: '10px 5px',
    height:cardHeight,
    flex:1,
    textAlign: 'center',
    alignItems:'flex-start',
    justifyContent: 'center',
    alignContent:'center'
  },
  statusIndicator:{
    display:'block',
    width:'3px',
    height:cardHeight
  },
  orderNameText:{

  },
  orderInfoText:{
    // maxHeight:'50%',
    flex:1,
    textOverflow: 'ellipsis',
    overflow:'hidden',
    display:'block',
  },
  orderStatusText:{

  },
  statusTodo:{
    backgroundColor:'#aaaaaa',
  },
  statusWIP:{
    backgroundColor:'blue',
  },
  statusDone:{
    backgroundColor:'green',
  },
  statusCancel:{
    backgroundColor:'red',
  },
  statusPending:{
    backgroundColor:'orange',
  },
});
