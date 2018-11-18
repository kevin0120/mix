import Fingerprint from '@material-ui/icons/Fingerprint';
import LockOpen from '@material-ui/icons/LockOpen';
import LockScreenPage from '../containers/lock';
import LoginPage from '../containers/login';

// @material-ui/icons

const pagesRoutes = [
  {
    path: '/pages/lock-screen-page',
    name: 'Lock Screen Page',
    short: 'Lock',
    mini: 'LSP',
    icon: LockOpen,
    component: LockScreenPage
  },
  {
    path: '/pages/login',
    name: 'Login Page',
    short: 'Login',
    mini: 'LP',
    icon: Fingerprint,
    component: LoginPage
  },
  {
    redirect: true,
    path: '/pages',
    pathTo: '/pages/lock-screen-page',
    name: 'Lock Screen Page'
  }
];

export default pagesRoutes;

// WEBPACK FOOTER //
// ./src/routes/pages.jsx
