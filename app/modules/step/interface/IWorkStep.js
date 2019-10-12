export interface IWorkStep {
  _onLeave: (any)=>void;
  run: ()=>any;

  timerStart: ()=>any;

  timerStop: ()=>any;

  updateData: ()=>any;
}
