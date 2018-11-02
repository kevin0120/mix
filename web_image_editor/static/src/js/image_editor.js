odoo.define('web_widget_darkroom.image_editor', function(require) {
    'use strict';

    var core = require('web.core');
    var Model = require('web.DataModel');
    var common = require('web.form_common');
    var session = require('web.session');
    var utils = require('web.utils');

    var QWeb = core.qweb;

    var ImageEditor =  common.AbstractField.extend(common.ReinitializeFieldMixin, {
        template: 'ImageEditor',
        placeholder: "/web/static/src/img/placeholder.png",
        events:{
            'click .js_add_mask': 'add_new_mask',
            'click .js_remove_mask': 'remove_last_mask',
            'click .js_remove_all_mask': 'remove_all_mask',
            'click .js_save': 'save_all_mask',
        },
        markPoints: [],

        init: function(field_manager, node) {
            this._super(field_manager, node);
            this.markPoints.splice(0, this.markPoints.length); //清空mark点位
        },

        save_all_mask: function(event) {
            var self = this;
            var active_id = this.view.dataset.context.active_record_id;
            var url = '/api/v1/mrp.routing.workcenter/';
            var markPoints =  JSON.stringify(self.markPoints);
            $.ajax({
                type: "PUT",
                url : url.concat(active_id, '/points_edit'),
                timeout : 2000, //超时时间设置，单位毫秒
                dataType: 'json',
                data: markPoints,
                beforeSend: function(xhr) {xhr.setRequestHeader('Content-Type', 'application/json');},
            }).then(function(response){
                
            })

        },

        inline_add_new_mask: function(top, left) {
            var self = this;
            var leftOffset = left || 0;
            var topOffset = top || 0;
            var t = _.str.sprintf('<div class="oe_mark_circle">%s</div>', _.str.escapeHTML(this.markPoints.length + 1));
            $(t).css({'left': "" + leftOffset + "px",'top': "" + topOffset + "px"}).draggable({addClasses: false,
                stop: this.MarkerDragstop.bind(this),
                containment: "parent",
            }).prependTo(self.$el.find('#img_block'));
            this.markPoints.push({sequence: this.markPoints.length + 1 ,x_offset: leftOffset, y_offset: topOffset});
        },

        add_new_mask: function(event) {
            var self = this;
            var t = _.str.sprintf('<div class="oe_mark_circle">%s</div>', _.str.escapeHTML(this.markPoints.length + 1));
            $(t).draggable({addClasses: false,
                stop: this.MarkerDragstop.bind(this),
                containment: "parent",
            }).prependTo(self.$el.find('#img_block'));
            this.markPoints.push({sequence: this.markPoints.length + 1 ,x_offset: 0, y_offset: 0});
        },

        MarkerDragstop: function(event, ui) {
            var self = this;
            console.log(this, event, ui);
            var idx = parseInt(event.target.textContent) -1;
            self.markPoints[idx].x_offset = ui.position.left;
            self.markPoints[idx].y_offset = ui.position.top;
        },

        remove_last_mask: function() {
            var self = this;
            self.$('.oe_mark_circle:first').remove(); //添加是prepend
            this.markPoints.pop();
        },

        remove_all_mask: function() {
            var self = this;
            self.$('.oe_mark_circle').remove();
            this.markPoints.splice(0, this.markPoints.length);
        },

        set_dimensions: function(height, width) {
            console.log('set_dimensions',height, width);
            this.$el.css({
                width: width,
                height: height,
            });
        },

        destroy: function() {
            this._super.apply(this, arguments);
            this.markPoints.splice(0, this.markPoints.length);
        },

        render_value: function() {
            var url = this.placeholder;
            if (this.get('value')) {
                if (!utils.is_bin_size(this.get('value'))) {
                    url = 'data:image/png;base64,' + this.get('value');
                } else {
                    url = session.url('/web/image', {
                        model: this.view.dataset.model,
                        id: JSON.stringify(this.view.datarecord.id || null),
                        field: (this.options.preview_image) ? this.options.preview_image : this.name,
                        unique: (this.view.datarecord.__last_update || '').replace(/[^0-9]/g, ''),
                    });
                }
            }

            var $img = $(QWeb.render("FieldBinaryImage-img", {widget: this, url: url}));

            var self = this;

            console.log(self);

            this.$('> img').remove();
            if (self.options.size) {
                $img.css("width", "" + self.options.size[0] + "px");
                $img.css("height", "" + self.options.size[1] + "px");
            }
            var t = _.str.sprintf('url(%s)',url);
            this.$el.find('#img_block').css('background-image',t);

            var active_id = this.view.dataset.context.active_record_id;
            var model = new Model(this.view.dataset.context.active_model).call("get_operation_points", [active_id]).then(function (result) {
                    _.each(result, function (ele) {
                        self.inline_add_new_mask(ele.y_offset, ele.x_offset);
                    })
            });

        },
    });

    core.form_widget_registry.add('image_editor',ImageEditor);

});