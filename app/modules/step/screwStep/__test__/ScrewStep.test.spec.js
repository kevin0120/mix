import { expectSaga } from 'redux-saga-test-plan';
import * as matchers from 'redux-saga-test-plan/matchers';
import ScrewStepMixin, {doPoint} from '../ScrewStep';
import Step from '../../Step';
import { STEP_STATUS } from '../../constants';
import { ORDER } from '../../../order/constants';
import { orderActions } from '../../../order/action';
import { demoOrder } from '../../../order/demoData';
import screwStepActions from '../action';
import { controllerModes, RESULT_STATUS } from '../constants';
import controllerModeTasks from '../controllerModeTasks';

const demoStepData = {
  id: '0001',
  name: 'name',
  desc: 'this is a screw step',
  info: {},
  type: 'screw'
};

const prepareDoingState = async () => {
  const step = new (ScrewStepMixin(Step))(demoOrder.steps[1]);
  await expectSaga(step._statusTasks[STEP_STATUS.ENTERING].bind(step), ORDER, orderActions).run();
  const demoPoints = demoOrder.steps[1].payload.points;
  const doingState = expectSaga(step._statusTasks[STEP_STATUS.DOING].bind(step), ORDER, orderActions);
  return { doingState, demoPoints, step };
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
    const { doingState, demoPoints, step } = await prepareDoingState();

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
    const { doingState, demoPoints, step } = await prepareDoingState();

    [RESULT_STATUS.nok, RESULT_STATUS.nok, RESULT_STATUS.nok].forEach(r =>
      doingState.dispatch(screwStepActions.result({
        data: [{
          ...demoPoints[0],
          result: r
        }]
      }))
    );
    return doingState.put(orderActions.stepStatus(step, STEP_STATUS.FAIL)).run();
  });

  it('Doing state: should do the passed point', async () => {
    const { doingState, demoPoints, step } = await prepareDoingState();
    const point = null;
    doingState.provide({
      call(effect, next) {
        console.log(effect.fn, doPoint);
        if (effect.fn === doPoint) {
          console.log('do point');
          // TODO
        }
        return next();
      }
    });
    return doingState.dispatch(screwStepActions.redoPoint(point)).run();
  });


});
