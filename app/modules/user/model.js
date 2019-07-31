export type tUuid = string;
export type tUserName = string;

export type tUser = {
  +name: string,
  +avatar: string,
  +uid: number,
  +uuid: tUuid,
  +role: string
};

export type tUserLoginAction = {
  +type: string,
  +userName: string,
  +password: string
};
