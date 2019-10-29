import type Device from '../../external/device/Device';
import { deviceType } from '../constants';

export type tArrayDevices = Set<Device>;

export type tDeviceType = $Values<typeof deviceType>;
