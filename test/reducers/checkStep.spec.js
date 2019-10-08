import counter from '../../app/reducers/counter';
import {
  CHECK_STEP,
} from '../../app/modules/step/checkStep/action';

describe('reducers', () => {
  describe('counter', () => {
    it('should handle initial state', () => {
      expect(counter(undefined, {})).toMatchSnapshot();
    });

    it('should handle CHECK_STEP.SUBMIT', () => {
      expect(counter(1, { type: CHECK_STEP.SUBMIT })).toMatchSnapshot();
    });

    it('should handle CHECK_STEP.CANCEL', () => {
      expect(counter(1, { type: CHECK_STEP.CANCEL })).toMatchSnapshot();
    });

    it('should handle unknown action type', () => {
      expect(counter(1, { type: 'unknown' })).toMatchSnapshot();
    });
  });
});
