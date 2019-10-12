import faker from 'faker';
import reducer from '../reducer';
import { USER } from '../action';
import type { tUser } from '../interface/typeDef';

describe('user Reducer', () => {
  it('User should Return the initial state', () => {
    expect(reducer(undefined, {})).toMatchSnapshot();
  });
  it('User Action UnSupport Action Type', () => {
    expect(reducer(undefined, { type: 'any' })).toMatchSnapshot();
  });
  it('User Action Login Success', () => {
    const data: tUser = {
      uid: 1,
      name: faker.name.findName(),
      uuid: faker.random.uuid(),
      avatar: faker.image.avatar(),
      role: 'admin'
    };
    expect(
      reducer(undefined, { type: USER.LOGIN.SUCCESS }, data)
    ).toMatchSnapshot();
  });
});
