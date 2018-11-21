import { USER_CONFIGS } from './actionTypes';

export default function saveConfigs(section, newConfigs) {
  return {
    type: USER_CONFIGS.PRE_SAVE,
    data: {
      section,
      newConfigs,
    },
  };
}
