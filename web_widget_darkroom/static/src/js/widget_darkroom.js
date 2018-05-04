/**
*    Copyright 2013 Matthieu Moquet
*    Copyright 2016-2017 LasLabs Inc.
*    License MIT (https://opensource.org/licenses/MIT)
**/

odoo.define('web_widget_darkroom.darkroom_widget', function(require) {
    'use strict';

    var core = require('web.core');
    var common = require('web.form_common');
    var session = require('web.session');
    var utils = require('web.utils');

    var QWeb = core.qweb;

    var FieldDarkroomImage = common.AbstractField.extend(common.ReinitializeFieldMixin, {
        className: 'darkroom-widget',
        template: 'FieldDarkroomImage',
        placeholder: "/web/static/src/img/placeholder.png",
        darkroom: null,
        no_rerender: false,

        defaults: {
            // Canvas initialization size
            cssMaxWidth: 700,
            cssMaxHeight: 500,
            selectionStyle: {
                cornerSize: 20,
                rotatingPointOffset: 70
            }
        },

        init: function(field_manager, node) {
            this._super(field_manager, node);
            this.options = _.defaults(this.options, this.defaults);
        },

        _init_darkroom: function() {
            if (!this.darkroom) {
                this._init_darkroom_icons();
                this.bind_events();
            }
        },

        _init_darkroom_icons: function() {
            var element = document.createElement('div');
            element.id = 'darkroom-icons';
            // element.style.height = 0;
            // element.style.width = 0;
            // element.style.position = 'absolute';
            // element.style.visibility = 'hidden';
            element.innerHTML = '<!-- inject:svg --><!-- endinject -->';
            this.el.appendChild(element);
        },

        bind_events: function () {
            var self = this;

            // event: click on '(Un)Follow' button, that toggles the follow for uid
            this.$el.on('click', '.o_darkroom_button_new', function () {
                if(self.darkroom){
                    console.log('new');
                    self.darkroom.addShape('circle', {
                        fill: 'red',
                        stroke: 'blue',
                        strokeWidth: 3,
                        rx: 10,
                        ry: 100,
                        isRegular: false
                    }).then(function (objectProps) {
                        console.log(objectProps.id);
                    });
                }
            });
            this.$el.on('click', '.o_darkroom_button_save', function () {
                console.log('save')
            });
            this.$el.on('click', '.o_darkroom_button_edit', function () {
                console.log('edit')
            });
        },

        destroy_content: function() {
            if (this.darkroom && this.darkroom.containerElement) {
                this.darkroom.containerElement.remove();
                this.darkroom = null;
            }
        },

        set_value: function(value) {
            return this._super(value);
        },

        render_value: function() {
            var self = this;
            this.destroy_content();
            this._init_darkroom();

            var url = null;
            if (this.get('value') && !utils.is_bin_size(this.get('value'))) {
                url = 'data:image/png;base64,' + this.get('value');
            } else if (this.get('value')) {
                var id = JSON.stringify(this.view.datarecord.id || null);
                var field = this.name;
                if (this.options.preview_image) {
                    field = this.options.preview_image;
                }
                url = session.url('/web/image', {
                    model: this.view.dataset.model,
                    id: id,
                    field: field,
                    unique: (this.view.datarecord.__last_update || '').replace(/[^0-9]/g, ''),
                });
            } else {
                url = this.placeholder;
            }

            var $img = $(QWeb.render("FieldBinaryImage-img", {widget: this, url: url}));
            this.$el.find('> img').remove();
            var _t = this.$('.tui-image-editor');
            // this.$el.append($img);
            this.darkroom = new tui.ImageEditor(_t, this.options);
            this.darkroom.loadImageFromURL('https://raw.githubusercontent.com/nhnent/tui.image-editor/master/examples/img/bg.jpg', 'tttt').then(
                function (x) {
                    console.log(x);
                    self.darkroom.clearUndoStack();
                }
            );
        },

        commit_value: function() {
            if (this.darkroom.sourceImage) {
                this.set_value(this.darkroom.sourceImage.toDataURL().split(',')[1]);
            }
        },
    });

    core.form_widget_registry.add("darkroom", FieldDarkroomImage);

    return {FieldDarkroomImage: FieldDarkroomImage};
});
