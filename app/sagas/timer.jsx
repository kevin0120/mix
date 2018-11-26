import { eventChannel, END } from 'redux-saga';
import {
  take,
  call,
  race,
  put,
  cancelled,
  cancel,
  fork
} from 'redux-saga/effects';

function intervalChannel(ms) {
  return eventChannel(emit => {
    const iv = setInterval(() => {
      emit('tick');
    }, ms);

    return () => {
      clearInterval(iv);
    };
  });
}

export function* Interval(ms, { onStart, onStop, onCancel, onTick }) {
  const chan = yield call(intervalChannel, ms);
  if (onStart) {
    yield call(...onStart);
  }
  let onTickTasks = [];
  try {
    while (true) {
      yield take(chan);
      // console.log('Tick');
      if (onTick) {
        onTickTasks.push(yield fork(...onTick));
        for (let task in onTickTasks) {
          if (!onTickTasks[task].isRunning()) {
            onTickTasks.splice(task, 1);
          }
        }
      }
    }
  } finally {
    if (onTickTasks && onTickTasks.length > 0) {
      yield cancel(...onTickTasks);
      onTickTasks = null;
    }
    if (!(yield cancelled())) {
      console.log('Interval Stopped');
      if (onStop) {
        yield call(...onStop);
      }
    } else if (onCancel) {
      console.log('Interval Cancelled');
      yield call(...onCancel);
    }
    chan.close();
  }
}

// export function* Interval(ms, tasks) {
//   const { onStart, onStop, onCancel, onTick } = tasks;
//   let interval=null;
//   try {
//     while (true) {
//       interval=yield fork(runInterval, ms, { onStart, onStop, onCancel, onTick });
//     }
//   } finally {
//     if(interval){
//       yield cancel(interval);
//     }
//     console.log('Interval terminated');
//   }
// }
