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
import saveConfigs from '../../actions/userConfigs';

import styles from './styles';

const lodash = require('lodash');

const override = css`
  display: block;
  margin: auto;
  border-color: red;
`;

const netmask2CIDR = netmask =>
  netmask
    .split('.')
    .map(Number)
    .map(part => (part >>> 0).toString(2))
    .join('')
    .split('1').length - 1;
//
// const CDIR2netmask = (bitCount) => {
//   let mask=[];
//   for(let i=0;i<4;i++) {
//     let n = Math.min(bitCount, 8);
//     mask.push(256 - Math.pow(2, 8-n));
//     bitCount -= n;
//   }
//   return mask.join('.');
// };

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
  ...ownProps
});

const mapDispatchToProps = {
  saveConfigs
};

class ConnectedNet extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      loading: false,
      section: 'network',
      data: props.storedConfigs,
      ssid: '',
      ssidSelectOpen: false,
      ssids: []
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleChange(e, key) {
    const tempData = cloneDeep(this.state.data);
    tempData[key].value = get(e, 'target.value', '').trim();
    this.setState({
      data: tempData
    });
  }

  handleSubmit() {
    const { section, ssid, data } = this.state;
    const { saveConfigs } = this.props;
    this.setState({
      loading: true
    });
    let ret = 0;
    const tempData = cloneDeep(data);
    const mask = netmask2CIDR(tempData.netmask.value);
    tempData.ssid.value = ssid;
    const { exec } = require('child_process');
    exec('nmcli con delete default', () => {
      exec(
        `nmcli dev wifi connect ${tempData.ssid.value} password ${
          tempData.password.value
        } name default`,
        (error, stdout, stderr) => {
          if (error) {
            console.error(`exec error: ${error}`);
            ret = -1;
          }
          if (stderr) {
            console.log(`stderr: ${stderr}`);
            ret = -1;
          }

          if (ret < 0) {
            this.setState({
              loading: false
            });
            return;
          }
          exec('nmcli con down default', () => {
            const cmd = `nmcli con mod default ipv4.method manual ipv4.address ${
              tempData.ipAddress.value
            }/${mask} ipv4.gateway ${tempData.gateway.value}`;
            exec(cmd, (error, stdout, stderr) => {
              if (error) {
                console.error(`exec error: ${error}`);
                ret = -1;
              }
              if (stderr) {
                console.log(`stderr: ${stderr}`);
                ret = -1;
              }
              this.setState({
                loading: false
              });
              if (ret === 0) {
                saveConfigs(section, tempData);
                exec('nmcli con up default');
              }
            });
          });
        }
      );
    });
  }

  handleChangeSSID(e) {
    this.setState({
      ssid: e.target.value
    });
  }

  validateData(data = this.state) {
    if (lodash.isEmpty(data.ssid)) {
      return false;
    }
    return Object.values(data.data).every(v => v.value);
  }

  getSSIDs = () => {
    const ret = [];
    const { exec } = require('child_process');
    exec('nmcli dev wifi', (error, stdout, stderr) => {
      if (error) {
        console.error(`exec error: ${error}`);
        return;
      }
      if (stdout) {
        const lines = stdout.toString().split('\n');
        let isHeader = true;
        for (let i = 0; i < lines.length - 1; i++) {
          if (isHeader) {
            isHeader = false;
          } else {
            const line = lines[i].toString();
            const x = lodash.words(line, /[^, ]+/g);
            if (x[0] === '*') {
              ret.push(x[1]);
            } else {
              ret.push(x[0]);
            }
          }
        }
        this.setState({
          ssids: lodash.uniq(ret),
          ssidSelectOpen: true
        });
      }
      if (stderr) {
        this.setState({
          ssids: [],
          ssidSelectOpen: true
        });
      }
    });
  };

  handleCloseSSID = () => {
    this.setState({
      ssidSelectOpen: false
    });
  };

  componentDidMount() {
    const { data } = this.state;
    const tempData = cloneDeep(this.state);
    const ret = [];
    const { exec } = require('child_process');
    exec('nmcli dev wifi', (error, stdout, stderr) => {
      if (error) {
        console.error(`exec error: ${error}`);
        return;
      }
      if (stdout) {
        const lines = stdout.toString().split('\n');
        let isHeader = true;
        for (let i = 0; i < lines.length - 1; i += 1) {
          if (isHeader) {
            isHeader = false;
          } else {
            const line = lines[i].toString();
            const x = lodash.words(line, /[^, ]+/g);
            if (x[0] === '*') {
              ret.push(x[1]);
            } else {
              ret.push(x[0]);
            }
          }
        }
        tempData.ssids = lodash.uniq(ret);
        if (lodash.includes(ret, data.ssid.value)) {
          tempData.ssid = data.ssid.value;
        }
        this.setState({
          ...tempData
        });
      }
    });
  }

  Transition(props) {
    return <Fade {...props} timeout={500} />;
  }

  render() {
    const { classes } = this.props;
    const { data, ssid, loading, ssidSelectOpen, ssids } = this.state;

    const validateData = lodash.omit(data, ['ssid']);

    const submitDisabled =
      ssid === '' || Object.values(validateData).some(v => v.value === '');

    const inputsItems = t =>
      sortObj(data, 'displayOrder')
        .slice(1)
        .map(({ key, value: item }) => (
          <div key={key}>
            <ListItem className={classes.inputItem}>
              <InputLabel className={classes.inputLabel} htmlFor="name-simple">
                {t(item.displayTitle)}
              </InputLabel>
              <Input
                id="name-simple"
                placeholder={t('Common.isRequired')}
                className={classes.input}
                value={item.value}
                onChange={e => this.handleChange(e, key)}
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
          <div>
            <Dialog
              fullScreen
              classes={{
                root: classes.loadModal
              }}
              open={loading}
              style={{ opacity: 0.7 }}
              TransitionComponent={this.Transition}
            >
              <GridLoader
                className={override}
                sizeUnit="px"
                size={50}
                color="#36D7B7"
                loading={loading}
              />
            </Dialog>
            <section>
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
                        value={ssid}
                        onChange={e => this.handleChangeSSID(e)}
                        open={ssidSelectOpen}
                        onOpen={this.getSSIDs}
                        onClose={this.handleCloseSSID}
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
                        {ssids.map((item, idx) =>
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
                      <Divider />
                    </li>
                  </div>
                  {inputsItems(t)}
                </List>
                <Button
                  disabled={submitDisabled}
                  color="primary"
                  onClick={this.handleSubmit}
                  className={classes.button}
                >
                  <SaveIcon className={classes.leftIcon} />
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

export default withStyles(styles)(Net);
