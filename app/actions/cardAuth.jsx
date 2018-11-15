// @flow
import { CARD_AUTH } from './actionTypes';

const pcsclite = require('pcsclite');

// 设置授权卡监听器
export const setCardAuthListener = (dispatch) => {
  const pcscd = pcsclite();

  pcscd.on('reader', (reader) => {
    readerProcess(reader, dispatch);
    dispatch(cardAuthReaderInserted());
  });

  pcscd.on('error', (err) => {
    console.log('PCSC error', err.message);
    dispatch(cardAuthError(err));
  });
};


const readerProcess = (reader, dispatch) => {
  reader.on('error', (err) => (
    dispatch(cardAuthError(err))
  ));
  reader.on('end', () => {
    dispatch(cardAuthReaderRemoved())
  });

};


// 授权卡读卡器插入
const cardAuthReaderInserted = () => ({
    type: CARD_AUTH.READER.INSERTED
});

// 授权卡读卡器拔出
const cardAuthReaderRemoved = () => ({
  type: CARD_AUTH.READER.REMOVED
});



export const cardAuthReaderReadStarted = () => ({
  type: CARD_AUTH.CARD.READ.STARTED
});

// 授权卡读取成功
export const cardAuthReaderReadSuccess = () => ({
  type: CARD_AUTH.CARD.READ.SUCCESS
});

// 授权卡错误
export const cardAuthError = (error) => {
  return{
    type: CARD_AUTH.ERROR,
    error
  }};

const parseBuffer=(data)=>{
  let uuid='112233';
  console.log(data.toString('utf8',0,-2));
  //uuid=data.toString('utf8',0,-2);
  // TODO

  return uuid;
};

export const alertRemoveCard=()=>(dispatch,getState)=>{
  if(getState().authCardInserted){
  }
};
