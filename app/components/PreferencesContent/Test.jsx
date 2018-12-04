import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { isEqual, isNil, cloneDeep } from 'lodash';

import Paper from '@material-ui/core/Paper';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import { I18n } from 'react-i18next';

import SaveIcon from '@material-ui/icons/Save';

import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import Button from '../CustomButtons/Button';

import { sortObj } from '../../common/utils';
import styles from './styles';
import saveConfigs from "../../actions/userConfigs";
import withKeyboard from "../Keyboard";
import { systemInit } from "../../actions/sysInit";


function handleTest(obj) {
  obj.test(obj.value);
}

/* eslint-disable react/prefer-stateless-function */
class ConnectedTest extends React.Component {
  // 获取 btns 的状态集
  static getBtnStatus(data) {
    return Object.keys(data).reduce(
      (pre, key) => ({
        ...pre,
        [key]: data[key].value === 0 || !!data[key].value
      }),
      {}
    );
  }

  constructor(props) {
    super(props);
    console.log(this.props, props);
    // this.handleTest = this.handleTest.bind(this);
    this.testMasterPC = this.testMasterPC.bind(this);
    this.testModbus = this.testModbus.bind(this);
    this.testAiis = this.testAiis.bind(this);
    this.handleSave = this.handleSave.bind(this);

    this.state = {
      testStatus: {
        masterPcUrl: 99,
        ioUrl: 99,
        aiisUrl: 99
      },
      data: this.formatConnInfo(this.props.connInfoData),
      btnGroupStatus: {},
    };
  }


  componentDidMount() {
    console.log(this.props);
    this.setState({data: this.formatConnInfo(this.props.connInfoData)});
    this.setBtnsStatus();
  }

  componentDidUpdate() {
    this.setBtnsStatus();
    const { connInfoData } = this.props;
    const { data } = this.state;
    const formatedData = this.formatConnInfo(connInfoData);

    if (isEqual(formatedData, data)) return;

    this.updateState(connInfoData);
  }

  // 设置 btns 的状态集
  setBtnsStatus(tData) {
    const { btnGroupStatus, data } = this.state;
    const testData = arguments.length > 0 ? tData : data;
    if (isEqual(btnGroupStatus, ConnectedTest.getBtnStatus(testData))) return;

    this.setState({
      btnGroupStatus: ConnectedTest.getBtnStatus(testData)
    });
  }

  updateState(connInfo) {
    this.setState({
      data: this.formatConnInfo(connInfo)
    });
  }

  formatConnInfo(connInfo) {
    return {
      masterPcUrl: {
        key: 'masterpc',
        displayOrder: 0,
        value: String(connInfo.masterpc),
        displayTitle: 'MasterPC URL',
        test: this.testMasterPC
      },
      aiisUrl: {
        key: 'aiis',
        displayOrder: 50,
        value: String(connInfo.aiis),
        displayTitle: 'Aiis服务 URL',
        test: this.testAiis
      },
      controllerSn: {
        key: 'controllers',
        displayOrder: 100,
        value: isNil(connInfo.controllers[0])
          ? ''
          : String(connInfo.controllers[0].serial_no),
        displayTitle: '控制器序列号'
      },
      rfidUrl: {
        key: 'rfid',
        displayOrder: 200,
        value: String(connInfo.rfid),
        displayTitle: 'RFID 链接地址'
      },
      ioUrl: {
        key: 'io',
        displayOrder: 300,
        value: String(connInfo.io),
        displayTitle: 'IO 模块链接地址'
        // test: this.testModbus,
      },
      workCenterCode: {
        key: 'workcenter',
        displayOrder: 200,
        value: String(connInfo.workcenterCode),
        displayTitle: '工位'
      },
      rework: {
        key: 'rework',
        displayOrder: 200,
        value: String(connInfo.rework_workcenter),
        displayTitle: '返修工位'
      },
    };
  }

  handleSave() {
    const { saveConfigs, systemInit } = this.props;
    const { data } = this.state;
    const section='connections';
    saveConfigs(section, {masterpc: data.masterPcUrl.value,
      aiis: data.aiisUrl.value,
      controllers: [{serial_no: data.controllerSn.value}],
      rfid: data.rfidUrl.value,
      io: data.ioUrl.value,
      workcenterCode: data.workCenterCode.value,
      rework_workcenter: data.rework.value,
    });
    systemInit()
  }


  testAiis(conn) {
    const url = `${conn}/aiis/v1/healthz`;
    const { testStatus } = this.state;
    fetch(url, {
      timeout: 3000
    })
      .then(() => {
        this.setState({
          testStatus: {
            ...testStatus,
            aiisUrl: false
          }
        });
      })
      .catch(() => {
        this.setState({
          testStatus: {
            ...testStatus,
            aiisUrl: false
          }
        });
      });
  }

  testMasterPC(conn) {
    const url = `${conn}/rush/v1/healthz`;
    const { testStatus } = this.state;
    fetch(url, {
      timeout: 3000
    })
      .then(response => {
        this.setState({
          testStatus: {
            ...testStatus,
            masterPcUrl: response.status === 204
          }
        });
      })
      .catch(() => {
        this.setState({
          testStatus: {
            ...testStatus,
            masterPcUrl: false
          }
        });
      });
  }

  testModbus() {
    const { testStatus } = this.state;
    this.setState({
      testStatus: {
        ...testStatus,
        ioUrl: true
      }
    });
  }

  render() {
    const { classes } = this.props;
    const { data, btnGroupStatus, testStatus } = this.state;
    const inputsItems = t =>
      sortObj(data, 'displayOrder').map(({ key, value: item }) => {
        const testPart = t =>
          item.test ? (
            <div>
              <Button
                variant="outlined"
                disabled={!btnGroupStatus[key]}
                size="sm"
                color="warning"
                onClick={() => handleTest(item)}
                className={classes.testButton}
              >
                {t('Common.Test')}
              </Button>
              {testStatus[key] !== 99 ? (
                <span
                  className={`${classes.statusCircle} ${
                    testStatus[key] ? classes.success : classes.fail
                  }`}
                />
              ) : null}
            </div>
          ) : null;

        return (
          <div key={key}>
            <ListItem className={classes.inputItem}>
              <InputLabel className={classes.inputLabel}>
                {item.displayTitle}
              </InputLabel>
              <Input
                key={item.key}
                placeholder={t('Common.isRequired')}
                className={classes.input}
                value={item.value}
                onClick={() => {
                  this.props.keyboardInput({
                    onSubmit: text => {
                      const tempData = cloneDeep(this.state.data);
                      tempData[key].value = text;
                      this.setState({
                        data: tempData
                      });
                    },
                    text: item.value,
                    title: item.displayTitle,
                    label: item.displayTitle
                  });
                }}
              />
              {testPart(t)}
            </ListItem>
            <li>
              <Divider />
            </li>
          </div>
        );
      });

    return (
      <I18n ns="translations">
        {t => (
          <Paper className={classes.paperWrap} elevation={1}>
            <List>{inputsItems(t)}</List>
            <Button
              color="info"
              onClick={this.handleSave}
              className={classes.button}
            >
              <SaveIcon className={classes.leftIcon} />
              {t('Common.Save')}
            </Button>
          </Paper>
        )}
      </I18n>
    );
  }
}

ConnectedTest.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  // connInfoData: PropTypes.shape({}).isRequired,
};


export default withKeyboard(withStyles(styles)(ConnectedTest));
