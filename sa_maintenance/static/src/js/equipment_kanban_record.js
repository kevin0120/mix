odoo.define('equipment.kanban_record', function (require) {
"use strict";

var core = require('web.core');
var KanbanRecord = require('web_kanban.Record');
var Model = require('web.Model');
var _t = core._t;

KanbanRecord.include({
    init: function (parent, record, options) {
        this._super(parent, record, options);
        this.mytimer = null;
    },

    ping: function(){
      var self = this;
      var url = self.values.healthz_url.value;
      $.ajax({
          type: 'GET',
          url: url,
          timeout : 2000, //超时时间设置，单位毫秒
          dataType: 'json',
          beforeSend: function(xhr) {xhr.setRequestHeader('Content-Type', 'application/json');},
          success: function(){
              if(self.$(".oe_kanban_equpiment_status").hasClass('oe_kanban_equpiment_status_red')){
                  self.$(".oe_kanban_equpiment_status").removeClass('oe_kanban_equpiment_status_red').addClass('oe_kanban_equpiment_status_green');
              }
          },
          error: function() {
              if(self.$(".oe_kanban_equpiment_status").hasClass('oe_kanban_equpiment_status_green')){
                  self.$(".oe_kanban_equpiment_status").removeClass('oe_kanban_equpiment_status_green').addClass('oe_kanban_equpiment_status_red');
              }
          },
      })
    },

    destroy: function(){
        clearInterval(this.mytimer);
        this._super();
    },

    start: function() {
        var self = this;
        if (this.model === 'maintenance.equipment'){
            this._super.apply(this, arguments);
            self.$(".oe_kanban_equpiment_status").addClass('oe_kanban_equpiment_status_red');
            if(self.values.healthz_url.value) {
                self.mytimer = setInterval(self.ping.bind(self), 3000);
            }
        } else{
            return this._super.apply(this, arguments);
        }

    },
});

});
