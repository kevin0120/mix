import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { isEqual, isNil } from 'lodash';

import Button from '@material-ui/core/Button';
import Paper from '@material-ui/core/Paper';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import { I18n } from 'react-i18next';

import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';

import { sortObj } from '../../common/utils';
import styles from './styles';

const mapStateToProps = (state, ownProps) => ({
  connInfo: state.connections,
  ...ownProps
});

const mapDispatchToProps = {};

function handleTest(obj) {
  obj.test(obj.value);
}

/* eslint-disable react/prefer-stateless-function */
class ConnectedTest extends React.PureComponent {
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

    const { connInfo } = this.props;

    // this.handleTest = this.handleTest.bind(this);
    this.testMasterPC = this.testMasterPC.bind(this);
    this.testModbus = this.testModbus.bind(this);
    this.testAiis = this.testAiis.bind(this);

    this.state = {
      testStatus: {
        masterPcUrl: 99,
        ioUrl: 99,
        aiisUrl: 99
      },
      data: this.formatConnInfo(connInfo),
      btnGroupStatus: {}
    };
  }

  // 重新获取 props 时修改 btnGroupStatus
  // https://github.com/reactjs/reactjs.org/issues/721
  //
  // static getDerivedStateFromProps(props, currentState) {
  //   console.log('getDerivedStateFromProps', props, currentState);

  //   if (currentState.value !== props.value) {
  //     return {
  //       value: props.value,
  //       btnGroupStatus: getBtns(props.value)
  //     }
  //   }
  // }

  componentDidMount() {
    this.setBtnsStatus();
  }

  componentDidUpdate() {
    this.setBtnsStatus();
    const { connInfo } = this.props;
    const {data} = this.state;
    const formatedData = this.formatConnInfo(connInfo);

    if (isEqual(formatedData, data)) return;

    this.updateState(connInfo);
  }

  // 设置 btns 的状态集
  setBtnsStatus(tData) {
    const {btnGroupStatus, data} = this.state;
    const testData = arguments.length > 0 ? tData : data;
    if (isEqual(btnGroupStatus, ConnectedTest.getBtnStatus(testData))) return;

    this.setState({
      btnGroupStatus: ConnectedTest.getBtnStatus(testData),
    });
  }

  updateState(connInfo) {
    this.setState({
      data: this.formatConnInfo(connInfo),
    });
  }

  formatConnInfo(connInfo) {
    const { aiisUrl } = this.props;
    return {
      masterPcUrl: {
        displayOrder: 0,
        value: String(connInfo.masterpc),
        displayTitle: 'MasterPC URL',
        test: this.testMasterPC
      },
      aiisUrl: {
        displayOrder: 50,
        value: String(aiisUrl),
        displayTitle: 'Aiis服务 URL',
        test: this.testAiis
      },
      controllerSn: {
        displayOrder: 100,
        value: isNil(connInfo.controllers[0].serial_no) ? '' : String(connInfo.controllers[0].serial_no),
        displayTitle: '控制器序列号',
      },
      rfidUrl: {
        displayOrder: 200,
        value: String(connInfo.rfid),
        displayTitle: 'RFID 链接地址'
      },
      ioUrl: {
        displayOrder: 300,
        value: String(connInfo.io),
        displayTitle: 'IO 模块链接地址'
        // test: this.testModbus,
      }
    };
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
            aiisUrl: false,
          }
        });
      })
      .catch(() => {
        this.setState({
          testStatus: {
            ...testStatus,
            aiisUrl: false,
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
            masterPcUrl: response.status === 204,
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
        ioUrl: true,
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
                size="small"
                color="primary"
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
                disabled
                placeholder={t('Common.isRequired')}
                className={classes.input}
                value={item.value}
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
          </Paper>
        )}
      </I18n>
    );
  }
}

ConnectedTest.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  connInfo: PropTypes.shape({}).isRequired
};

const Test = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedTest);

export default withStyles(styles)(Test);
