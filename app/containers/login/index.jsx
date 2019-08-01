import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';
import InputAdornment from '@material-ui/core/InputAdornment';

// @material-ui/icons
import Face from '@material-ui/icons/Face';
import LockOutlined from '@material-ui/icons/LockOutlined';

import { I18n } from 'react-i18next';

// core components
import Fingerprint from '@material-ui/icons/Fingerprint';
import Keyboard from 'react-simple-keyboard';
import GridContainer from '../../components/Grid/GridContainer';
import GridItem from '../../components/Grid/GridItem';
import CustomInput from '../../components/CustomInput/CustomInput';
import Button from '../../components/CustomButtons/Button';
import Card from '../../components/Card/Card';
import CardBody from '../../components/Card/CardBody';
import CardHeader from '../../components/Card/CardHeader';
import CardFooter from '../../components/Card/CardFooter';

import loginPageStyle from '../../common/jss/views/loginPageStyle';
import { loginRequest } from '../../modules/user/action';

const lodash = require('lodash');

class LoginPage extends React.Component {
  constructor(props) {
    super(props);
    // we use this to make the card to appear after the page has been rendered
    this.keyboard = null;

    this.state = {
      cardAnimaton: 'cardHidden',
      layoutName: 'default',
      inputName: '',
      allInput: {
        firstname: '',
        password: ''
      }
    };
  }

  componentDidMount() {
    // we add a hidden class to the card and after 700 ms we delete it and the transition appears
    this.timeOutFunction = setTimeout(() => {
      this.setState({ cardAnimaton: '' });
    }, 700);
  }

  componentWillUnmount() {
    clearTimeout(this.timeOutFunction);
    this.timeOutFunction = null;
  }

  setActiveInput = event => {
    this.setState({
      inputName: event.target.id
    });
  };

  removeActiveInput = () => {
    // this.setState({
    //   inputName: "",
    // });
  };

  // onChangeAll = (input) => {
  //   this.setState({
  //     allInput: input
  //   }, () => {
  //     console.log("Inputs changed", input);
  //   });
  // };

  handleInput = input => {
    const { inputName, allInput } = this.state;
    if (inputName === 'firstname') {
      this.setState({
        allInput: {
          ...allInput,
          firstname: input
        }
      });
    }
    if (inputName === 'password') {
      this.setState({
        allInput: {
          ...allInput,
          password: input
        }
      });
    }
  };

  handleSubmit = () => {
    const { firstname, password } = this.state.allInput;
    const { login } = this.props;
    login(firstname, password);

    this.setState({
      inputName: '',
      allInput: {
        firstname: '',
        password: ''
      }
    });

  };

  handlePress = press => {
    const { allInput } = this.state;
    if (
      lodash.isEqual(press, '{enter}') &&
      allInput.firstname.length !== 0 &&
      allInput.password.length !== 0
    ) {
      this.handleSubmit();
    }
    if (lodash.isEqual(press, '{shift}')) {
      const { layoutName } = this.state;

      this.setState({
        layoutName: layoutName === 'default' ? 'shift' : 'default'
      });
    }
  };

  render() {
    const { classes } = this.props;

    const { inputName, allInput, layoutName, cardAnimaton } = this.state;

    const disabled = !allInput.firstname.length || !allInput.password.length;
    return (
      <I18n ns="translations">
        {t => (
          <div className={classes.container}>
            <GridContainer justify="center">
              <GridItem xs={12} sm={6} md={4}>
                <form>
                  <Card login className={classes[cardAnimaton]}>
                    <CardHeader
                      className={`${classes.cardHeader} ${classes.textCenter}`}
                      color="info"
                    >
                      <h4 className={classes.cardTitle}>Log in</h4>
                    </CardHeader>
                    <CardBody>
                      <CustomInput
                        labelText={t('Login.Name')}
                        id="firstname"
                        formControlProps={{
                          fullWidth: true,
                          required: true
                        }}
                        inputProps={{
                          endAdornment: (
                            <InputAdornment position="end">
                              <Face className={classes.inputAdornmentIcon}/>
                            </InputAdornment>
                          ),
                          onFocus: this.setActiveInput,
                          onBlur: this.removeActiveInput,
                          required: true,
                          value: allInput.firstname || ''
                        }}
                      />
                      <CustomInput
                        labelText={t('Login.Pwd')}
                        id="password"
                        formControlProps={{
                          fullWidth: true,
                          required: true
                        }}
                        inputProps={{
                          endAdornment: (
                            <InputAdornment position="end">
                              <LockOutlined
                                className={classes.inputAdornmentIcon}
                              />
                            </InputAdornment>
                          ),
                          type: 'password',
                          onFocus: this.setActiveInput,
                          onBlur: this.removeActiveInput,
                          required: true,
                          value: allInput.password || ''
                        }}
                      />
                    </CardBody>
                    <CardFooter className={classes.justifyContentCenter}>
                      <Button
                        disabled={disabled}
                        color="rose"
                        size="lg"
                        round
                        onClick={this.handleSubmit}
                      >
                        <Fingerprint/>
                        {t('Login.Submit')}
                      </Button>
                    </CardFooter>
                  </Card>
                </form>
              </GridItem>
            </GridContainer>
            {inputName.length !== 0 ? (
              <div className={classes.keyboard}>
                <Keyboard
                  ref={r => (this.keyboard = r)}
                  inputName={inputName}
                  layoutName={layoutName}
                  onChange={this.handleInput}
                  // onChangeAll={inputs => this.onChangeAll(inputs)}
                  onKeyPress={this.handlePress}
                />
              </div>
            ) : null}
          </div>
        )}
      </I18n>
    );
  }
}

LoginPage.propTypes = {
  classes: PropTypes.shape({}).isRequired
};

const mapState = (state, props) => ({});

const mapDispatch = {
  login: loginRequest
};

export default withStyles(loginPageStyle)(connect(mapState, mapDispatch)(LoginPage));

// WEBPACK FOOTER //
// ./src/views/Pages/LoginPage.jsx
