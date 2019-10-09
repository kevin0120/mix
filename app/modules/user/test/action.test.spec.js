import * as actions from '../action';

const user = 'test';

const uuid = '05e1ddca-dfd2-4ae0-9f7c-451f0bf94f19';

describe('User Auth Actions', () => {
  it('Local User Login Request action', () => {
    expect(actions.loginRequest(user)).toMatchSnapshot();
  });

  it('Online User Login Request action', () => {
    expect(actions.loginRequest(user, '', 'online')).toMatchSnapshot();
  });

  it('Local User UUID Login Request action', () => {
    expect(actions.loginRequestUuid(uuid)).toMatchSnapshot();
  });
});
