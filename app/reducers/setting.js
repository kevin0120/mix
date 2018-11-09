// @flow
import type { Action } from './types';

import { defaultConfigs } from '../shared/config/defaultConfig'


export default function setting(state: object = defaultConfigs, action: Action) {
  switch (action.type) {
    default:
      return state;
  }
}
