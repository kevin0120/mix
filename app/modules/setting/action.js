export const USER_CONFIGS = {
  PRE_SAVE: 'USER_CONFIGS_PRE_SAVE',
  SAVE: 'USER_CONFIGS_SAVE',
  SET_UUID: 'USER_CONFIGS_SET_UUID'
};

export default function saveConfigs(section, newConfigs) {
  return {
    type: USER_CONFIGS.PRE_SAVE,
    data: {
      section,
      newConfigs
    }
  };
}
