odoo.define('approval_enhance.approval_form_view_hile_edit', function (require) {
'use strict';
    var session = require('web.session');
    var FormView = require('web.FormView');
    var core = require('web.core');
    
    $.fn.openerpBounce = function() {
        return this.each(function() {
            if (!$(this).hasClass('o_form_button_edit')){
                $(this).css('box-sizing', 'content-box').effect('bounce', {distance: 18, times: 5}, 250);
            }
        });
    };
    
    FormView.include({  
        
        
        load_record: function() {
            var self = this;
            return this._super.apply(this, arguments).then(function() {
                if (session.uid != 1)
                    self.show_hide_edit_button()
            });
        },
        
        show_hide_edit_button: function() {
            if (this.$buttons) {
                var button = this.$buttons.find('.o_form_button_edit');
                if(button) {
                    if (this.fields['approval_state'])
                    {
                        this.dataset.call('compute_show_edit',[[this.datarecord.id]]).then(function(result){
                            button.toggle(result);    
                        });
                    }
                }
            }
        }
             
    });
});
