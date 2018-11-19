/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

// @flow

const defaultUsers = [
  // {
  //   name: "dummy",
  //   avatar: "",
  //   uid: 1,
  // },
];

type actionType = {
  +type: string,
  +name: string,
  +avatar: string,
  +uid: number
};

export default function users(
  state: object = defaultUsers,
  action: actionType
) {
  switch (action.type) {
    default:
      return state;
  }
}
