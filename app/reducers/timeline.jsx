import { TIMELINE_STORY } from '../actions/actionTypes';

const defaultTimeline = [
  // {
  //   // First story
  //   inverted: true,
  //   badgeColor: "danger",
  //   badgeIcon: CardTravel,
  //   title: "Some Title",
  //   titleColor: "danger",
  //   body: (
  //     <p>
  //       Wifey made the best Father's Day meal ever. So thankful so happy so
  //       blessed. Thank you for making my family We just had fun with the
  //       “future” theme !!! It was a fun night all together ... The always rude
  //       Kanye Show at 2am Sold Out Famous viewing @ Figueroa and 12th in
  //       downtown.
  //     </p>
  //   ),
  //   footerTitle: "11 hours ago via Twitter"
  // },
];

type actionType = {
  +type: string,
  +story: object
};

export default function timeline(
  state: object = defaultTimeline,
  action: actionType
) {
  switch (action.type) {
    case TIMELINE_STORY.NEW:
      return NewStory(state, action.story);
    case TIMELINE_STORY.CLEAR:
      return ClearStory();
    default:
      return state;
  }
}

export function NewStory(state, story) {
  return [story, ...state];
}

export function ClearStory() {
  return [];
}
