import React, { useEffect, useState } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { makeStyles } from '@material-ui/core/styles';
import { get, cloneDeep } from 'lodash';

import Paper from '@material-ui/core/Paper';
import Divider from '@material-ui/core/Divider';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import SaveIcon from '@material-ui/icons/Save';
import { I18n } from 'react-i18next';
import Dialog from '@material-ui/core/Dialog';
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import Fade from '@material-ui/core/Fade';

import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';

import { css } from '@emotion/core';

// First way to import
import { GridLoader } from 'react-spinners';
import Button from '../../components/CustomButtons/Button';

import { sortObj } from '../../common/utils';
import saveConfigs from '../../modules/setting/action';
import { networkScan, networkCheck, networkSet } from '../../modules/network/action';

import styles from './styles';
import withKeyboard from '../../components/Keyboard';

const lodash = require('lodash');

const override = css`
  display: block;
  margin: auto;
  border-color: red;
`;

function renderSSIDs(ssid, index, classes) {
  return (
    <MenuItem
      key={ssid}
      classes={{
        root: classes.selectMenuItem,
        selected: classes.selectMenuItemSelected
      }}
      value={ssid}
    >
      {ssid}
    </MenuItem>
  );
}

const mapStateToProps = (state, ownProps) => ({
  storedConfigs: state.setting.page.network,
  network: state.network,
  ...ownProps
});

const mapDispatchToProps = {
  saveConfigs,
  doNetworkScan: networkScan,
  doNetworkCheck: networkCheck,
  doNetworkSet: networkSet
};

const filter = (value, defaultValue) => value === undefined ? defaultValue : value;

function ConnectedNet({ doNetworkScan, doNetworkCheck, doNetworkSet, network, keyboardInput }) {
  const [data, setData] = useState({});
  const [ssidSelectOpen, setSSIDSelectOpen] = useState(false);

  useEffect(() => {
    doNetworkCheck();
  }, [doNetworkCheck]);

  function handleChange(e, key) {
    const tempData = cloneDeep(data);
    tempData[key].value = get(e, 'target.value', '').trim();
    setData(tempData);
  }

  function handleSubmit() {
    const tempData = cloneDeep(network.config);
    Object.keys(data).forEach((key) => {
      tempData[key].value = data[key];
    });
    doNetworkSet(tempData);
  }

  function handleChangeSSID(e) {
    setData({
      ...data,
      ssid: e.target.value
    });
  }

  function validateData(d = data) {
    if (lodash.isEmpty(d.ssid)) {
      return false;
    }
    return Object.values(d.data).every(v => v.value);
  }

  const getSSIDs = () => {
    doNetworkScan();
    setSSIDSelectOpen(true);
  };

  const handleCloseSSID = () => {
    setSSIDSelectOpen(true);
  };


  const Transition = React.forwardRef((p, r) => <Fade {...p} timeout={500} ref={r}/>);

  const classes = makeStyles(styles.content)();

  const dataToValidate = lodash.omit(data, ['ssid']);

  const submitDisabled = Object.values(dataToValidate).some(v => v === '');

  const inputsItems = t =>
    sortObj(network.config, 'displayOrder')
      .slice(1)
      .map(({ key, value: item }) => (
        <div key={key}>
          <ListItem className={classes.inputItem}>
            <InputLabel className={classes.inputLabel} htmlFor="name-simple">
              {t(item.displayTitle)}
            </InputLabel>
            <Input
              id="name-simple"
              type={item.isPWD ? 'password' : null}
              placeholder={t('Common.isRequired')}
              className={classes.input}
              value={filter(data[key], item.value)}
              onClick={() => {
                keyboardInput({
                  onSubmit: text => {
                    const tempData = cloneDeep(data);
                    tempData[key] = text;
                    setData(tempData);
                  },
                  text: filter(data[key], item.value),
                  title: item.displayTitle,
                  label: item.displayTitle
                });
              }}
              // onChange={e => handleChange(e, key)}
            />
          </ListItem>
          <li>
            <Divider/>
          </li>
        </div>
      ));

  return (
    <I18n ns="translations">
      {t => (
        <div>
          <Dialog
            fullScreen
            classes={{
              root: classes.loadModal
            }}
            open={network.connecting}
            style={{ opacity: 0.7 }}
            TransitionComponent={Transition}
          >
            <GridLoader
              className={override}
              sizeUnit="px"
              size={50}
              color="#36D7B7"
              loading={network.connecting}
            />
          </Dialog>
          <section className={classes.section}>
            <h3 className={classes.sectionTitle}>
              {t('Configuration.network.name')}
            </h3>
            <Paper className={classes.paperWrap} elevation={1}>
              <List>
                <div>
                  <ListItem className={classes.inputItem}>
                    <InputLabel className={classes.inputLabel} htmlFor="ssid">
                      {t('Configuration.network.SSID')}
                    </InputLabel>
                    <Select
                      displayEmpty
                      MenuProps={{
                        className: classes.selectMenu
                      }}
                      classes={{
                        select: classes.select
                      }}
                      value={filter(data.ssid, network.config.ssid.value)}
                      onChange={e => handleChangeSSID(e)}
                      open={ssidSelectOpen}
                      onOpen={getSSIDs}
                      onClose={handleCloseSSID}
                      renderValue={(v) => v}
                      inputProps={{
                        name: 'ssid',
                        id: 'ssid',
                        className: classes.input
                      }}
                    >
                      <MenuItem
                        disabled
                        classes={{
                          root: classes.selectMenuItem
                        }}
                      >
                        {t('Configuration.network.SSID')}
                      </MenuItem>
                      {network.ssidList.map((item, idx) =>
                        renderSSIDs(item, idx, classes)
                      )}
                    </Select>
                    {/* <Input */}
                    {/* id="ssid" */}
                    {/* placeholder={t('Common.isRequired')} */}
                    {/* className={classes.input} */}
                    {/* value={state.ssid} */}
                    {/* onChange={e => handleChangeSSID(e)} */}
                    {/* /> */}
                  </ListItem>
                  <li>
                    <Divider/>
                  </li>
                </div>
                {inputsItems(t)}
              </List>
              <Button
                size="lg"
                disabled={submitDisabled}
                color="info"
                onClick={handleSubmit}
                className={classes.button}
              >
                <SaveIcon className={classes.leftIcon}/>
                {t('Common.Submit')}
              </Button>
            </Paper>
          </section>
        </div>
      )}
    </I18n>
  );
}

ConnectedNet.propTypes = {
  storedConfigs: PropTypes.shape({}).isRequired,
  saveConfigs: PropTypes.func.isRequired
};

const Net = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedNet);

export default withKeyboard(Net);
