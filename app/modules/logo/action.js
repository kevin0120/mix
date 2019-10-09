// @flow

import type { tActionLogo } from './type';

const LOGO = {
  FETCH_OK: 'LOGO_FETCH_OK'
};

export function fetchLogoOK(logo: string): tActionLogo {
  return {
    type: LOGO.FETCH_OK,
    logo
  };
}

export default LOGO;
