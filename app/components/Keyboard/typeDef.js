// @flow
export type tKeyboardConfig = {
  // eslint-disable-next-line flowtype/no-weak-types
  onSubmit: ?(text: string)=>void,
  title: string,
  label: string,
  text: string,
  required?: boolean
};


export type tKeyboardBinder = (tKeyboardConfig) => void;

export type tKeyboard = { keyboard: { setInput: Function } };
