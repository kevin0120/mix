const cardHeight='120px';
import {successColor, warningColor, dangerColor,grayColor,primaryColor} from '../../common/jss/material-react-pro';

export default (theme) => ({
  root: {
    width: '100%',
    height: '100%',
    // padding: '20px',
    backgroundColor: '#7c7c7c',
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
    justifyContent:'flex-start',
  },
  bgOdd: {
    backgroundColor: '#747474'
  },
  bgEven: {
    backgroundColor: '#7c7c7c'
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
    width:'5px',
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
    backgroundColor:grayColor,
  },
  statusWIP:{
    backgroundColor:primaryColor,
  },
  statusDone:{
    backgroundColor:successColor,
  },
  statusCancel:{
    backgroundColor:dangerColor,
  },
  statusPending:{
    backgroundColor:warningColor,
  },
});
