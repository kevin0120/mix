import i18n from '../i18n';
import stepWorking from '../containers/stepWorking/local';

const trans = [
  ...stepWorking

];

export function loadLocales(){
  trans.forEach((t) => {
    i18n.addResourceBundle(t.lng, t.ns, t.resources);
  });
}

