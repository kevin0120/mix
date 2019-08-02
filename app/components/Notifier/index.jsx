/* eslint-disable react/prop-types */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { withSnackbar } from 'notistack';
import { CloseSharp } from '@material-ui/icons';
import actions from '../../modules/Notifier/action';
import Button from '../CustomButtons/Button';

class Notifier extends Component {
  // eslint-disable-next-line react/sort-comp
  displayed = [];

  storeDisplayed = (id) => {
    this.displayed = [...this.displayed, id];
  };

  shouldComponentUpdate({ notifications: newSnacks = [] }) {
    if (!newSnacks.length) {
      this.displayed = [];
      return false;
    }

    const { notifications: currentSnacks } = this.props;
    let notExists = false;
    for (let i = 0; i < newSnacks.length; i += 1) {
      const newSnack = newSnacks[i];
      if (newSnack.dismissed) {
        const {removeSnackbarAction, closeSnackbar} = this.props;
        removeSnackbarAction(newSnack.key);
        closeSnackbar(newSnack.key);
      }

      if (!notExists) {
        // update notExists Flag
        notExists = notExists || !currentSnacks.filter(({ key }) => newSnack.key === key).length;
      }
    }
    return notExists;
  }

  componentDidUpdate() {
    const {
      notifications = [],
      enqueueSnackbar,
      removeSnackbarAction,
      closeSnackbarAction,
    } = this.props;

    notifications.forEach(({ key, message, options = {} }) => {
      // Do nothing if snackbar is already displayed
      if (this.displayed.includes(key)) return;
      // Display snackbar using notistack
      enqueueSnackbar(message, {
        ...options,
        action: (k) => (
          <Button
            justIcon
            regular
            onClick={() => {
              closeSnackbarAction(k);
            }}
            simple
            round
          >
            <CloseSharp color="inherit"/>
          </Button>
        ),
        onClose: (event, reason, k) => {
          if (options.onClose) {
            options.onClose(event, reason, k);
          }

          // Dispatch action to remove snackbar from redux store
          removeSnackbarAction(k);
        }
      });
      // Keep track of snackbars that we've displayed
      this.storeDisplayed(key);
    });
  }

  render() {
    return null;
  }
}

const mapStateToProps = (store, props) => ({
  ...props,
  notifications: store.Notifier.notifications
});

const mapDispatchToProps = {
  removeSnackbarAction: actions.removeSnackbar,
  closeSnackbarAction: actions.closeSnackbar
};

export default withSnackbar(connect(
  mapStateToProps,
  mapDispatchToProps
)(Notifier));
