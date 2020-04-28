// @flow
import type { Node } from 'react';
import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import { cloneDeep, get } from 'lodash';
import Paper from '@material-ui/core/Paper';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import SaveIcon from '@material-ui/icons/Save';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import { makeStyles } from '@material-ui/styles';
import Button from '../../components/CustomButtons/Button';
import { ioDirection } from '../../modules/device/io/constants';
import { ioInputs, ioOutputs } from '../../modules/io/constants';
import styles from './styles';
import withKeyboard from '../../components/Keyboard';
import { withI18n } from '../../i18n';
import io from '../../modules/io';

const mapStateToProps = (state, ownProps) => ({
  ...ownProps,
  ioEnabled: io.select(state).enabled,
  ioModule: io.select(state).ioModule,
  testStatus: io.select(state).testStatus,
  ioPorts: io.select(state).ioPorts
});

const mapDispatchToProps = {
  testIO: io.action.testPort,
  setPortsConfig: io.action.setPortsConfig
};

function ConnectedIo(props): ?Node {
  const {
    ioEnabled,
    ioModule,
    testIO,
    testStatus,
    ioPorts,
    setPortsConfig
  } = props;
  const classes = makeStyles(styles.content)();

  const [config, setConfig] = useState(ioPorts);
  const [isConfigValid, setIsConfigValid] = useState(true);

  useEffect(() => {
    // eslint-disable-next-line no-unused-vars
    function validateConfig(d = config) {
      // return Object.keys(d).every(io => d[io].every(item => item.label && item.function));
      return true;
    }

    setIsConfigValid(validateConfig(config));
  }, [config]);

  if (!ioModule) {
    return null;
  }

  const { ports } = ioModule;

  function handlePortLabelChange() {
  }

  function handlePortFunctionChange(e, port) {
    const { idx, direction } = port;
    const tempConfig = cloneDeep(config);
    const newFunction = get(e, 'target.value', '').trim();
    const prevFunction = Object.keys(config[direction]).find(
      k => config[direction][k] === idx
    );
    if (prevFunction) {
      delete tempConfig[direction][prevFunction];
    }
    if (newFunction) {
      tempConfig[direction][newFunction] = idx;
    }
    setConfig(tempConfig);
  }

  function handleSubmit() {
    setPortsConfig({ ioPorts: config });
    // setIO(config);
  }

  function renderTestResult(status) {
    const cls = {
      '0': {
        color: classes.info,
        text: 'OFF',
        textColor: classes.infoText
      },
      '1': {
        color: classes.success,
        text: 'ON',
        textColor: classes.successText
      },
      Fail: {
        color: classes.fail,
        text: 'Fail',
        textColor: classes.failText
      }
    };
    if (!(status in Object.keys(cls))) {
      return null;
    }
    const {
      color = classes.info,
      text = '',
      textColor = classes.infoText
    } = cls[status];
    return (
      <div className={classes.statusWrap}>
        <span className={`${classes.statusCircle} ${color}`}/>
        <span className={textColor}>{text}</span>
      </div>
    );
  }

  function generatorItems(direction) {
    if (!ports) {
      return null;
    }
    const filteredPorts = ports.filter(p => p.direction === direction);

    const portHasFunction = p =>
      Object.values(config[p.direction]).some(v => v === p.idx);
    let options = [];
    if (direction === ioDirection.output) {
      options = Object.values(ioOutputs);
    } else if (direction === ioDirection.input) {
      options = Object.values(ioInputs);
    }

    return filteredPorts.map((port, idx) => {
      const currentFunction =
        Object.keys(config[direction]).find(
          k => port.idx === config[direction][k]
        ) || '';

      const renderTest = () => {
        if (direction === ioDirection.input) {
          return null;
        }
        return withI18n(t => (
          <React.Fragment>
            {renderTestResult(testStatus[port.idx])}
            <Button
              color="warning"
              size="lg"
              onClick={() => testIO({ port })}
              className={classes.testButton}
              disabled={!portHasFunction(port) || !ioEnabled}
            >
              {t('Common.Test')}
            </Button>
          </React.Fragment>
        ));
      };
      return withI18n(t => (
        <List key={`${port.idx}_${port.direction}`}>
          <ListItem className={classes.inputItem}>
            <InputLabel className={classes.ioInputLabel} htmlFor="name-simple">
              {port.idx}
            </InputLabel>
            <Input
              id="name-simple"
              placeholder={t('Common.isRequired')}
              className={classes.ioInput}
              value={port.idx}
              onClick={() => {
                // props.keyboardInput({
                //   onSubmit: text => {
                //     const tempConfig = cloneDeep(config[port.idx]);
                //     tempConfig[port.idx].label = text;
                //     setConfig({
                //       ...config,
                //       [port.idx]: tempConfig
                //     });
                //   },
                //   text: port.idx,
                //   title: port.idx,
                //   label: port.idx
                // });
              }}
              onChange={e => handlePortLabelChange(e, port)}
            />
            <Select
              value={currentFunction}
              onChange={e => {
                handlePortFunctionChange(e, port);
              }}
              displayEmpty
              name="function"
              className={classes.ioFunctionSelect}
            >
              {['', ...options].map(v => (
                <MenuItem key={v} value={v}>
                  {v}
                </MenuItem>
              ))}
            </Select>

            {renderTest()}
          </ListItem>
          {idx !== config.length - 1 ? <Divider/> : null}
        </List>
      ));
    });
  }

  const inItems = generatorItems(ioDirection.input);
  const outItems = generatorItems(ioDirection.output);

  return withI18n(t => (
    <section className={classes.section}>
      {ioModule.maxInputs > 0 ?
        <React.Fragment>
          <h3 className={classes.sectionTitle}>{t('Configuration.IO.INname')}</h3>
          <Paper className={classes.paperWrap} elevation={1}>
            {inItems}
          </Paper>
        </React.Fragment> : null}
      {ioModule.maxOutputs > 0 ?
        <React.Fragment>
          <h3 className={`${classes.sectionTitle} ${classes.sectionTitleInner}`}>
            {t('Configuration.IO.OUTname')}
          </h3>
          <Paper className={classes.paperWrap} elevation={1}>
            {outItems}
          </Paper>
        </React.Fragment> : null}
      <Button
        variant="contained"
        disabled={!isConfigValid || !ioEnabled}
        color="info"
        className={classes.button}
        onClick={handleSubmit}
      >
        <SaveIcon className={classes.leftIcon}/>
        {t('Common.Submit')}
      </Button>
    </section>
  ));
}

export default withKeyboard(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )(ConnectedIo)
);
