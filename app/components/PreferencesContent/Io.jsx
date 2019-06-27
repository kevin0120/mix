import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';
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
import { bindActionCreators } from 'redux';
import Button from '../CustomButtons/Button';

import saveConfigs from '../../modules/setting/userConfigs';
import { IO_FUNCTION } from '../../modules/io/model';

import { testIO } from '../../modules/io/saga';
import { resetIO } from '../../modules/io/action';

import styles from './styles';
import withKeyboard from '../Keyboard';

const mapStateToProps = (state, ownProps) => ({
  storedConfigs: state.setting.page.modbus,
  ioEnabled: state.setting.systemSettings.modbusEnable,
  ...ownProps
});

function mapDispatchToProps(dispatch) {
  return {
    saveConfigs: (...args) => dispatch(saveConfigs(...args)),
    resetIO: (...args) => dispatch(resetIO(...args))
  };
}

/* eslint-disable react/prefer-stateless-function */
class ConnectedIo extends React.PureComponent {
  // 获取 btns 的状态集
  static getBtnStatus(data) {
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

  constructor(props) {
    super(props);
    this.state = {
      isDataValid: true,
      data: props.storedConfigs,
      btnGroupStatus: {},
      section: 'modbus',
      ioTestRsp: Array(props.storedConfigs.in.length).fill(0)
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleTest = this.handleTest.bind(this);

    this.IO_FUNCTION = IO_FUNCTION;
  }

  handleChange(e, idx, io, key) {
    const { data } = this.state;
    const tempData = cloneDeep(data[io]);
    tempData[idx][key] = get(e, 'target.value', '').trim();
    this.setState({
      data: {
        ...data,
        [io]: tempData
      },
      isDataValid: this.validateData({
        ...data,
        [io]: tempData
      })
    });
  }

  handleSubmit() {
    const { saveConfigs, resetIO } = this.props;
    const { data, section } = this.state;
    resetIO(data);
    saveConfigs(section, data);
  }

  handleTest(obj) {
    const { ioTestRsp } = this.state;
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

  validateData(data = this.state.data) {
    // return Object.keys(data).every(io => data[io].every(item => item.label && item.function));
    return true;
  }

  generatorItems(data, t) {
    const { classes, ioEnabled } = this.props;
    const { btnGroupStatus, ioTestRsp } = this.state;
    return data.map((item, idx) => {
      const options = get(
        this.IO_FUNCTION,
        String.prototype.toUpperCase.call(item.io),
        this.IO_FUNCTION.IN
      );
      const selectItems = Object.keys(options).map(v => (
        <MenuItem key={v} value={v}>
          {v}
        </MenuItem>
      ));

      const generatorStatus = () => {
        if (item.io === 'out') return null;

        if (ioTestRsp[item.bit] === 0) {
          return (
            <div className={classes.statusWrap}>
              <span className={`${classes.statusCircle} ${classes.info}`}/>
              <span className={classes.infoText}>OFF</span>
            </div>
          );
        }
        if (ioTestRsp[item.bit] === 1) {
          return (
            <div className={classes.statusWrap}>
              <span className={`${classes.statusCircle} ${classes.success}`}/>
              <span className={classes.successText}>ON</span>
            </div>
          );
        }
        if (ioTestRsp[item.bit] === 'fail') {
          return (
            <div className={classes.statusWrap}>
              <span className={`${classes.statusCircle} ${classes.fail}`}/>
              <span className={classes.failText}>Fail</span>
            </div>
          );
        }
        return null;
      };
      return (
        <div key={`${item.io}_${item.bit}`}>
          <ListItem className={classes.inputItem}>
            <InputLabel className={classes.ioInputLabel} htmlFor="name-simple">
              {item.bit}
            </InputLabel>
            <Input
              id="name-simple"
              placeholder={t('Common.isRequired')}
              className={classes.ioInput}
              value={item.label}
              onClick={() => {
                this.props.keyboardInput({
                  onSubmit: text => {
                    const tempData = cloneDeep(this.state.data[item.io]);
                    tempData[idx]['label'] = text;
                    this.setState({
                      ...this.state,
                      data: {
                        ...this.state.data,
                        [item.io]: tempData
                      },
                      isDataValid: this.validateData({
                        ...this.state.data,
                        [item.io]: tempData
                      })
                    });
                  },
                  text: item.label,
                  title: item.bit,
                  label: item.bit
                });
              }}
              // onChange={e => this.handleChange(e, idx, item.io, 'label')}
            />
            <Select
              value={item.function}
              onChange={e => this.handleChange(e, idx, item.io, 'function')}
              displayEmpty
              name="function"
              className={classes.ioFunctionSelect}
            >
              <MenuItem key="dummy" value=""/>
              {selectItems}
            </Select>
            <span>{String.prototype.toUpperCase.call(item.io)}</span>
            <Button
              color="warning"
              size="lg"
              onClick={() => this.handleTest(item)}
              className={classes.testButton}
              disabled={(!get(btnGroupStatus, `${item.io}.${item.bit}`, [])) || !ioEnabled}
            >
              {t('Common.Test')}
            </Button>
            {generatorStatus()}
          </ListItem>
          <li>
            <Divider/>
          </li>
        </div>
      );
    });
  }

  render() {
    const { classes,ioEnabled } = this.props;
    const { data, isDataValid } = this.state;
    const inItems = t => this.generatorItems(data.in, t);
    const outItems = t => this.generatorItems(data.out, t);

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
              <Button
                variant="contained"
                disabled={!isDataValid|| !ioEnabled}
                color="info"
                onClick={this.handleSubmit}
                className={classes.button}
              >
                <SaveIcon className={classes.leftIcon}/>
                {t('Common.Submit')}
              </Button>
            </Paper>
          </section>
        )}
      </I18n>
    );
  }
}

ConnectedIo.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  storedConfigs: PropTypes.shape({}).isRequired,
  // in: PropTypes.array,
  // out: PropTypes.array,
  saveConfigs: PropTypes.func.isRequired,
  resetIO: PropTypes.func.isRequired
};

const Io = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedIo);

export default withKeyboard(withStyles(styles)(Io));
