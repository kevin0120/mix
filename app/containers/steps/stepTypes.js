// @flow
import InputStep from './InputStep';
import ScannerStep from './scannerStep';
import InstructionStep from './instructionStep';
import TextStep from './TextStep';
import ScrewStep from './screwStep';
import MaterialStep from './materialStep';
import PassFailStep from './PassFailStep';
import MeasureStep from './MeasureStep';
import VideoStep from './VideoStep';
import { stepTypeKeys as s } from '../../modules/step/constants';

export default {
  [s.input]: {
    component: InputStep
  },
  [s.scanner]: {
    component: ScannerStep
  },
  [s.instruction]: {
    component: InstructionStep
  },
  [s.text]: {
    component: TextStep
  },
  [s.promiscuousTightening]: {
    component: ScrewStep
  },
  [s.screw]: {
    component: ScrewStep
  },
  [s.material]: {
    component: MaterialStep
  },
  [s.measure]: {
    component: MeasureStep
  },
  [s.passFail]: {
    component: PassFailStep
  },
  [s.video]: {
    component: VideoStep
  }
};

