import { expectSaga, testSaga } from 'redux-saga-test-plan';
import { push } from 'connected-react-router';
import userAuth from '../saga';
import { configureStore } from '../../../store/configureStore';
import * as actions from '../action';
import config from '../../../shared/config';

const user: string = 'ming';
const { password } = config.defaultConfig.authorization.localUsers[user];

const userInfo = {
  uid: config.defaultConfig.authorization.localUsers[user].uid,
  name: config.defaultConfig.authorization.localUsers[user].name,
  uuid: config.defaultConfig.authorization.localUsers[user].uuid,
  avatar: config.defaultConfig.authorization.localUsers[user].avatar,
  role: config.defaultConfig.authorization.localUsers[user].role
};

describe('User Auth Saga Authorize Login Logout Workflow', () => {
  it('Local User Test Auth Login Workflow Success', () => expectSaga(userAuth)
      .withState(configureStore())
      .put(actions.loginSuccess(userInfo))
      .put(push('/app'))
      .dispatch(actions.loginRequest(user, password))
      .run());
});
