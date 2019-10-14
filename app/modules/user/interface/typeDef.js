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

export type tAuthRespData = {
  +id: number,
  +name: string,
  +uuid: string,
  +image_small: string
};

export type tAuthInfo = {
  +name: string | null,
  +password: string | null,
  +method: string
};

export type tAuthLogout = {
  +uuid: string | null
};
