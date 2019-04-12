odoo.define('web_widget_echart', function (require) {
    "use strict";

    var core = require('web.core');
    var form_common = require('web.form_common');
    var formats = require('web.formats');
    var Model = require('web.Model');

    var QWeb = core.qweb;


    var EChartWidget = form_common.AbstractField.extend({
        render_value: function () {
            var val = JSON.parse(this.get('value'));
            console.log(val);
            var chart = echarts.init(this.el, null, {width: 900, height: 500});
            var xLabel = 'cur_w';
            var yLabel = 'cur_m';

            let series = [];
            for (var v in val) {
                if (val[v]) {
                    series.push({
                        data: this.genData(val[v][xLabel], val[v][yLabel]),
                        type: 'line',
                        name: val[v]['name'],
                        markLine: {
                            data: [
                                {type: 'min', name: '最小值'},
                                {type: 'max', name: '最大值'}
                            ]
                        },
                        smooth: true
                    })
                }
            }
            console.log(series);
            var option = {
                title: {
                    text: ''
                },
                tooltip: {
                    trigger: 'axis'
                },
                xAxis: {},
                yAxis: {},
                dataZoom: [
                    {
                        type: 'slider',
                        xAxisIndex: 0,
                        filterMode: 'empty'
                    },
                    {
                        type: 'slider',
                        yAxisIndex: 0,
                        filterMode: 'empty'
                    },
                    {
                        type: 'inside',
                        xAxisIndex: 0,
                        filterMode: 'empty'
                    },
                    {
                        type: 'inside',
                        yAxisIndex: 0,
                        filterMode: 'empty'
                    }
                ],
                series
            };
            chart.setOption(option);
            chart.resize();
        },

        genData: function (xData, yData) {
            var data = [];
            for (var x in xData) {
                data.push([xData[x], yData[x]])
            }
            return data;
        }
    });
    core.form_widget_registry.add('echart', EChartWidget);
    return {
        EChartWidget: EChartWidget
    };
});
