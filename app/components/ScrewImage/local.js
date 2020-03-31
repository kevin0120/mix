import { makeLocalBundle, lng } from '../../i18n';


export const translation = {
  tool:'tool'
};


const trans = [
  makeLocalBundle(lng.zh_CN, 'point', {
    [translation.tool]: '拧紧工具',
  }),
  makeLocalBundle(lng.en, 'point', {
    [translation.tool]: 'tool',
  })
];

export default trans;
