export const WORKABLE={
  DID_FINISH: 'WORKABLE_DID_FINISH'
};

export function workableDidFinish(workable){
  return{
    type:WORKABLE.DID_FINISH,
    workable
  }
}