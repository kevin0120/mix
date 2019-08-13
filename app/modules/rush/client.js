

let ws=null;

export function getWSClient(){
  return ws;
}

export function setWSClient(newWS){
  ws=newWS;
}