export interface IWorkStep {
  id: string,
  _onLeave: (any)=>void,
  run: ()=>any,
  timerStart: ()=>any,
  timerStop: ()=>any,
  updateData: ()=>any
}
