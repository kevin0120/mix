export type tNotifyKey = string;

export type tNotifyVariant = 'info' | 'error' | 'info';

export type tNotifyState = {
  notifications: Array<tNotification>
};

export type tNotification = {
  message: string,
  options: {
    key: tNotifyKey,
    variant: tNotifyVariant,
    anchorOrigin: {
      vertical: 'top',
      horizontal: 'left'
    }
  }
}