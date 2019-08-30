import React  from 'react';
import i18n from 'i18next';
import Backend from 'i18next-xhr-backend';
import LngDetector from 'i18next-browser-languagedetector';
import { I18n, reactI18nextModule } from 'react-i18next';

export const lng = {
  'zh_CN': 'zh_CN',
  'en': 'en'
};

i18n
  .use(Backend)
  .use(LngDetector)
  .use(reactI18nextModule)
  .init({
    lng: 'zh_CN',
    fallbackLng: 'zh_CN',
    initImmediate: false,
    load: 'all',
    preload: ['en', 'zh_CN'],
    // have a common namespace used around the full app
    ns: ['translations'],
    defaultNS: 'translations',
    backend: {
      loadPath: './locales/{{lng}}/{{ns}}.json'
    },

    debug: false,

    interpolation: {
      escapeValue: false // not needed for react!!
    },

    react: {
      wait: true
    }
  });

export default i18n;

export function makeLocalBundle(lng, ns, resources) {
  return { lng, ns, resources };
}

export function tNS(str, ns) {
  return i18n.t(str, { ns: ns || 'translations' });
}

export function withI18n(tComp,ns){
  return <I18n ns={ns} >
    {tComp}
  </I18n>
}
