/* eslint flowtype-errors/show-errors: 0 */

// color
import * as colors from '@material-ui/core/colors';

// icons
import * as icons from '@material-ui/icons';

import WorkOrders from './orders';
import Working from './working';
import StepWorking from './stepWorking';
import ConnResult from './result';
import Event from './event';
import Preferences from './config';
import Help from './help';
import LockLayout from './lockLayout';
import LoginPage from './login';
import Curve from './curve';
import Viewer from './viewer';
import Layout from './appLayout';
import WorkingTemplate from './workingTemplate';
import OperationList from './OperationList';
import HomePage from './home';

// imgs
import helpImg from '../../resources/imgs/help.png';
import viewerImg from '../../resources/imgs/craft.jpeg';
import editorImg from '../../resources/imgs/operaIns.png';
import WorkingImg from '../../resources/imgs/operation.png';
import settingImg from '../../resources/imgs/setting.png';
import LockingImg from '../../resources/imgs/lock.jpeg';
import CurveImg from '../../resources/imgs/curveImg.jpeg';
import LoginImg from '../../resources/imgs/login.jpeg';

import {
  grayColor,
  warningColor
} from '../common/jss/material-react-pro';

const shade = 500;

const pages = {
  '/app': {
    component: Layout,
    navBarContents: ['clock', 'pages', 'language'], //'avatar','clock', 'pages', 'language', 'sysInfo', 'healthCheck'
    DefaultContent: OperationList,
    title: 'main.home',
    icon: icons.Home,
    color: colors.indigo[shade],
    '/working': {
      title: 'main.operation',
      component: WorkingTemplate,
      icon: icons.Build,
      color: colors.cyan[shade],
      image: WorkingImg,
      exact: true
    },
    '/viewer': {
      title: 'main.operationViewer',
      component: Viewer,
      icon: icons.Image,
      color: colors.lightGreen[shade],
      image: viewerImg,
      exact: true

    },
    '/order': {
      title: 'main.orders',
      component: WorkOrders,
      icon: icons.Collections,
      color: warningColor,
      image: editorImg,
      exact: true

    },
    '/preference': {
      title: 'main.preferences',
      component: Preferences,
      icon: icons.SettingsApplications,
      color: colors.orange[shade],
      image: settingImg,
      exact: true

    },
    '/event': {
      title: 'main.event',
      component: Event,
      icon: icons.Mail,
      color: colors.blue[shade],
      image: LoginImg,
      exact: true

    },
    '/result': {
      title: 'main.resultQuery',
      component: ConnResult,
      icon: icons.Save,
      color: grayColor,
      image: LockingImg,
      exact: true

    },
    '/curve': {
      title: 'main.curve',
      component: Curve,
      icon: icons.TrendingUp,
      color: colors.teal[shade],
      image: CurveImg,
      exact: true

    },
    '/help': {
      title: 'main.help',
      component: Help,
      icon: icons.Help,
      color: colors.pink[shade],
      image: helpImg,
      exact: true

    }
  },
  '/pages': {
    component: LockLayout,
    '/login': {
      title: 'main.login',
      component: LoginPage,
      icon: icons.Fingerprint,
      color: grayColor,
      image: LoginImg,
      exact: true
    }
    // lock: {
    //   url: '/pages/lock-screen-page',
    //   title: 'main.lock',
    //   component: Pages,
    //   icon: LockIcon,
    //   color: grayColor,
    //   image: LockingImg,
    //   enName: 'Lock'
    // },
  }
};

export const getContentByUrl = (url) => {
  const arr = url.split('/').filter((u) => u !== '');
  const page = arr.reduce((p, r) => (p && p[`/${r}`]) || null, pages);
  return page && {
    ...page,
    url,
    name: url.slice(-1)[0]
  };
};


export default pages;
