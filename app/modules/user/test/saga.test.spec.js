import { expectSaga } from 'redux-saga-test-plan';
import * as matchers from 'redux-saga-test-plan/matchers';
import { push } from 'connected-react-router';
import userAuth from '../saga';
import { configureStore } from '../../../store/configureStore';
import * as actions from '../action';
import configs from '../../../shared/config';
import { getUserInfo } from '../../../api/user';
import type { tAuthRespData } from '../interface/typeDef';
import status from 'http-status';

const user: string = 'ming';
const { password } = configs.authorization.localUsers[user];

const state = configureStore().getState();

const LocalUserInfo = {
  uid: configs.authorization.localUsers[user].uid,
  name: user,
  uuid: configs.authorization.localUsers[user].uuid,
  avatar: configs.authorization.localUsers[user].avatar,
  role: configs.authorization.localUsers[user].role
};

const respAckUser: tAuthRespData = {
  id: 1,
  name: user,
  uuid: 'uuid',
  image_small: 'avatar',
};

const RemoteUserInfo = {
  uid: respAckUser.id,
  name: respAckUser.name,
  uuid: respAckUser.uuid,
  avatar: respAckUser.image_small,
  role: "admin"
};

describe('User Auth Saga Authorize Login Logout Workflow', () => {
  it('Local User Test Auth Login Workflow Success', () =>
    expectSaga(userAuth)
      .withState(state)
      .put(actions.loginSuccess(LocalUserInfo))
      .put(push('/app'))
      .dispatch(actions.loginRequest(user, password))
      .run(3000));
  it('Online(Remote) User Test Auth Login Workflow Success', () =>
    expectSaga(userAuth)
      .withState(state)
      .provide([
        // Use the `call.fn` matcher from Redux Saga Test Plan
        [matchers.call.fn(getUserInfo), {status: status.OK, data: respAckUser }],
      ])
      .put(actions.loginSuccess(RemoteUserInfo))
      .put(push('/app'))
      .dispatch(actions.loginRequest(user, password, 'online'))
      .run(3000));
});
