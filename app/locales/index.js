import i18n from '../i18n';
import stepWorking from '../containers/stepWorking/local';
import navBar from '../components/NavBar/local';
import rework from '../modules/reworkPattern/local';
import points from '../components/ScrewImage/local';



const trans = [
  ...rework,
  ...stepWorking,
  ...navBar,
  ...points
];

export function loadLocales() {
  trans.forEach((t) => {
    i18n.addResourceBundle(t.lng, t.ns, t.resources);
  });
}

