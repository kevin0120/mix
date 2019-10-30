import { expectSaga } from 'redux-saga-test-plan';
import ScrewStepMixin, { doPoint } from '../ScrewStep';
import Step from '../../Step';
import { demoOrder } from '../../../order/demoData';
import { STEP_STATUS } from '../../constants';
import { ORDER } from '../../../order/constants';
import { orderActions } from '../../../order/action';
import screwStepActions from '../action';
import { RESULT_STATUS } from '../constants';
import { mockGetTools } from './mocks/mocks';
import { demoPoints } from './mocks/demoPoints';

const stepData = (points) => {
  const demoStep = demoOrder.steps[1];

  if (!points) {
    return demoStep;
  }
  demoStep.payload.points = points;
  return demoStep;
};

const prepareDoingState = async (p) => {
  const step = new (ScrewStepMixin(Step))(stepData(p));
  await expectSaga(step._statusTasks[STEP_STATUS.ENTERING].bind(step), ORDER, orderActions).provide([
    mockGetTools
  ]).run();
  const { points } = demoOrder.steps[1].payload;
  const doingState = expectSaga(step._statusTasks[STEP_STATUS.DOING].bind(step), ORDER, orderActions);
  return { doingState, points, step };
};

describe('DoingState', () => {
  it('auto start with doing the first point', async () => {
    const { doingState, step } = await prepareDoingState();
    const point = step._pointsManager.points[0];
    return doingState
      .call([step, doPoint], [point], false, orderActions)
      .dispatch(screwStepActions.redoPoint(point))
      .run();
  });

  it('step finishes when all points are passed', async () => {
    const { doingState, points, step } = await prepareDoingState();

    points.forEach((p, idx) => {
      doingState.dispatch(screwStepActions.result({
        data: [{
          ...p,
          result: RESULT_STATUS.ok
        }]
      }));
    });

    return doingState.put(orderActions.stepStatus(step, STEP_STATUS.FINISHED)).run();


  });

  it('step fails when point fails', async () => {
    const { doingState, points, step } = await prepareDoingState();

    [RESULT_STATUS.nok, RESULT_STATUS.nok, RESULT_STATUS.nok].forEach(r =>
      doingState.dispatch(screwStepActions.result({
        data: [{
          ...points[0],
          result: r
        }]
      }))
    );
    return doingState.put(orderActions.stepStatus(step, STEP_STATUS.FAIL)).run();
  });

  it('call doPoint with the point passed on redo', async () => {
    const { doingState, step } = await prepareDoingState();
    step._pointsManager.points.forEach((p) => {
      doingState
        .call([step, doPoint], [p], false, orderActions)
        .dispatch(screwStepActions.redoPoint(p));
    });
    return doingState.run();
  });

  it.todo('should start next group when condition matches', async () => {
    const { doingState, step } = await prepareDoingState(demoPoints);
      // TODO
  });
});