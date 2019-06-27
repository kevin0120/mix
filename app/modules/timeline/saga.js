import React from 'react';
import CheckCircle from '@material-ui/icons/CheckCircle';
import Close from '@material-ui/icons/Close';
import InfoRounded from '@material-ui/icons/InfoRounded';

import { put } from 'redux-saga/effects';

import { TIMELINE_STORY } from './action';
import {STORY_TYPE} from './model';

const dayjs = require('dayjs');



export function* addNewStory(level, title, msg) {
  try{
    const story = createNewStory(level, title, msg);
    yield put({ type: TIMELINE_STORY.NEW, story });
  }catch (e) {
    console.error(e);
  }
}

export function* clearStories() {
  try {
    yield put({ type: TIMELINE_STORY.CLEAR });
  }catch (e) {
    console.error(e);
  }
}

function createNewStory(level, title, msg) {
  let icon = null;
  let badgeColor = null;
  switch (level) {
    case STORY_TYPE.PASS:
      icon = CheckCircle;
      badgeColor = 'success';
      break;
    case STORY_TYPE.FAIL:
      icon = Close;
      badgeColor = 'warning';
      break;
    case STORY_TYPE.INFO:
      icon = InfoRounded;
      badgeColor = 'info';
      break;
    default:
      icon = InfoRounded;
      badgeColor = 'info';
      break;
  }

  return {
    inverted: true,
    badgeColor,
    badgeIcon: icon,
    title,
    timestamp: dayjs().format('YYYY-MM-DD HH:mm:ss'),
    titleColor: badgeColor,
    body: <p>{msg}</p>
  };
}
