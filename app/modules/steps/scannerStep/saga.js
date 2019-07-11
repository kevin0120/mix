import { take } from 'redux-saga/effects';
import { SCANNER_STEP } from './action';

export default function* root(ORDER, orderActions) {
  try {
    while (true) {
      const { value } = yield take(SCANNER_STEP.GET_VALUE);

    }
  } catch (e) {
    console.error(e);
  }
}
