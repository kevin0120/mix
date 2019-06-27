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
  welcome: {
    url: '/app/welcome',
    title: 'main.home',
    main: HomePage,
    icon: HomeIcon,
    color: indigo[shade],
  },
  working: {
    url: '/app/working',
    title: 'main.operation',
    main: Working,
    icon: BuildIcon,
    color: cyan[shade],
    image: WorkingImg,
  },
  viewer: {
    url: '/app/viewer',
    title: 'main.operationViewer',
    main: Viewer,
    icon: ViewerIcon,
    color: lightGreen[shade],
    image: viewerImg,
  },
  orders: {
    url: '/app/orders',
    title: 'main.orders',
    main: WorkOrders,
    icon: CollectionsIcon,
    color: warningColor,
    image: editorImg,
  },
  preferences: {
    url: '/app/preferences',
    title: 'main.preferences',
    main: Preferences,
    icon: SettingsApplicationsIcon,
    color: orange[shade],
    image: settingImg,
  },
  events: {
    url: '/app/events',
    title: 'main.event',
    main: Event,
    icon: Mail,
    color: blue[shade],
    image: LoginImg,
  },
  result: {
    url: '/app/results',
    title: 'main.resultQuery',
    main: ConnResult,
    icon: Save,
    color: grayColor,
    image: LockingImg,
  },
  curve: {
    url: '/app/curves',
    title: 'main.curve',
    main: Curve,
    icon: Trend,
    color: teal[shade],
    image: CurveImg,
  },
  // lock: {
  //   url: '/pages/lock-screen-page',
  //   title: 'main.lock',
  //   main: Pages,
  //   icon: LockIcon,
  //   color: grayColor,
  //   image: LockingImg,
  //   enName: 'Lock'
  // },
  help: {
    url: '/app/help',
    title: 'main.help',
    main: Help,
    icon: HelpIcon,
    color: pink[shade],
    image: helpImg,
  },
  login: {
    url: '/pages/login',
    title: 'main.login',
    main: Pages,
    icon: Fingerprint,
    color: grayColor,
    image: LoginImg,
  }
};

export default function filterRoutesByConfig(config) {
  const routeNames = Object.keys(pages);
  return config.map((c) => {
    const name= routeNames.find((n)=>n===c.name);
    return {
      ...pages[name],
      roles:c.roles
    };
  });
}
