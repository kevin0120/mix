import i18n from 'i18next';
import Backend from 'i18next-xhr-backend';
import LngDetector from 'i18next-browser-languagedetector';
import { reactI18nextModule } from 'react-i18next';

i18n
  .use(Backend)
  .use(LngDetector)
  .use(reactI18nextModule)
  .init({
    lng: 'en',
    fallbackLng: 'en',
    initImmediate: false,
    load: ['en', 'zh_CN'],
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
