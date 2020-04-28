const { events } = require('../constants');

const odooTypes = {
  ODOO_STATUS: 'WS_ODOO_STATUS'
};

const odooHandlers = {
  [odooTypes.ODOO_STATUS]: (data, reply) => {
    reply(
      {
        sn: 0,
        data: {
          name: 'ODOO',
          status: 'online'
        }
      },
      events.odoo
    );
  }
};

module.exports = {
  odooHandlers,
  odooTypes
};
