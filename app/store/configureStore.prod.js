// @flow
import { createStore, applyMiddleware } from 'redux';
import createSagaMiddleware from 'redux-saga';
import { createHashHistory } from 'history';
import { routerMiddleware } from 'connected-react-router';
import createRootReducer from '../modules/indexReducer';
import type { tAction, Dispatch, StateType } from '../modules/typeDef';

const history = createHashHistory();
const rootReducer = createRootReducer(history);
const router = routerMiddleware(history);

function configureStore(initialState?: StateType) {
  const sagaMiddleware = createSagaMiddleware();
  const enhancer = applyMiddleware(sagaMiddleware, router);

  return {
    // eslint-disable-next-line flowtype/no-weak-types
    ...createStore<StateType, tAction<string, any>, Dispatch>(rootReducer, initialState, enhancer),
    runSaga: sagaMiddleware.run
  };
}

export default { configureStore, history };
