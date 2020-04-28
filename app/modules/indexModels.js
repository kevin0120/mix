import { combineModels } from 'saga-modux';
import io from './io';
import viewOperationInfo from './viewOperationInfo';

const models = [
  io,
  viewOperationInfo
];

export const { reducer: rootReducer, saga: rootSaga } = combineModels(
  models,
  s => s.modules
);
