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

import { sortObj } from '../../common/utils';
import Test from './Test';
import styles from './styles';
import withKeyboard from '../Keyboard';

const mapStateToProps = (state, ownProps) => ({
  storedConfigs: state.setting.page.odooConnection,
  ...ownProps
});

const mapDispatchToProps = {
  saveConfigs,
  systemInit
};

/* eslint-disable react/prefer-stateless-function */
class ConnectedConnect extends React.PureComponent {
  constructor(props) {
    super(props);
    this.state = {
      isDataValid: true,
      data: props.storedConfigs,
      section: 'odooConnection'
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

  handleSubmit() {
    const { saveConfigs, systemInit, storedConfigs } = this.props;
    const { section, data } = this.state;
    saveConfigs(section, data);
    systemInit(storedConfigs.odooUrl.value, storedConfigs.hmiSn.value);
  }

  validateData(data = this.state.data) {
    return Object.values(data).every(v => v.value);
  }

  render() {
    const { classes } = this.props;
    const { data, isDataValid } = this.state;

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
              onClick={()=>{
                this.props.keyboardInput({
                  onSubmit:(text)=>{
                    const tempData = cloneDeep(this.state.data);
                    tempData[key].value = text;
                    this.setState({
                      ...this.state,
                      data: tempData,
                    });
                  },
                  text:item.value,
                  title:item.displayTitle,
                  label:item.displayTitle
                })
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
          <section>
            <h3 className={classes.sectionTitle}>
              {t('Configuration.connections.name')}
            </h3>
            <Paper className={classes.paperWrap} elevation={1}>
              <List>{baseItems(t)}</List>
              <Button
                disabled={!isDataValid}
                color="primary"
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
              <Test aiisUrl={data.aiisUrl.value} />
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
  systemInit: PropTypes.func.isRequired,
  keyboardInput:PropTypes.func.isRequired
};

const Connect = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedConnect);

export default withKeyboard(withStyles(styles)(Connect));
