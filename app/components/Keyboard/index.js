/* eslint-disable react/no-this-in-sfc */
import React, { useEffect, useState } from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import Keyboard from 'react-simple-keyboard';
import Slide from '@material-ui/core/Slide';
import InputAdornment from '@material-ui/core/InputAdornment';
import IconButton from '@material-ui/core/IconButton';
import Clear from '@material-ui/icons/CancelOutlined';
import makeStyles from '@material-ui/core/styles/makeStyles';
import { I18n } from 'react-i18next';
import customSelectStyle from '../../common/jss/customSelectStyle';
import Button from '../CustomButtons/Button';
import CustomInput from '../CustomInput/CustomInput';

const lodash = require('lodash');

const customStyles = theme => ({
  ...customSelectStyle(theme),
  bigInput: {
    '&,&::placeholder': {
      fontSize: '50px'
    }
  },
  leftButton: {
    marginLeft: '40%',
    marginRight: '40%'
  },
  rightButton: {
    marginLeft: '40%'
  }
});

export default function withKeyboard(SubComponents) {
  function KeyboardDialog(props) {
    let keyboard = null;
    const [layout, setLayout] = useState('default');
    const [text, setText] = useState('');
    const [config, setConfig] = useState({
      onSubmit: null,
      title: '',
      label: '',
      required: false
    });
    const [show, setShow] = useState(false);

    useEffect(() => {
      if (keyboard) {
        keyboard.keyboard.setInput(text || '');
      }
    }, [keyboard, text]);

    const handleSubmit = () => {
      config.onSubmit(text);
      setShow(false);
    };

    const onKeyPress = press => {
      if (lodash.isEqual(press, '{enter}') && text.length !== 0) {
        handleSubmit();
      }
      if (lodash.isEqual(press, '{shift}')) {
        setLayout(layout === 'default' ? 'shift' : 'default');
      }
    };

    const bindKeyboardInput = c => {
      setConfig(c);
      setText(c.text || '');
      setShow(true);
    };

    const onChangeInput = event => {
      const tx = event.target.value;
      setText(tx);
      keyboard.setInput(tx);
    };

    const { ...restProps } = props;
    const classes = makeStyles(customStyles)();
    const submitEnable = config.required ? text.length !== 0 : true;
    return (
      <I18n ns="translations">
        {t => (
          <React.Fragment>
            <SubComponents
              keyboardInput={c => bindKeyboardInput(c)}
              {...restProps}
            />
            <Dialog
              classes={{
                root: classes.modalRoot,
                paper: `${classes.modal} ${classes.modalLarge}`
              }}
              TransitionComponent={Slide}
              keepMounted
              open={show}
              onClose={() => setShow(false)}
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
                    // onFocus: setActiveInput,
                    required: config.required,
                    value: text || '',
                    onChange: onChangeInput,
                    endAdornment: (
                      <InputAdornment position="end">
                        <IconButton
                          aria-label="clear input"
                          onClick={() => {
                            setText('');
                            keyboard.setInput('');
                          }}
                        >
                          <Clear/>
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
                  // className={classes.modalFooterCenter}
                  style={{ margin: '0 80px' }}
                  onClick={() => setShow(false)}
                  color="danger"
                  // size='lg'
                  autoFocus
                  round
                >
                  {t('Common.Close')}
                </Button>
                <Button
                  style={{ margin: '0 80px' }}
                  onClick={handleSubmit}
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
                    keyboard = r;
                  }}
                  inputName="keyboard-dialog-input"
                  layoutName={layout}
                  onChange={setText}
                  onKeyPress={onKeyPress}
                />
              </div>
            </Dialog>
          </React.Fragment>
        )}
      </I18n>
    );
  }

  return KeyboardDialog;
}
