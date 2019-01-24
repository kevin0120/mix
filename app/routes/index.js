/* eslint flowtype-errors/show-errors: 0 */
import React from 'react';

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

import WorkOrders from '../containers/orders';
import Working from '../containers/working';
import ConnResult from '../containers/result';
import Event from '../containers/event';
import Preferences from '../containers/config';
import Help from '../containers/help';
import Pages from '../layouts/Pages';
import HomePage from '../containers/home';
import Curve from '../containers/curve';
import Viewer from '../containers/viewer';

// imgs
import helpImg from '../../resources/imgs/help.png';
import editorImg from '../../resources/imgs/operaIns.png';
import WorkingImg from '../../resources/imgs/operation.png';
import settingImg from '../../resources/imgs/setting.png';
import LockingImg from '../../resources/imgs/lock.jpeg';
import LoginImg from '../../resources/imgs/login.jpeg';

import {
  grayColor,
  roseColor,
  primaryColor,
  infoColor,
  successColor,
  warningColor,
  dangerColor
} from '../common/jss/material-react-pro';

const lodash = require('lodash');

import configs from '../shared/config/index';

const routes = [
  {
    name: 'welcome',
    url: '/welcome',
    title: 'main.home',
    main: HomePage,
    icon: HomeIcon,
    color: '#00abf3',
    showLayout: true,
    roles: ['user', 'admin']
  },
  {
    name: 'working',
    url: '/working',
    title: 'main.operation',
    main: Working,
    icon: BuildIcon,
    color: '#00abbf',
    image: WorkingImg,
    showLayout: true,
    roles: ['user', 'admin']
  },
  {
    name: 'viewer',
    url: '/viewer',
    title: 'main.operationViewer',
    main: Viewer,
    icon: BuildIcon,
    color: '#00abbf',
    image: WorkingImg,
    showLayout: true,
    roles: ['user', 'admin']
  },
  {
    name: 'orders',
    url: '/orders',
    title: 'main.orders',
    main: WorkOrders,
    icon: CollectionsIcon,
    color: warningColor,
    image: editorImg,
    showLayout: true,
    roles: ['user', 'admin']
  },
  {
    name: 'preferences',
    url: '/preferences',
    title: 'main.preferences',
    main: Preferences,
    icon: SettingsApplicationsIcon,
    color: '#fba53d',
    image: settingImg,
    showLayout: true,
    roles: ['admin']
  },
  {
    name: 'event',
    url: '/events',
    title: 'main.event',
    main: Event,
    icon: Mail,
    color: '#008adf',
    image: LoginImg,
    showLayout: true,
    roles: ['user', 'admin']
  },
  {
    name: 'result',
    url: '/results',
    title: 'main.resultQuery',
    main: ConnResult,
    icon: Save,
    color: grayColor,
    image: LockingImg,
    showLayout: true,
    roles: ['user', 'admin']
  },
  {
    name: 'curve',
    url: '/curves',
    title: 'main.curve',
    main: Curve,
    icon: Trend,
    color: '#009688',
    image: LockingImg,
    showLayout: true,
    roles: ['user', 'admin']
  },
  // {
  //   name: 'lock',
  //   url: '/pages/lock-screen-page',
  //   title: 'main.lock',
  //   main: Pages,
  //   icon: LockIcon,
  //   color: grayColor,
  //   image: LockingImg,
  //   enName: 'Lock'
  // },
  {
    name: 'help',
    url: '/help',
    title: 'main.help',
    main: Help,
    icon: HelpIcon,
    color: '#8a6c63',
    image: helpImg,
    showLayout: true,
    roles: ['user', 'admin']
  },
  {
    name: 'login',
    url: '/pages/login',
    title: 'main.login',
    main: Pages,
    icon: Fingerprint,
    color: grayColor,
    image: LoginImg,
    showLayout: false,
    roles: ['user', 'admin']
  }
];

export const routeConfigs = lodash.filter(
  routes,
  ele =>
    !(configs.operationSettings.opMode !== 'order' && ele.name === 'orders')
);

const indexRoutes = [
  { url: '/pages', main: Pages, showLayout: false },
  { url: '/', main: HomePage, showLayout: true }
];

export default indexRoutes;
