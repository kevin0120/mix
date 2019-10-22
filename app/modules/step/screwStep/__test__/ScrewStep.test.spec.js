import { expectSaga } from 'redux-saga-test-plan';
import ScrewStepMixin from '../ScrewStep';
import Step from '../../Step';
import STEP_STATUS from '../../constants';
import { ORDER } from '../../../order/constants';
import { orderActions } from '../../../order/action';
import dataImg from '../../../../../resources/imgs/workingImg.png';

const demoStepData = {
  id: '0001',
  name: 'name',
  desc: 'this is a screw step',
  info: {},
  type: 'screw'
};

const step = new (ScrewStepMixin(Step))(demoStepData);


describe('Class ScrewStep', () => {
  it('should be initialized successfully', () => {
    expect(step).toMatchSnapshot();
  });

  it('should be updated successfully', () => {
    expect(() => {
      step.update({
        payload: {
          controllerMode: 'pset',
          points: [
            {
              sequence: 1,
              group_sequence: 1,
              x: 10,
              y: 10,
              maxRetryTimes: 3,
              pset: 1,
              toolSN: 'demo'
            },
            {
              sequence: 2,
              group_sequence: 2,
              x: 20,
              y: 20,
              maxRetryTimes: 3,
              pset: 1,
              toolSN: 'demo'
            },
            {
              sequence: 3,
              group_sequence: 3,

              x: 30,
              y: 30,
              maxRetryTimes: 3,
              pset: 1,
              toolSN: 'demo'
            },
            {
              sequence: 4,
              group_sequence: 3,
              x: 40,
              y: 40,
              maxRetryTimes: 3,
              pset: 1,
              toolSN: 'demo'
            },
            {
              sequence: 5,
              group_sequence: 4,
              x: 50,
              y: 50,
              maxRetryTimes: 3,
              pset: 1,
              toolSN: 'demo'
            }
          ],
          image: dataImg
        }
      });
    }).not.toThrow();
  });

  it('Entering state ', () => {
    expectSaga(step._statusTasks[STEP_STATUS.ENTERING].bind(step), ORDER, orderActions)

      .run();
  });


});
