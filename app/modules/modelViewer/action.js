// @flow
export const MODEL_VIEWER = {
  OPEN: 'MODEL_VIEWER_OPEN'
};

export default {
  open: (url: string) => ({
    type: MODEL_VIEWER.OPEN,
    url
  })
};
