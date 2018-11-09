import React from "react";
import { connect } from 'react-redux';
import PropTypes from "prop-types";

// @material-ui/core components
import withStyles from "@material-ui/core/styles/withStyles";

import { I18n } from 'react-i18next';

import Keyboard from 'react-simple-keyboard'


// core components
import Button from "../../components/CustomButtons/Button";
import CustomInput from "../../components/CustomInput/CustomInput";
import Card from "../../components/Card/Card";
import CardBody from "../../components/Card/CardBody";
import CardAvatar from "../../components/Card/CardAvatar";
import CardFooter from "../../components/Card/CardFooter";

import lockScreenPageStyle from "../../common/jss/views/LockScreenStyle";

const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  name: state.userInfo.name,
  avatarImg: state.userInfo.image_small,
  ...ownProps,
});

const mapDispatchToProps = {
};

class ConnectedLockScreenPage extends React.Component {
  constructor(props) {
    super(props);
    // we use this to make the card to appear after the page has been rendered

    this.keyboard = null;

    this.state = {
      cardAnimaton: "cardHidden",
      layoutName: "default",
      inputName: "",
      value: ""
    };
  }

  componentDidMount() {
    // we add a hidden class to the card and after 700 ms we delete it and the transition appears
    this.timeOutFunction = setTimeout(
      () => {
        this.setState({ cardAnimaton: "" });
      },
      700
    );
  }

  componentWillUnmount(){
    clearTimeout(this.timeOutFunction);
    this.timeOutFunction = null;
  }

  setActiveInput = (event) => {
    this.setState({
      inputName: event.target.id
    });
  };

  removeActiveInput = () => {
    // const { inputName, value} = this.state;
    // this.setState({
    //   inputName: "",
    //   init: false,
    // });
  };

  handleInput = (input) => {
    this.setState(
      {
        value: input
      }
    )
  };

  handleUnlock = () => {
    this.setState(
      {
        value: "",
        inputName: ""
      }
    )
  };

  handlePress = (press) => {
    const { value} = this.state;
    if (lodash.isEqual(press, '{enter}') && value.length !== 0){
      // 发送unlock命令
      this.handleUnlock();
    }
    if (lodash.isEqual(press, '{shift}')) {
      const { layoutName }= this.state;

      this.setState({
        layoutName: layoutName === "default" ? "shift" : "default"
      });
    }
  };

  render() {
    const { classes, name, avatarImg } = this.props;

    const { inputName, value, cardAnimaton, layoutName} = this.state;

    const disabled = !value.length;


    return (
      <I18n ns="translations">
        {
          t => (
            <div className={classes.container}>
              <form>
                  <Card
                    profile
                    className={
                      `${classes.customCardClass  } ${  classes[cardAnimaton]}`
                    }
                  >
                    <CardAvatar profile className={classes.cardAvatar}>
                      <a href="#pablo" onClick={e => e.preventDefault()}>
                        <img src={avatarImg} alt="..."/>
                      </a>
                    </CardAvatar>
                    <CardBody profile>
                      <h4 className={classes.cardTitle}>{name}</h4>
                      <CustomInput
                        labelText={t('Lock.Pwd')}
                        id="company-disabled"
                        formControlProps={{
                          fullWidth: true
                        }}
                        inputProps={{
                          type: "password",
                          onFocus: this.setActiveInput,
                          onBlur: this.removeActiveInput,
                          value
                        }}
                      />
                    </CardBody>
                    <CardFooter className={classes.justifyContentCenter}>
                      <Button disabled={disabled} color="rose" round onClick={this.handleUnlock}>
                        {t('Lock.Unlock')}
                      </Button>
                    </CardFooter>
                  </Card>
              </form>
              {
                inputName.length !== 0 ? <div className={classes.keyboard} >
                  <Keyboard
                    ref={r => this.keyboard = r}
                    inputName={inputName}
                    layoutName={layoutName}
                    onChange={this.handleInput}
                    onKeyPress={this.handlePress}
                  />
                </div>: null
              }

            </div>
          )
        }
      </I18n>
    );
  }
}

ConnectedLockScreenPage.propTypes = {
  classes: PropTypes.object.isRequired,
  name: PropTypes.string.isRequired,
  avatarImg: PropTypes.string.isRequired,
};

const LockScreenPage = connect(mapStateToProps, mapDispatchToProps)(ConnectedLockScreenPage);


export default withStyles(lockScreenPageStyle)(LockScreenPage);



// WEBPACK FOOTER //
// ./src/views/Pages/LockScreenPage.jsx
