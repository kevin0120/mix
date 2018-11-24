// @flow
import { CARD_AUTH } from './actionTypes';

const pcsclite = require('pcsclite');

// 设置授权卡监听器
export const setCardAuthListener = dispatch => {
  const pcscd = pcsclite();

  pcscd.on('reader', reader => {
    readerProcess(reader, dispatch);
    dispatch(cardAuthReaderInserted());

    reader.on('error', (err) => {
      dispatch(cardAuthError(err));
    });

    reader.on('status', (status) => {
      /* check what has changed */
      const changes = reader.state ^ status.state;
      if (changes) {
        if ((changes & reader.SCARD_STATE_EMPTY) && (status.state & reader.SCARD_STATE_EMPTY)) {
          /* card removed */
          return cardAuthReaderStateChange('card removed', reader, dispatch);
        } else if ((changes & reader.SCARD_STATE_PRESENT) && (status.state & reader.SCARD_STATE_PRESENT)) {
          /* card inserted */
          return cardAuthReaderStateChange('card inserted', reader, dispatch);
        }
      }
    });

    reader.on('end', () => {
      console.log('Reader', 'removed');
      dispatch(cardAuthReaderRemoved());
    });
  });

  pcscd.on('error', err => {
    console.log('PCSC error', err.message);
    dispatch(cardAuthError(err));
  });
};

const readerProcess = (reader, dispatch) => {
  reader.on('error', err => dispatch(cardAuthError(err)));
  reader.on('end', () => {
    dispatch(cardAuthReaderRemoved());
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
export const cardAuthReaderReadSuccess = (data) => ({
  type: CARD_AUTH.CARD.READ.SUCCESS,
  data
});

// 授权卡错误
export const cardAuthError = error => {
  return {
    type: CARD_AUTH.ERROR,
    error
  };
};

// 授权卡状态变化
export const cardAuthReaderStateChange = (state, reader, dispatch) => {
  if (state === 'card removed') {
    dispatch({
      type: CARD_AUTH.CARD.REMOVED
    });
    // dispatch(lockScreen());
    reader.disconnect(reader.SCARD_LEAVE_CARD, (err) => {
      if (err) {
        dispatch(cardAuthError(err));
      } else {
        // console.log('disconnected');
      }
    });

  } else {
    dispatch({
      type: CARD_AUTH.CARD.INSERTED
    });
    reader.connect({ share_mode: reader.SCARD_SHARE_SHARED }, (err, protocol) => {
      if (err) {
        dispatch(cardAuthError());
      } else {
        // read data
        dispatch(cardAuthReaderReadStarted());
        reader.transmit(new Buffer([0xff, 0xb2, 0x00, 0x00, 0x04]), 40, protocol, (err, data) => {
          if (err) {
            dispatch(cardAuthError(err));
          } else {
            const uuid = parseBuffer(data);
            dispatch(cardAuthReaderReadSuccess(uuid));
            // dispatch(login(uuid));
          }
        });
      }
    });

  }
};

const parseBuffer = data => {
  let uuid='112233';
  console.log(data.toString('utf8',0,-2));
  // uuid=data.toString('utf8',0,-2);
  // TODO

  return uuid;
};

export const alertRemoveCard = () => (dispatch, getState) => {
  if (getState().authCardInserted) {
  }
};
