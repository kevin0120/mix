/* eslint flowtype-errors/show-errors: 0 */

// icons
import HomeIcon from '@material-ui/icons/Home';
import SettingsApplicationsIcon from '@material-ui/icons/SettingsApplications';
import BuildIcon from '@material-ui/icons/Build';
import HelpIcon from '@material-ui/icons/Help';
import CollectionsIcon from '@material-ui/icons/Collections';
import LockIcon from '@material-ui/icons/Lock';
import Fingerprint from '@material-ui/icons/Fingerprint';
import Trend from '@material-ui/icons/TrendingUp';
import Mail from '@material-ui/icons/Mail';
import Save from '@material-ui/icons/Save';
import ViewerIcon from '@material-ui/icons/Image';

// color
import cyan from '@material-ui/core/colors/cyan';
import indigo from '@material-ui/core/colors/indigo';
import blue from '@material-ui/core/colors/blue';
import pink from '@material-ui/core/colors/pink';
import teal from '@material-ui/core/colors/teal';
import orange from '@material-ui/core/colors/orange';
import lightGreen from '@material-ui/core/colors/lightGreen';

import WorkOrders from './orders';
import Working from './working';
import ConnResult from './result';
import Event from './event';
import Preferences from './config';
import Help from './help';
import Pages from './layouts/Pages';
import HomePage from './home';
import Curve from './curve';
import Viewer from './viewer';
import Layout from './Layout/layout';
import page from './layouts/Pages';

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

export const pages = {
  '/app': {
    component: Layout,
    // '/welcome': {
    //   title: 'main.home',
    //   component: HomePage,
    //   icon: HomeIcon,
    //   color: indigo[shade],
    //   exact: true
    // },
    '/working': {
      title: 'main.operation',
      component: Working,
      icon: BuildIcon,
      color: cyan[shade],
      image: WorkingImg,
      exact: true
    },
    '/viewer': {
      title: 'main.operationViewer',
      component: Viewer,
      icon: ViewerIcon,
      color: lightGreen[shade],
      image: viewerImg,
      exact: true

    },
    '/order': {
      title: 'main.orders',
      component: WorkOrders,
      icon: CollectionsIcon,
      color: warningColor,
      image: editorImg,
      exact: true

    },
    '/preference': {
      title: 'main.preferences',
      component: Preferences,
      icon: SettingsApplicationsIcon,
      color: orange[shade],
      image: settingImg,
      exact: true

    },
    '/event': {
      title: 'main.event',
      component: Event,
      icon: Mail,
      color: blue[shade],
      image: LoginImg,
      exact: true

    },
    '/result': {
      title: 'main.resultQuery',
      component: ConnResult,
      icon: Save,
      color: grayColor,
      image: LockingImg,
      exact: true

    },
    '/curve': {
      title: 'main.curve',
      component: Curve,
      icon: Trend,
      color: teal[shade],
      image: CurveImg,
      exact: true

    },
    '/help': {
      title: 'main.help',
      component: Help,
      icon: HelpIcon,
      color: pink[shade],
      image: helpImg,
      exact: true

    }
  },
  '/pages': {
    component: page,
    '/login': {
      title: 'main.login',
      component: Pages,
      icon: Fingerprint,
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



function flattenRoutes(route) {
  console.log('flattenRoutes');
  const subRouteUrls = Object.keys(route).filter((key) => /^\//.test(key));
  const subRouteList = subRouteUrls.map(u => {
    return {
      ...route[u],
      url: (route.url || '') + u,
      name: u.slice(1),
    };
  });
  let flatRoutes = subRouteList;
  subRouteList.forEach((subRoute) => {
    flatRoutes = flatRoutes.concat(flattenRoutes(subRoute));
  });
  return flatRoutes;
}



export default function filterRoutesByConfig(config) {

  return config.map((c) => {
    const p = pages.find((p) => p.name === c.name);
    return {
      ...p,
      roles: c.roles
    };
  });
}
