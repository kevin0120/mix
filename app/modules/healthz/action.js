export const HEALTHZ = {
  UPDATE: 'HEALTHZ_UPDATE',
  DATA: 'HEALTHZ_DARA'
};

export default {
  data: status => ({
    type: HEALTHZ.DATA,
    status
  }),
  update: () => ({
    type: HEALTHZ.UPDATE
  })
};
