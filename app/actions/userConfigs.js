import { USER_CONFIGS } from './actionTypes';

export default function saveConfigs(section, newConfigs) {
  return {
    type: USER_CONFIGS.SAVE,
    data: {
      section,
      newConfigs,
    },
  };
}
