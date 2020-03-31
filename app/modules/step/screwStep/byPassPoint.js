import { put, race, take } from 'redux-saga/effects';
import { CommonLog } from '../../../common/utils';
import dialogActions from '../../dialog/action';
import screwStepActions from './action';
import { SCREW_STEP } from './constants';
import { orderActions } from '../../order/action';
import { STEP_STATUS } from '../constants';

export function* byPassPoint(controls, retry = false, disableByPass = false) {
  try {
    const n: string = controls
      .map((c) => c.sequence)
      .join(',');
    if (controls.length > 0) {
      CommonLog.Debug('Show Next Point By Pass Diag');
      yield put(
        dialogActions.dialogShow({
          maxWidth: 'md',
          buttons: [
            {
              label: 'Order.Next',
              color: 'danger',
              action: screwStepActions.confirmFailSpecPoint()

            },
            ...(retry ? [{
              label: 'Common.Retry',
              color: 'info',
              action: screwStepActions.bypassRetry()
            }, {
              label: 'Common.Manual',
              color: 'info',
              // action: reworkActions.tryRework(null, null, finalFailPoints[0])
              action: { type: SCREW_STEP.MANUAL }
            }] : []),
            ...disableByPass ? [] : [{
              label: 'Screw.Next',
              color: 'warning',
              action: screwStepActions.byPassSpecPoint()
            }]
          ],
          // eslint-disable-next-line camelcase
          title: `拧紧点失败：${n}`,
          content: `${this.failureMsg}`
        })
      );
      const { retry: doRetry, bypass, fail, manual } = yield race({
        bypass: take(SCREW_STEP.BYPASS_SPEC_POINT),
        retry: take(SCREW_STEP.BYPASS_RETRY),
        fail: take(SCREW_STEP.CONFIRM_FAIL_SPEC_POINT),
        manual: take(SCREW_STEP.MANUAL)
      });
      if (fail) {
        yield put(orderActions.stepStatus(this, STEP_STATUS.FAIL, { error: '拧紧失败' })); // 失败退出
      }
      if (bypass) {
        this._pointsManager.byPassControls(controls);
      }
      if (doRetry) {
        yield retry;
      }
      if (manual) {
        // yield take()]
        console.log('manual');
        return;
      }
    }
  } catch (e) {
    CommonLog.lError(e);
    throw e;
  }
}
