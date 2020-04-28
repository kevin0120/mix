// @flow
import { makeStyles, Paper } from '@material-ui/core';
import { connect } from 'react-redux';
import React, { useEffect, useRef } from 'react';
import ReactDOM from 'react-dom';
import styles from './styles';
import Panel from './Panel';
import { withI18n } from '../../../i18n';
import * as oSel from '../../../modules/order/selector';

type DProps = {||};

type ownProps = {||};

type SProps = {|
  ...ownProps
|};

type Props = {|
  ...ownProps,
  ...SProps,
  ...DProps
|};

const mapState = (state, props: ownProps): SProps => ({
  ...props,
  payload: oSel.viewingOrder(state.order)?.payload || {}
});

const mapDispatch: DProps = {};

const data = {
  tools: [{
    name: '工具1aaaaaaaaaaaaaaaaaaa',
    model: 'model1',
    sn: '11',
    num: 1
  }]
};


// 订单人机料法环信息
function OrderInfoLeft({ payload }: Props) {
  const classes = makeStyles(styles.OrderInfoStyles)();

  const viewingRef = useRef(null);
  const viewingNode = viewingRef?.current;
  useEffect(() => {
    if (viewingNode) {
      // eslint-disable-next-line react/no-find-dom-node
      const node: null | Element | Text = ReactDOM.findDOMNode(viewingNode);
      if (node && node.scrollIntoView && typeof node.scrollIntoView === 'function') {
        ((node: any): Element).scrollIntoView({ block: 'center', behavior: 'smooth' });
      }
    }
  }, [viewingNode]);
  const {
    STARTEMPLOYEE: startEmployee, // 人员
    MOMCONFIG: {
      TOOLS: tools,
      DEVICES: devices,
      MOMWIPORDERCOMS: materials,
      SECURITY: security,
      RESOURCENO,
      PRODUCTNAME
    } = {},
    environments
  } = payload;
  console.log(payload);
  return withI18n(t => (
    <Paper square className={classes.root}>
      <Panel
        title="产品"
        cols={[{
          label: '产品名称',
          name: 'PRODUCTNAME'
        }, {
          label: '物料编码',
          name: 'RESOURCENO'
        }]}
        data={[{ RESOURCENO, PRODUCTNAME }]}
        classes={classes}
      />
      <Panel
        title="物料"
        cols={[{
          label: '物料编号',
          name: 'NAME'
        }, {
          label: '物料名称',
          name: 'NO'
        }, {
          label: '数量',
          name: 'QUANTITY'
        }]}
        data={materials}
        classes={classes}
      />
    </Paper>
  ));
}

function OrderInfoRight({ payload }: Props) {
  const classes = makeStyles(styles.OrderInfoStyles)();

  const viewingRef = useRef(null);
  const viewingNode = viewingRef?.current;
  useEffect(() => {
    if (viewingNode) {
      // eslint-disable-next-line react/no-find-dom-node
      const node: null | Element | Text = ReactDOM.findDOMNode(viewingNode);
      if (node && node.scrollIntoView && typeof node.scrollIntoView === 'function') {
        ((node: any): Element).scrollIntoView({ block: 'center', behavior: 'smooth' });
      }
    }
  }, [viewingNode]);
  const {
    STARTEMPLOYEE: startEmployee, // 人员
    MOMCONFIG: {
      TOOLS: tools,
      DEVICES: devices,
      MOMWIPORDERCOMS: materials,
      SECURITY: security,
      RESOURCENO,
      PRODUCTNAME
    } = {},
    environments
  } = payload;
  console.log(payload);
  return withI18n(t => (
    <Paper square className={classes.root}>
      <Panel
        title="人员"
        cols={[{
          label: '用户名',
          name: 'name'
        }]}
        data={[{ name: startEmployee }]}
        classes={classes}
      />
      <Panel
        title="工具"
        cols={[{
          label: '工具名称',
          name: 'NAME'
        }, {
          label: '规格型号',
          name: 'TYPE'
        }]}
        data={tools}
        classes={classes}
      />
      <Panel
        title="设备"
        cols={[{
          label: '设备名称',
          name: 'name'
        }]}
        data={(devices || []).map(d => ({ name: d }))}
        classes={classes}
      />
      <Panel
        title="环境要求"
        cols={[{
          label: '环境要求',
          name: 'name'
        }]}
        data={(environments || []).map(e => ({ name: e }))}
        classes={classes}
      />
      <Panel
        title="劳动防护"
        cols={[{
          label: '劳动防护',
          name: 'CPCLASS'
        }]}
        data={security}
        classes={classes}
      />
    </Paper>
  ));
}

export default {
  OrderInfoLeft: connect(mapState, mapDispatch)(OrderInfoLeft),
  OrderInfoRight: connect(mapState, mapDispatch)(OrderInfoRight)
}
