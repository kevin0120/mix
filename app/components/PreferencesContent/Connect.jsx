import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';
import { get, cloneDeep } from 'lodash';

import Paper from '@material-ui/core/Paper';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import SaveIcon from '@material-ui/icons/Save';
import { I18n } from 'react-i18next';
import Button from '../CustomButtons/Button';

import saveConfigs from '../../actions/userConfigs';
import { systemInit } from '../../actions/sysInit';

import {toggleRFID} from '../../actions/rfid';

import { sortObj, defaultClient } from '../../common/utils';
import Test from './Test';
import styles from './styles';
import withKeyboard from '../Keyboard';

const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  storedConfigs: state.setting.page.odooConnection,
  connInfo: state.setting.system.connections,
  rfidEnabled: state.setting.systemSettings.rfidEnabled,
  ...ownProps
});

const mapDispatchToProps = {
  saveConfigs,
  systemInit,
  toggleRFID
};

/* eslint-disable react/prefer-stateless-function */
class ConnectedConnect extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isDataValid: true,
      data: props.storedConfigs,
      section: 'odooConnection',
      connInfoData: props.connInfo
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e, key) {
    const { data } = this.state;
    const tempData = cloneDeep(data);
    tempData[key].value = get(e, 'target.value', '').trim();
    this.setState({
      data: tempData,
      isDataValid: this.validateData(tempData)
    });
  }

  handleTestKeyBoardSubmit = (key, text) => {
    const { connInfoData } = this.state;
    this.setState({
      connInfoData: {
        ...connInfoData,
        [key]: text
      }
    });
  };

  handleSubmit() {
    const { saveConfigs, connInfo } = this.props;
    const { section, data } = this.state;
    const fullUrl = `${data.odooUrl.value}/hmi.connections/${data.hmiSn.value}`;
    defaultClient
      .get(fullUrl)
      .then(resp => {
        const { masterpc, rfid, io, controllers, info } = resp.data;
        const d = {
          masterpc: masterpc.connection ? masterpc.connection : '',
          aiis: connInfo.aiis,
          rfid: rfid.connection ? rfid.connection : '',
          io: io.connection ? io.connection : '',
          workcenterCode: info.workcenter_code ? info.workcenter_code : '',
          rework_workcenter: info.qc_workcenter ? info.qc_workcenter : '',
          controllers: lodash.isArray(controllers) ? controllers : []
        };
        this.setState({
          connInfoData: d
        });
      })
      .catch(e => console.log(e.toString()));
    saveConfigs(section, data);
  }

  validateData(data = this.state.data) {
    return Object.values(data).every(v => v.value);
  }

  render() {
    const { classes, systemInit, saveConfigs,toggleRFID,rfidEnabled } = this.props;
    const { data, isDataValid, connInfoData } = this.state;

    const baseItems = t =>
      sortObj(data, 'displayOrder').map(({ key, value: item }) => (
        <div key={key}>
          <ListItem className={classes.inputItem}>
            <InputLabel className={classes.inputLabel} htmlFor={key}>
              {t(item.displayTitle)}
            </InputLabel>
            <Input
              id={key}
              placeholder={t('Common.isRequired')}
              className={classes.input}
              value={item.value}
              // onChange={e => this.handleChange(e, key)}
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
          </ListItem>
          <li>
            <Divider />
          </li>
        </div>
      ));

    return (
      <I18n ns="translations">
        {t => (
          <section className={classes.section}>
            <h3 className={classes.sectionTitle}>
              {t('Configuration.connections.name')}
            </h3>
            <Paper className={classes.paperWrap} elevation={1}>
              <List>{baseItems(t)}</List>
              <Button
                disabled={!isDataValid}
                color="info"
                onClick={this.handleSubmit}
                className={classes.button}
              >
                <SaveIcon className={classes.leftIcon} />
                {t('Common.Submit')}
              </Button>
            </Paper>
            <h3
              className={`${classes.sectionTitle} ${classes.sectionTitleInner}`}
            >
              {t('Common.Test')}
            </h3>
            <Paper className={classes.paperWrap} elevation={1}>
              <Test
                connInfoData={connInfoData}
                rfidEnabled={rfidEnabled}
                systemInit={systemInit}
                saveConfigs={saveConfigs}
                toggleRFID={toggleRFID}
                keyBoardSubmit={this.handleTestKeyBoardSubmit}
              />
            </Paper>
          </section>
        )}
      </I18n>
    );
  }
}

ConnectedConnect.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  storedConfigs: PropTypes.shape({}).isRequired,
  saveConfigs: PropTypes.func.isRequired,
  rfidEnabled: PropTypes.bool.isRequired,
  toggleRFID: PropTypes.func.isRequired,
  connInfo: PropTypes.shape({}).isRequired,
  keyboardInput: PropTypes.func.isRequired
};

const Connect = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedConnect);

export default withKeyboard(withStyles(styles)(Connect));
