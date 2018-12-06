import React from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { connect } from 'react-redux';

import Button from '../../components/CustomButtons/Button';
import { withStyles } from '@material-ui/core/styles';

import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogActions from '@material-ui/core/DialogActions';
import Dialog from '@material-ui/core/Dialog';
import Slide from '@material-ui/core/Slide';

import { NewCar } from '../../actions/scannerDevice';

import customSelectStyle from '../../common/jss/customSelectStyle';

import { I18n } from 'react-i18next';
import InputAdornment from '@material-ui/core/InputAdornment';
import Face from '@material-ui/icons/Face';
import CustomInput from '../CustomInput/CustomInput';
import Keyboard from 'react-simple-keyboard';

const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  ...ownProps
});

const mapDispatchToProps = {
  NewCar
};

const customerStyles = theme => ({
  ...customSelectStyle(theme),
  bigInput: {
    '&,&::placeholder': {
      fontSize: '50px'
    }
  }
});

/* eslint-disable react/prefer-stateless-function */
class ConnectedManualDiag extends React.Component {
  constructor(props) {
    super(props);
    this.keyboard = null;

    this.state = {
      vin: '',
      layoutName: 'default'
    };
    this.handleClose = this.handleClose.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.Transition = this.Transition.bind(this);
  }

  handleClose = () => {
    this.props.close();
  };

  handleSubmit = () => {
    console.log('manual scanner');
    this.props.NewCar(this.state.vin);
    this.props.close();
  };

  setActiveInput = event => {
    // this.setState({
    //   vin: event.target.id
    // });
  };

  handleInput = input => {
    // console.log(input);
    const { vin } = this.state;

    this.setState({
      ...vin,
      vin: input
    });
  };

  handlePress = press => {
    const { vin } = this.state;
    if (lodash.isEqual(press, '{enter}') && vin.length !== 0) {
      this.handleSubmit();
    }
    if (lodash.isEqual(press, '{shift}')) {
      const layoutName = this.state.layoutName;

      this.setState({
        layoutName: layoutName === 'default' ? 'shift' : 'default'
      });
    }
  };

  handleChangeInputValue = event => {
    this.setState({
      vin: event.target.value
    });
    this.keyboard.setInput(event.target.value, 'vin');
  };

  Transition(props) {
    return <Slide direction="down" {...props} />;
  }

  render() {
    const { show, classes } = this.props;

    const { vin } = this.state;

    const submitEnable = vin.length !== 0;

    return (
      <I18n ns="translations">
        {t => (
          <Dialog
            classes={{
              root: classes.modalRoot,
              paper: classes.modal + ' ' + classes.modalLarge
            }}
            TransitionComponent={this.Transition}
            keepMounted
            open={show}
            onClose={this.handleClose}
            aria-labelledby="form-dialog-title"
            scroll="paper"
          >
            <DialogTitle id="form-dialog-title" className={classes.modalHeader}>
              {t('ManualDiag.title')}
            </DialogTitle>
            <DialogContent className={classes.modalBody}>
              <CustomInput
                large
                labelText={t('VIN/KNR')}
                id="vin"
                formControlProps={{
                  fullWidth: true,
                  required: true
                }}
                inputProps={{
                  onFocus: this.setActiveInput,
                  required: true,
                  value: vin || '',
                  onChange: this.handleChangeInputValue
                }}
              />
            </DialogContent>
            <DialogActions
              className={classes.modalFooter + ' ' + classes.modalFooterCenter}
            >
              <Button
                onClick={this.handleClose}
                color="primary"
                autoFocus
                round
              >
                {t('Common.Close')}
              </Button>
              <Button
                onClick={this.handleSubmit}
                color="success"
                round
                disabled={!submitEnable}
              >
                {t('Common.Submit')}
              </Button>
            </DialogActions>
            {
              <div className={classes.keyboard}>
                <Keyboard
                  ref={r => (this.keyboard = r)}
                  inputName={'vin'}
                  layoutName={this.state.layoutName}
                  onChange={this.handleInput}
                  onKeyPress={this.handlePress}
                />
              </div>
            }
          </Dialog>
        )}
      </I18n>
    );
  }
}

ConnectedManualDiag.propTypes = {
  show: PropTypes.bool.isRequired,
  close: PropTypes.func.isRequired
};

const ManualDiag = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedManualDiag);

export default withStyles(customerStyles)(ManualDiag);
