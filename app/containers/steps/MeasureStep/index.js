// @flow
import React, { useEffect, useRef, useState } from 'react';
import { connect } from 'react-redux';
import { makeStyles, Typography } from '@material-ui/core';
import classNames from 'classnames';
import Assignment from '@material-ui/icons/Assignment';
import { stepData, stepPayload, viewingStep } from '../../../modules/order/selector';
import stepActions from '../../../modules/step/actions';
import type { tStepProps } from '../types';
import Button from '../../../components/CustomButtons/Button';
import type { Dispatch } from '../../../modules/typeDef';
import type { tStepPayload } from '../../../modules/step/interface/typeDef';
import withKeyboard from '../../../components/Keyboard';
import styles from './styles';
import Card from '../../../components/Card/Card';
import CardHeader from '../../../components/Card/CardHeader';
import CardIcon from '../../../components/Card/CardIcon';
import CardBody from '../../../components/Card/CardBody';

type tOP = {|
  ...tStepProps
|};

type tSP = {|
  ...tOP,
  type: string,
  payload: ?tStepPayload
|};

type tDP = {|
  submit: Dispatch
|};

const mapState = (state, props: tOP): tSP => ({
  ...props,
  type: stepPayload(viewingStep(state.order))?.type || '',
  payload: stepPayload(viewingStep(state.order)),
  value: stepData(viewingStep(state.order)).result,
  target: viewingStep(state.order).target,
  max: viewingStep(state.order).toleranceMax,
  min: viewingStep(state.order).toleranceMin
});

const mapDispatch: tDP = {
  submit: stepActions.submit,
  measureInput: stepActions.input
};

type Props = {|
  ...tSP,
  ...tDP
|};

const isValidMeasureValue = (v) => typeof v === 'string' || typeof v === 'number';

function MeasureStep({
  step,
  isCurrent,
  submit,
  bindAction,
  keyboardInput,
  target,
  max,
  min,
  measureInput,
  value
}: Props) {
  const classes = makeStyles(styles)();
  useEffect(
    () => {
      bindAction(
        <Button
          type="button"
          color="primary"
          onClick={() => {
            submit(value);
          }}
          disabled={!isCurrent}
        >
          完成
        </Button>
      );
      return () => bindAction(null);
    },
    [step, bindAction, isCurrent, value, submit]
  );
  const { desc, text } = step;

  return <div className={classes.root}>
    {/*<Card className={classes.card} login>*/}
    {/*  <CardHeader color="info" icon>*/}
    {/*    <CardIcon color="info">*/}
    {/*      <Assignment/>*/}
    {/*    </CardIcon>*/}
    {/*    <Typography variant="h6" className={classes.cardIconTitle}>测量</Typography>*/}
    {/*  </CardHeader>*/}
    {/*  <CardBody className={classes.cardContent}>*/}
    {/*    <table>*/}
    {/*      <tbody>*/}
    {/*      {[*/}
    {/*        {*/}
    {/*          label: '目标值：',*/}
    {/*          content: isValidMeasureValue(target) ? target : '未指定'*/}
    {/*        }, {*/}
    {/*          label: '最小值：',*/}
    {/*          content: isValidMeasureValue(min) ? min : '未指定'*/}
    {/*        }, {*/}
    {/*          label: '最大值：',*/}
    {/*          content: isValidMeasureValue(max) ? max : '未指定'*/}
    {/*        }, {*/}
    {/*          label: '测量值：',*/}
    {/*          content: <span*/}
    {/*            className={classNames({*/}
    {/*              [classes.inputContainer]: true,*/}
    {/*              [classes.inputContainerDisabled]: !isCurrent*/}
    {/*            })}*/}
    {/*            onClick={() => {*/}
    {/*              if (isCurrent) {*/}
    {/*                keyboardInput({*/}
    {/*                  onSubmit: text => {*/}
    {/*                    measureInput({*/}
    {/*                      data: text,*/}
    {/*                      time: new Date(),*/}
    {/*                      source: 'input'*/}
    {/*                    });*/}
    {/*                  },*/}
    {/*                  text: value || '',*/}
    {/*                  title: '请输入测量值',*/}
    {/*                  label: '请输入测量值'*/}
    {/*                });*/}
    {/*              }*/}
    {/*            }}*/}
    {/*          ><Typography variant="h4" className={classNames({*/}
    {/*            [classes.inputText]: !!value && isCurrent,*/}
    {/*            [classes.inputTextDisabled]: !value || !isCurrent*/}
    {/*          })}>{value || '请输入'}</Typography></span>*/}
    {/*        }*/}
    {/*      ].map(row => (*/}
    {/*        <tr className={classes.row} key={row.label}>*/}
    {/*          <td><Typography variant="h4">{row.label}</Typography></td>*/}
    {/*          <td className={classes.rowContent}>{*/}
    {/*            typeof row.content === 'string' || typeof row.content === 'number' ?*/}
    {/*              <Typography variant="h4" className={classes.inputText}>{row.content}</Typography>*/}
    {/*              : row.content*/}
    {/*          }*/}
    {/*          </td>*/}
    {/*        </tr>*/}
    {/*      ))}*/}
    {/*      </tbody>*/}
    {/*    </table>*/}
    {/*  </CardBody>*/}
    {/*</Card>*/}
    <span className={classes.desc}>
    {text || desc ? (
      <Typography variant="h4">
        {text || desc}
      </Typography>
    ) : null}
    </span>
    <span
      className={classNames({
        [classes.inputContainer]: true,
        [classes.inputContainerDisabled]: !isCurrent
      })}
      onClick={() => {
        if (isCurrent) {
          keyboardInput({
            onSubmit: text => {
              measureInput({
                data: text,
                time: new Date(),
                source: 'input'
              });
            },
            text: value || '',
            title: '请输入测量值',
            label: '请输入测量值'
          });
        }
      }}
    ><Typography variant="h4" className={classNames({
      [classes.inputText]: !!value && isCurrent,
      [classes.inputTextDisabled]: !value || !isCurrent
    })}>{value || '请输入测量值'}</Typography></span>
  </div>;
}


export default connect<Props, tOP, tSP, tDP, _, _>(
  mapState,
  mapDispatch
)(withKeyboard(MeasureStep));
