import { combineModels } from 'saga-modux';
import io from './io';

const models = [
  io

];

export const { reducer: rootReducer, saga: rootSaga } = combineModels(
  models,
  s => s.modules
);