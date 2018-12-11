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
    showLayout:true
  },
  {
    name: 'working',
    url: '/working',
    title: 'main.operation',
    main: Working,
    icon: BuildIcon,
    color: '#00abbf',
    image: WorkingImg,
    enName: 'Operation',
    showLayout:true
  },
  {
    name: 'orders',
    url: '/orders',
    title: 'main.orders',
    main: WorkOrders,
    icon: CollectionsIcon,
    color: warningColor,
    image: editorImg,
    enName: 'Vehicle Queue',
    showLayout:true
  },
  {
    name: 'preferences',
    url: '/preferences',
    title: 'main.preferences',
    main: Preferences,
    icon: SettingsApplicationsIcon,
    color: '#fba53d',
    image: settingImg,
    enName: 'Preference',
    showLayout:true
  },
  {
    name: 'event',
    url: '/events',
    title: 'main.event',
    main: Event,
    icon: Mail,
    color: '#008adf',
    image: LoginImg,
    enName: 'Event Log',
    showLayout:true
  },
  {
    name: 'result',
    url: '/results',
    title: 'main.resultQuery',
    main: ConnResult,
    icon: Save,
    color: grayColor,
    image: LockingImg,
    enName: 'Result',
    showLayout:true
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
    enName: 'Help',
    showLayout:true
  },
  {
    name: 'login',
    url: '/pages/login',
    title: 'main.login',
    main: Pages,
    icon: Fingerprint,
    color: grayColor,
    image: LoginImg,
    enName: 'Login',
    showLayout:false
  }
];

export const routeConfigs = lodash.filter(
  routes,
  ele =>
    !(configs.operationSettings.opMode !== 'order' && ele.name === 'orders')
);

const indexRoutes = [
  { url: '/pages', main: Pages, showLayout:false },
  { url: '/', main: HomePage, showLayout:true }
];

export default indexRoutes;
