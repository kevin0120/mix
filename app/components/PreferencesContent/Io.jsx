import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';
import { get, cloneDeep, isEqual } from 'lodash';

import Button from '../../components/CustomButtons/Button';
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

import saveConfigs from '../../actions/userConfigs';
import { IO_FUNCTION } from '../../reducers/io';


import {
  testModbus,
  resetIO,
} from '../../actions/ioModbus';

import styles from './styles';

const mapStateToProps = (state, ownProps) => ({
  storedConfigs: state.setting.page.modbus,
  section: 'modbus',
  ioTestRsp: state.ioTestRsp,
  ...ownProps,
});

const mapDispatchToProps = {
  saveConfigs,
  testModbus,
  resetIO,
};

/* eslint-disable react/prefer-stateless-function */
class ConnectedIo extends React.PureComponent {
  // 获取 btns 的状态集
  static getBtnStatus(data) {
    const checkEveryBtns = items => items.reduce((pre, item) => ({
      ...pre,
      [item.bit]: !(get(item, 'label', '').length === 0 || get(item, 'function', '').length === 0),
    }), {});

    const status = Object.keys(data).reduce((pre, io) => ({
      ...pre,
      [io]: checkEveryBtns(data[io]),
    }), {});
    return status;
  }

  constructor(props) {
    super(props);
    this.state = {
      isDataValid: true,
      data: props.storedConfigs,
      btnGroupStatus: {},
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleTest = this.handleTest.bind(this);

    this.IO_FUNCTION = IO_FUNCTION;
  }

  handleChange(e, idx, io, key) {
    const tempData = cloneDeep(this.state.data[io]);
    tempData[idx][key] = get(e, 'target.value', '').trim();
    this.setState({
      ...this.state,
      data: {
        ...this.state.data,
        [io]: tempData,
      },
      isDataValid: this.validateData({
        ...this.state.data,
        [io]: tempData,
      }),
    });
  }

  handleSubmit() {
    this.props.resetIO(this.state.data);
    this.props.saveConfigs(this.props.section, this.state.data);
  }

  handleTest(obj) {
    this.props.testModbus(obj.io, obj.bit);
  }
  validateData(data = this.state.data) {
    // return Object.keys(data).every(io => data[io].every(item => item.label && item.function));
    return true
  }
  generatorItems(data, t) {
    const { classes, ioTestRsp } = this.props;
    const { btnGroupStatus } = this.state;
    return data.map((item, idx) => {
      const options = get(
        this.IO_FUNCTION,
        String.prototype.toLowerCase.call(item.io),
        this.IO_FUNCTION.IN,
      );
      const selectItems = options.map(v => (<MenuItem key={v} value={v}>{v}</MenuItem>));

      const generatorStatus = () => {
        if (item.io === 'out') return null;

        if (ioTestRsp[item.bit] === 0) {
          return (
            <div className={classes.statusWrap}>
              <span className={`${classes.statusCircle} ${classes.info}`} />
              <span className={classes.infoText}>OFF</span>
            </div>
          );
        } else if (ioTestRsp[item.bit] === 1) {
          return (
            <div className={classes.statusWrap}>
              <span className={`${classes.statusCircle} ${classes.success}`} />
              <span className={classes.successText}>ON</span>
            </div>
          );
        } else if (ioTestRsp[item.bit] === 'fail') {
          return (
            <div className={classes.statusWrap}>
              <span className={`${classes.statusCircle} ${classes.fail}`} />
              <span className={classes.failText}>Fail</span>
            </div>
          );
        }
        return null;
      };

      return (
        <div key={`${item.io}_${item.bit}`}>
          <ListItem className={classes.inputItem}>
            <InputLabel
              className={classes.ioInputLabel}
              htmlFor="name-simple"
            >
              {item.bit}
            </InputLabel>
            <Input
              id="name-simple"
              placeholder={t('Common.isRequired')}
              className={classes.ioInput}
              value={item.label}
              onChange={e => this.handleChange(e, idx, item.io, 'label')}
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
            <span>
              {String.prototype.toUpperCase.call(item.io)}
            </span>
            <Button
              color="info"
              onClick={() => this.handleTest(item)}
              className={classes.testButton}
              disabled={!get(btnGroupStatus, `${item.io}.${item.bit}`, [])}
            >
              {t('Common.Test')}
            </Button>
            {generatorStatus()}
          </ListItem>
          <li>
            <Divider />
          </li>
        </div>
      );
    });
  }

  render() {
    const { classes } = this.props;
    const { data, isDataValid } = this.state;
    const inItems = t => this.generatorItems(data.in, t);
    const outItems = t => this.generatorItems(data.out, t);

    return (
      <I18n ns="translations">
        {
          t => (
            <section>
              <h3 className={classes.sectionTitle}>
                {t('Configuration.IO.INname')}
              </h3>
              <Paper className={classes.paperWrap} elevation={1}>
                <List>
                  {inItems(t)}
                </List>
              </Paper>
              <h3 className={`${classes.sectionTitle} ${classes.sectionTitleInner}`}>
                {t('Configuration.IO.OUTname')}
              </h3>
              <Paper className={classes.paperWrap} elevation={1}>
                <List>
                  {outItems(t)}
                </List>
                <Button
                  variant="contained"
                  disabled={!isDataValid}
                  color="primary"
                  onClick={this.handleSubmit}
                  className={classes.button}
                >
                  <SaveIcon className={classes.leftIcon} />
                  {t('Common.Submit')}
                </Button>
              </Paper>
            </section>
          )
        }
      </I18n>
    );
  }
}

ConnectedIo.propTypes = {
  classes: PropTypes.shape({
  }).isRequired,
  storedConfigs: PropTypes.shape({
  }).isRequired,
  section: PropTypes.string.isRequired,
  ioTestRsp: PropTypes.shape({
  }).isRequired,
  saveConfigs: PropTypes.func.isRequired,
  testModbus: PropTypes.func.isRequired,
};

const Io = connect(mapStateToProps, mapDispatchToProps)(ConnectedIo);

export default withStyles(styles)(Io);
