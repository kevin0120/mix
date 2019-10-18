// @flow
import type { tHealthzStatus } from './typeDef';

export const HEALTHZ = {
  UPDATE: 'HEALTHZ_UPDATE',
  DATA: 'HEALTHZ_DARA'
};

export default {
  data: (status: tHealthzStatus) => ({
    type: HEALTHZ.DATA,
    status
  }),
  update: () => ({
    type: HEALTHZ.UPDATE
  })
};
