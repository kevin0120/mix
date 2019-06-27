import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';
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

import { css } from 'react-emotion';

// First way to import
import { GridLoader } from 'react-spinners';
import Button from '../CustomButtons/Button';

import { sortObj } from '../../common/utils';
import saveConfigs from '../../modules/setting/action';
import { networkScan, networkCheck, networkSet } from '../../modules/network/action';

import styles from './styles';
import withKeyboard from '../Keyboard';

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

class ConnectedNet extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      data: {},
      ssidSelectOpen: false
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    const { doNetworkCheck } = this.props;
    doNetworkCheck();
  }

  handleChange(e, key) {
    const { data } = this.state;
    const tempData = cloneDeep(data);
    tempData[key].value = get(e, 'target.value', '').trim();
    this.setState({
      data: tempData
    });
  }

  handleSubmit() {
    const { data } = this.state;
    const { doNetworkSet, network } = this.props;
    const tempData = cloneDeep(network.config);
    Object.keys(data).forEach((key) => {
      tempData[key].value = data[key];
    });
    doNetworkSet(tempData);
  }

  handleChangeSSID(e) {
    const { data } = this.state;
    this.setState({
      data: {
        ...data,
        ssid: e.target.value
      }
    });
  }

  validateData(data = this.state) {
    if (lodash.isEmpty(data.ssid)) {
      return false;
    }
    return Object.values(data.data).every(v => v.value);
  }

  getSSIDs = () => {
    const { doNetworkScan } = this.props;
    doNetworkScan();
    this.setState({
      ssidSelectOpen: true
    });
  };

  handleCloseSSID = () => {
    this.setState({
      ssidSelectOpen: false
    });
  };


  Transition(props) {
    return <Fade {...props} timeout={500}/>;
  }

  render() {
    const { classes, network, keyboardInput } = this.props;
    const { data, ssidSelectOpen } = this.state;

    const validateData = lodash.omit(data, ['ssid']);

    const submitDisabled =
      Object.values(validateData).some(v => v === '');



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
                      this.setState({
                        data: tempData
                      });
                    },
                    text: filter(data[key], item.value),
                    title: item.displayTitle,
                    label: item.displayTitle
                  });
                }}
                // onChange={e => this.handleChange(e, key)}
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
              TransitionComponent={this.Transition}
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
                        onChange={e => this.handleChangeSSID(e)}
                        open={ssidSelectOpen}
                        onOpen={this.getSSIDs}
                        onClose={this.handleCloseSSID}
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
                      {/* value={this.state.ssid} */}
                      {/* onChange={e => this.handleChangeSSID(e)} */}
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
                  onClick={this.handleSubmit}
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
}

ConnectedNet.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  storedConfigs: PropTypes.shape({}).isRequired,
  saveConfigs: PropTypes.func.isRequired
};

const Net = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedNet);

export default withKeyboard(withStyles(styles)(Net));
