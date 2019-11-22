import React, { useState } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { get, cloneDeep } from 'lodash';
import Paper from '@material-ui/core/Paper';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import { I18n } from 'react-i18next';

import SaveIcon from '@material-ui/icons/Save';

import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import Button from '../../components/CustomButtons/Button';

import saveConfigs from '../../modules/setting/action';
import { ioDirection } from '../../modules/device/io/constants';
import { ioOutputs, ioInputs } from '../../modules/io/constants';
// import { testIO } from '../../modules/external/device/io/saga';

import styles from './styles';
import withKeyboard from '../../components/Keyboard';
import { makeStyles } from '@material-ui/styles';

const testIO = () => {
};

const mapStateToProps = (state, ownProps) => ({
  storedConfigs: state.setting.page.modbus,
  ioEnabled: state.setting.systemSettings.modbusEnable,
  ...ownProps,
  ioModule: state.io.ioModule
});

function mapDispatchToProps(dispatch) {
  return {
    saveConfigs: (...args) => dispatch(saveConfigs(...args)),
    resetIO: (...args) => dispatch(resetIO(...args))
  };
}

/* eslint-disable react/prefer-stateless-function */
function ConnectedIo(props) {
  const { ioEnabled, ioModule } = props;
  const classes = makeStyles(styles.content)();


  const [data, setData] = useState(props.storedConfigs);
  const [isDataValid, setIsDataValid] = useState(true);
  const [section, setSection] = useState('modbus');
  const [btnGroupStatus, setBtnGroupStatus] = useState({});
  const [ioTestRsp, setIoTestRsps] = useState(Array(props?.storedConfigs?.in?.length).fill(0));

  if (!ioModule) {
    return null;
  }
  // 获取 btns 的状态集
  function getBtnStatus(data) {
    const checkEveryBtns = items =>
      items.reduce(
        (pre, item) => ({
          ...pre,
          [item.bit]: !(
            get(item, 'label', '').length === 0 ||
            get(item, 'function', '').length === 0
          )
        }),
        {}
      );

    const status = Object.keys(data).reduce(
      (pre, io) => ({
        ...pre,
        [io]: checkEveryBtns(data[io])
      }),
      {}
    );
    return status;
  }


  function handlePortChange(e, idx, io, key) {
    const tempData = cloneDeep(data[io]);
    tempData[idx][key] = get(e, 'target.value', '').trim();
    setData({
      ...data,
      [io]: tempData
    });
    setIsDataValid(validateData({
      ...data,
      [io]: tempData
    }));
  }

  function handleSubmit() {
    const { saveConfigs, resetIO } = props;
    resetIO(data);
    saveConfigs(section, data);
  }

  function handleTest(obj) {
    switch (obj.io) {
      case 'in': {
        testIO(obj.io, obj.bit)
          .then(resp => {
            const d = resp.response.body.valuesAsArray[0];
            const [...retIO] = ioTestRsp;
            retIO[obj.bit] = d;
            this.setState({
              ioTestRsp: retIO
            });
            return true;
          })
          .catch(() => {
            const [...retIO] = ioTestRsp;
            retIO[obj.bit] = 'fail';
            this.setState({
              ioTestRsp: retIO
            });
            return true;
          });
        break;
      }
      case 'out': {
        return testIO(obj.io, obj.bit);
      }
      default:
        break;
    }
  }

  function validateData(d = data) {
    // return Object.keys(data).every(io => data[io].every(item => item.label && item.function));
    return true;
  }

  function renderTestResult(status) {
    const cls = {
      0: {
        color: classes.info,
        text: 'OFF',
        textColor: classes.infoText
      },
      1: {
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
    const { color = classes.info, text = '', textColor = classes.infoText } = cls[status];
    return <div className={classes.statusWrap}>
      <span className={`${classes.statusCircle} ${color}`}/>
      <span className={textColor}>{text}</span>
    </div>;
  }

  function generatorItems(ports = [], direction, t) {
    if (!ports) {
      return null;
    }
    const data = ports.filter(p => p.direction === direction);

    return data.map((item, idx) => {
      let options = [];
      if (direction === ioDirection.output) {
        options = Object.keys(ioOutputs);
      } else if (direction === ioDirection.input) {
        options = Object.keys(ioInputs);
      }
      const renderTest = () => {
        if (direction === ioDirection.output) return null;
        return renderTestResult(ioTestRsp[item.bit]);
      };
      return (
        <React.Fragment key={`${item.io}_${item.bit}`}>
          <ListItem className={classes.inputItem}>
            <InputLabel className={classes.ioInputLabel} htmlFor="name-simple">
              {item.idx}
            </InputLabel>
            <Input
              id="name-simple"
              placeholder={t('Common.isRequired')}
              className={classes.ioInput}
              value={item.label}
              // onClick={() => {
              //   this.props.keyboardInput({
              //     onSubmit: text => {
              //       const tempData = cloneDeep(this.state.data[item.io]);
              //       tempData[idx]['label'] = text;
              //       this.setState({
              //         ...this.state,
              //         data: {
              //           ...this.state.data,
              //           [item.io]: tempData
              //         },
              //         isDataValid: this.validateData({
              //           ...this.state.data,
              //           [item.io]: tempData
              //         })
              //       });
              //     },
              //     text: item.label,
              //     title: item.bit,
              //     label: item.bit
              //   });
              // }}
              // onChange={e => this.handleChange(e, idx, item.io, 'label')}
            />
            <Select
              value={item.function}
              onChange={e => handlePortChange(e, idx, item, 'function')}
              displayEmpty
              name="function"
              className={classes.ioFunctionSelect}
            >
              {['', ...options].map(v => (
                <MenuItem key={v} value={v}>{v}</MenuItem>
              ))}
            </Select>
            <Button
              color="warning"
              size="lg"
              onClick={() => handleTest(item)}
              className={classes.testButton}
              // disabled={(!get(btnGroupStatus, `${item.io}.${item.bit}`, [])) || !ioEnabled}
            >
              {t('Common.Test')}
            </Button>
            {renderTest()}
          </ListItem>
          {idx !== data.length - 1 ?
            <Divider/>
            : null}
        </React.Fragment>
      );
    });
  }

  const { ports } = ioModule;
  const inItems = t => generatorItems(ports, ioDirection.input, t);
  const outItems = t => generatorItems(ports, ioDirection.output, t);

  return (
    <I18n ns="translations">
      {t => (
        <section className={classes.section}>
          <h3 className={classes.sectionTitle}>
            {t('Configuration.IO.INname')}
          </h3>
          <Paper className={classes.paperWrap} elevation={1}>
            <List>{inItems(t)}</List>
          </Paper>
          <h3
            className={`${classes.sectionTitle} ${classes.sectionTitleInner}`}
          >
            {t('Configuration.IO.OUTname')}
          </h3>
          <Paper className={classes.paperWrap} elevation={1}>
            <List>{outItems(t)}</List>
          </Paper>
          <Button
            variant="contained"
            disabled={!isDataValid || !ioEnabled}
            color="info"
            onClick={handleSubmit}
            className={classes.button}
          >
            <SaveIcon className={classes.leftIcon}/>
            {t('Common.Submit')}
          </Button>
        </section>
      )}
    </I18n>
  );
}

ConnectedIo.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  storedConfigs: PropTypes.shape({}).isRequired,
  saveConfigs: PropTypes.func.isRequired,
  resetIO: PropTypes.func.isRequired
};

const Io = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedIo);

export default withKeyboard(Io);
