import reducer from '../reducer';

describe('user Reducer', () => {
  it('User should Return the initial state', () => {
    expect(reducer(undefined, {})).toMatchSnapshot();
  });
});
