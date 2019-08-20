export const MODEL_VIEWER={
  OPEN:'MODEL_VIEWER_OPEN',
};

export default {
  open:(url)=>({
    type:MODEL_VIEWER.OPEN,
    url
  })
}
