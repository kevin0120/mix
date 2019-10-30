import * as matchers from 'redux-saga-test-plan/matchers';
import { expectSaga } from 'redux-saga-test-plan';
import ScrewStepMixin from '../ScrewStep';
import Step from '../../Step';
import { demoOrder } from '../../../order/demoData';
import { STEP_STATUS } from '../../constants';
import { ORDER } from '../../../order/constants';
import { orderActions } from '../../../order/action';
import { mockGetTools } from './mocks/mocks';


describe('Class ScrewStep Entering state', () => {
  it('goes to doing', () => {
    const step = new (ScrewStepMixin(Step))(demoOrder.steps[1]);

    return expectSaga(step._statusTasks[STEP_STATUS.ENTERING].bind(step), ORDER, orderActions)
      .provide([
        mockGetTools,
        [matchers.call.fn(step.updateData)]
      ])
      .put(orderActions.stepStatus(step, STEP_STATUS.DOING))
      .run();
  });

  it.todo('restores result data', () => {
    // TODO
  });
});