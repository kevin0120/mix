odoo.define('web_widget_darkroom.image_editor', function(require) {
    'use strict';

    var core = require('web.core');
    var DataModel = require('web.DataModel');
    var common = require('web.form_common');
    var session = require('web.session');
    var utils = require('web.utils');

    var QWeb = core.qweb;

    var ImageEditor =  common.AbstractField.extend(common.ReinitializeFieldMixin, {
        template: 'ImageEditor',
        placeholder: "/web/static/src/img/placeholder.png",
        defaults: {
             includeUI: {
                 initMenu: 'filter',
                 menuBarPosition: 'bottom'
             },
             cssMaxWidth: 700,
             cssMaxHeight: 500
         },
        events:{
            'click .js_add_mask': 'add_new_mask',
        },

        init: function(field_manager, node) {
            this._super(field_manager, node);
            this.options = _.defaults(this.options);
        },

        add_new_mask: function() {
            console.log('add new mask')
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

            this.$('> img').remove();
            if (self.options.size) {
                $img.css("width", "" + self.options.size[0] + "px");
                $img.css("height", "" + self.options.size[1] + "px");
            }
            this.$el.prepend($img);

            // var el = self.$("#tui-image-editor-container");
            // console.log(el);
            // this.imageEditor = new tui.ImageEditor(this.$el,  {
            //       includeUI: {
            //         menu: ['shape', 'filter'],
            //         initMenu: 'filter',
            //         uiSize: {
            //             width: '800px',
            //             height: '900px'
            //         },
            //         menuBarPosition: 'bottom'
            //       },
            //       cssMaxWidth: 700,
            //       cssMaxHeight: 500,
            //       selectionStyle: {
            //         cornerSize: 20,
            //         rotatingPointOffset: 70
            //       }
            //     });
            // this.imageEditor.loadImageFromURL(url,'Image Editor').then(function (result) {
            //     console.log(result);
            // })
        },
    });

    core.form_widget_registry.add('image_editor',ImageEditor);

});