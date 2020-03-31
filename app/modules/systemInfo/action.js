export const SYSTEM_INFO = {
  SET_WORKCENTER: 'SYSTEM_SET_WORKCENTER'
};

const setWorkcenter = (workcenter) => ({
  type: SYSTEM_INFO.SET_WORKCENTER,
  workcenter
});

export default {
  setWorkcenter
}
