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
    newSnacks.forEach(s => {
      if (s.dismissed) {
        const { removeSnackbarAction, closeSnackbar } = this.props;
        console.log('removing');
        removeSnackbarAction(s.key);
        closeSnackbar(s.key);
        return;
      }
      if (!notExists) {
        // update notExists Flag
        notExists = notExists || !currentSnacks.filter(({ key }) => s.key === key).length;
      }
    });
    return notExists;
  }

  componentDidUpdate() {
    const {
      notifications = [],
      enqueueSnackbar,
      closeSnackbar,
      removeSnackbarAction,
      closeSnackbarAction
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
              closeSnackbar(k);
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
          closeSnackbar(k);
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
