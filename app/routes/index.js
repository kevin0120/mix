/* eslint flowtype-errors/show-errors: 0 */
import React from 'react';

// import Editor from '../containers/editor';
// import Help from '../containers/help';
import Preferences from '../containers/config';
import Event from '../containers/event';
import ConnResult from '../containers/result';
import Working from '../containers/working';

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
import Pages from '../layouts/Pages';
import HomePage from '../containers/Home';

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

export const routeConfigs = [
  {
    name: 'welcome',
    url: '/welcome',
    title: 'main.home',
    main: HomePage,
    icon: HomeIcon,
    color: grayColor
  },
  {
    name: 'working',
    url: '/working',
    title: 'main.operation',
    main: Working,
    icon: BuildIcon,
    color: grayColor,
    image: WorkingImg,
    enName: 'Operating'
  },
  // {
  //   name: 'editor',
  //   url: '/editor',
  //   title: 'main.worksheet_img',
  //   main: Editor,
  //   icon: CollectionsIcon,
  //   color: grayColor,
  //   image: editorImg,
  //   enName: 'Operation Instructions',
  // },
  {
    name: 'preferences',
    url: '/preferences',
    title: 'main.preferences',
    main: Preferences,
    icon: SettingsApplicationsIcon,
    color: grayColor,
    image: settingImg,
    enName: 'Preferences'
  },
  {
    name: 'event',
    url: '/events',
    title: 'main.event',
    main: Event,
    icon: Mail,
    color: grayColor,
    image: LoginImg,
    enName: 'Event Log'
  },
  {
    name: 'result',
    url: '/results',
    title: 'main.resultQuery',
    main: ConnResult,
    icon: Save,
    color: grayColor,
    image: LoginImg,
    enName: 'Query Results'
  },
  // {
  //   name: 'lock',
  //   url: '/pages/lock-screen-page',
  //   title: 'main.lock',
  //   main: Pages,
  //   icon: LockIcon,
  //   color: grayColor,
  //   image: LockingImg,
  //   enName: 'Lock',
  // },
  // {
  //   name: 'help',
  //   url: '/help',
  //   title: 'main.help',
  //   main: Help,
  //   icon: HelpIcon,
  //   color: grayColor,
  //   image: helpImg,
  //   enName: 'Help',
  // },
  {
    name: 'login',
    url: '/pages/login',
    title: 'main.login',
    main: Pages,
    icon: Fingerprint,
    color: grayColor,
    image: LoginImg,
    enName: 'Login'
  }
];

const indexRoutes = [
  { path: '/pages', component: Pages },
  { path: '/', component: HomePage }
];

export default indexRoutes;
