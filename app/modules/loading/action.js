export const LOADING = {
  START: 'LOADING_START',
  STOP: 'LOADING_STOP'
};

export default {
  start: () => ({
    type: LOADING.START
  }),
  stop: () => ({
    type: LOADING.STOP
  })
};
