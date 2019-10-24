import { expectSaga } from 'redux-saga-test-plan';
import * as matchers from 'redux-saga-test-plan/matchers';
import ScrewStepMixin from '../ScrewStep';
import Step from '../../Step';
import { STEP_STATUS } from '../../constants';
import { ORDER } from '../../../order/constants';
import { orderActions } from '../../../order/action';
import { demoOrder } from '../../../order/demoData';
import screwStepActions from '../action';
import { RESULT_STATUS } from '../constants';

const demoStepData = {
  id: '0001',
  name: 'name',
  desc: 'this is a screw step',
  info: {},
  type: 'screw'
};


describe('Class ScrewStep', () => {
  // it('should be initialized successfully', () => {
  //   expect(step).toMatchSnapshot();
  // });
  //
  // it('should be updated successfully', () => {
  //   expect(() => {
  //     step.update(demoOrder.steps[1]);
  //   }).not.toThrow();
  // });

  it('Entering state', () => {
    const step = new (ScrewStepMixin(Step))(demoOrder.steps[1]);

    return expectSaga(step._statusTasks[STEP_STATUS.ENTERING].bind(step), ORDER, orderActions)
      .provide([
        // Use the `call.fn` matcher from Redux Saga Test Plan
        [matchers.call.fn(step.updateData)]
      ])
      .put(orderActions.stepStatus(step, STEP_STATUS.DOING))
      .run();
  });

  it('Doing state: step finishes when all points are passed', async () => {
    const step = new (ScrewStepMixin(Step))(demoOrder.steps[1]);
    await expectSaga(step._statusTasks[STEP_STATUS.ENTERING].bind(step), ORDER, orderActions).run();
    const demoPoints = demoOrder.steps[1].payload.points;
    const doingState = expectSaga(step._statusTasks[STEP_STATUS.DOING].bind(step), ORDER, orderActions);

    demoPoints.forEach((p, idx) => {
      doingState.dispatch(screwStepActions.result({
        data: [{
          ...p,
          result: RESULT_STATUS.ok
        }]
      }));
    });

    return doingState.put(orderActions.stepStatus(step, STEP_STATUS.FINISHED)).run();


  });

  it('Doing state: step fails when point fails', async () => {
    const step = new (ScrewStepMixin(Step))(demoOrder.steps[1]);
    await expectSaga(step._statusTasks[STEP_STATUS.ENTERING].bind(step), ORDER, orderActions).run();
    const demoPoints = demoOrder.steps[1].payload.points;
    const doingState = expectSaga(step._statusTasks[STEP_STATUS.DOING].bind(step), ORDER, orderActions);

    [RESULT_STATUS.nok, RESULT_STATUS.nok, RESULT_STATUS.nok].forEach(r =>
      doingState.delay(10).dispatch(screwStepActions.result({
        data: [{
          ...demoPoints[0],
          result: r
        }]
      }))
    );
    return doingState.put(orderActions.stepStatus(step, STEP_STATUS.FAIL)).run();
  });


});
