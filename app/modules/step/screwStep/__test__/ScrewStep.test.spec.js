import { demoOrder } from '../../../order/demoData';
import ScrewStepMixin from '../ScrewStep';
import Step from '../../Step';

const demoStepListData = {
  code: '0001',
  name: 'name',
  desc: 'this is a screw step',
  info: {},
  type: 'screw'
};

describe('Class ScrewStep', () => {

  const step = new (ScrewStepMixin(Step))(demoStepListData);

  it('should be initialized successfully', () => {
    expect(step).toMatchSnapshot();
  });

  it('should be updated successfully', () => {
    expect(() => {
      step.update(demoOrder.steps[1]);
    }).not.toThrow();
    expect(step).toMatchSnapshot();
  });

});
