import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import CheckCircleIcon from '@material-ui/icons/CheckCircleOutlineOutlined';
import ErrorIcon from '@material-ui/icons/ErrorOutlineOutlined';
import InfoIcon from '@material-ui/icons/InfoOutlined';
import WarningIcon from '@material-ui/icons/WarningOutlined';
import { withStyles } from '@material-ui/core/styles';
import Snackbar from '../Snackbar/Snackbar';
import * as NotificationActions from '../../modules/notification/notification';

import notificationsStyle from '../../common/jss/views/notificationsStyle';

const mapStateToProps = (state, ownProps) => ({
  message: state.notify.message,
  variant: state.notify.variant,
  isShow: state.notify.isShow,
  ...ownProps
});

function mapDispatchToProps(dispatch) {
  return bindActionCreators(NotificationActions, dispatch);
}

const variantIcon = {
  success: CheckCircleIcon,
  warning: WarningIcon,
  maintenance: InfoIcon,
  error: ErrorIcon,
  info: InfoIcon
};

/* eslint-disable react/prefer-stateless-function */
class ConnectedNoty extends React.PureComponent {
  constructor(props) {
    super(props);
    this.handleClose = this.handleClose.bind(this);
  }

  handleClose() {
    const { closeNotification } = this.props;
    closeNotification();
  }

  render() {
    const { message, variant, isShow } = this.props;

    const Icon = variantIcon[variant];

    return (
      <Snackbar
        place="tl"
        color={variant}
        icon={Icon}
        message={message}
        open={isShow}
        closeNotification={this.handleClose}
        close
      />
    );
  }
}

ConnectedNoty.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  message: PropTypes.string.isRequired,
  isShow: PropTypes.bool.isRequired,
  variant: PropTypes.oneOf([
    'success',
    'warning',
    'error',
    'info',
    'danger',
    'maintenance'
  ]).isRequired,
  closeNotification: PropTypes.func.isRequired
};

const Notify = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedNoty);

export default withStyles(notificationsStyle)(Notify);
