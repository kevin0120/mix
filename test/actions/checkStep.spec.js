import * as actions from '../../app/modules/step/checkStep/action';

describe('Check Step actions', () => {
  it('Check Step Submit should create check step submit action', () => {
    expect(actions.default.submit()).toMatchSnapshot();
  });

  it('Check Step Submit should cancel check step cancel action', () => {
    expect(actions.default.cancel()).toMatchSnapshot();
  });
});
