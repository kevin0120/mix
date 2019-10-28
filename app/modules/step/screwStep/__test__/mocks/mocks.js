import * as matchers from 'redux-saga-test-plan/matchers';
import { getTools } from '../../ScrewStep';

export const mockGetTools = [matchers.call.fn(getTools), []];