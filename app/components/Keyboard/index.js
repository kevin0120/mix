/* eslint-disable react/no-this-in-sfc */
import React from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Keyboard from 'react-simple-keyboard';
import { I18n } from 'react-i18next';
import { withStyles } from '@material-ui/core/styles';
import Slide from '@material-ui/core/Slide';
import InputAdornment from '@material-ui/core/InputAdornment';
import IconButton from '@material-ui/core/IconButton';
import Clear from '@material-ui/icons/CancelOutlined';
import customSelectStyle from '../../common/jss/customSelectStyle';
import Button from '../CustomButtons/Button';
import CustomInput from '../CustomInput/CustomInput';

const lodash = require('lodash');

const customerStyles = theme => ({
  ...customSelectStyle(theme),
  bigInput: {
    '&,&::placeholder': {
      fontSize: '50px'
    }
  }
});

export default function withKeyboard(SubComponents) {
  class KeyboardDialog extends React.Component {
    constructor(props) {
      super(props);
      this.state = {
        layoutName: 'default',
        text: '',
        show: false,
        config: {
          onSubmit: null,
          title: '',
          label: '',
          required: false
        }
      };
    }

    handleSubmit = () => {
      const { text, config } = this.state;
      config.onSubmit(text);
      this.handleClose();
    };

    handleClose = () => {
      const { onClose } = this.state;
      this.setState({ show: false }, () => {
        // eslint-disable-next-line no-unused-expressions
        onClose && onClose();
      });
    };

    onKeyPress = press => {
      const { text } = this.state;
      if (lodash.isEqual(press, '{enter}') && text.length !== 0) {
        this.handleSubmit();
      }
      if (lodash.isEqual(press, '{shift}')) {
        const { layoutName } = this.state;

        this.setState({
          layoutName: layoutName === 'default' ? 'shift' : 'default'
        });
      }
    };

    onChange = text => {
      this.setState({ text });
    };

    bindKeyboardInput = config => {
      this.keyboard.setInput(config.text || '');
      this.setState({
        show: true,
        config,
        text: config.text || ''
      });
    };

    onChangeInput = event => {
      const text = event.target.value;
      this.setState(
        {
          text
        },
        () => {
          this.keyboard.setInput(text);
        }
      );
    };

    Transition = props => <Slide direction="down" {...props} />;

    render() {
      const { classes, ...restProps } = this.props;
      const { text, config, layoutName, show } = this.state;
      const submitEnable = config.required ? text.length !== 0 : true;
      return (
        <div style={{ width: '100%', height: '100%' }}>
          <SubComponents
            keyboardInput={c => this.bindKeyboardInput(c)}
            {...restProps}
          />
          <I18n ns="translations">
            {t => (
              <Dialog
                classes={{
                  root: classes.modalRoot,
                  paper: `${classes.modal} ${classes.modalLarge}`
                }}
                TransitionComponent={this.Transition}
                keepMounted
                open={show}
                onClose={this.handleClose}
                aria-labelledby="form-dialog-title"
                scroll="paper"
              >
                <DialogTitle
                  id="form-dialog-title"
                  className={classes.modalHeader}
                >
                  {t(config.title)}
                </DialogTitle>
                <DialogContent className={classes.modalBody}>
                  <CustomInput
                    large
                    labelText={t(config.label)}
                    id="keyboard-dialog-input"
                    formControlProps={{
                      fullWidth: true,
                      required: config.required
                    }}
                    inputProps={{
                      // onFocus: this.setActiveInput,
                      required: config.required,
                      value: text || '',
                      onChange: this.onChangeInput,
                      endAdornment: (
                        <InputAdornment position="end">
                          <IconButton
                            aria-label="clear input"
                            onClick={() => {
                              this.setState({ text: '' }, () => {
                                this.keyboard.setInput('');
                              });
                            }}
                          >
                            <Clear />
                          </IconButton>
                        </InputAdornment>
                      )
                    }}
                  />
                </DialogContent>
                <DialogActions
                  className={`${classes.modalFooter} ${
                    classes.modalFooterCenter
                  }`}
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
                <div className={classes.keyboard}>
                  <Keyboard
                    ref={r => {
                      this.keyboard = r;
                    }}
                    inputName="keyboard-dialog-input"
                    layoutName={layoutName}
                    onChange={this.onChange}
                    onKeyPress={this.onKeyPress}
                  />
                </div>
              </Dialog>
            )}
          </I18n>
        </div>
      );
    }
  }

  return withStyles(customerStyles)(KeyboardDialog);
}
