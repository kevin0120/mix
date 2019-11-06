import type Device from '../../device/Device';
import { deviceType } from '../constants';

export type tArrayDevices = Set<Device>;

export type tDeviceType = $Values<typeof deviceType>;
