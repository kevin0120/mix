import { put, race, take } from 'redux-saga/effects';
import type { ClsOperationPoint } from './classes/ClsOperationPoint';
import { CommonLog } from '../../../common/utils';
import dialogActions from '../../dialog/action';
import screwStepActions from './action';
import { SCREW_STEP } from './constants';
import { orderActions } from '../../order/action';
import { STEP_STATUS } from '../constants';

export function* byPassPoint(finalFailPoints) {
  try {
    const n: string = finalFailPoints
      .map((p: ClsOperationPoint) => p.point.nut_no)
      .join(',');
    if (finalFailPoints.length > 0) {
      CommonLog.Debug('Show Next Point By Pass Diag');
      yield put(
        dialogActions.dialogShow({
          buttons: [
            {
              label: 'Order.Next',
              color: 'danger',
              action: screwStepActions.confirmFailSpecPoint()

            },
            {
              label: 'Screw.Next',
              color: 'warning',
              action: screwStepActions.byPassSpecPoint()
            }
          ],
          // eslint-disable-next-line camelcase
          title: `拧紧点失败：${n}`,
          content: `${this.failureMsg}`
        })
      );
      const { bypass, fail } = yield race({
        bypass: take(SCREW_STEP.BYPASS_SPEC_POINT),
        fail: take(SCREW_STEP.CONFIRM_FAIL_SPEC_POINT)
      });
      if (fail) {
        yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: '拧紧失败' })); // 失败退出
      }
    }
  } catch (e) {
    CommonLog.lError(e);
    throw e;
  }
}
