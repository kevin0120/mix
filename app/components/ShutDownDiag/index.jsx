import React from 'react';
import PropTypes from 'prop-types';

import SweetAlert from 'react-bootstrap-sweetalert';

import withStyles from '@material-ui/core/styles/withStyles';

import Slide from '@material-ui/core/Slide';
import { I18n } from 'react-i18next';
import sweetAlertStyle from '../../common/jss/views/sweetAlertStyle';

/* eslint-disable react/prefer-stateless-function */
class ShutdownDiag extends React.Component {

  Transition = props => <Slide direction="up" {...props} />;

  render() {
    const { classes, show, content, title, onConfirm, onCancel, showCancel } = this.props;

    return (
      <I18n ns="translations">
        {t =>
          <SweetAlert
            warning
            show={show}
            title={title}
            onConfirm={onConfirm}
            onCancel={onCancel}
            style={{ display: 'block', marginTop: '-100px' }}
            confirmBtnCssClass={`${classes.button} ${classes.success}`}
            cancelBtnCssClass={`${classes.button} ${classes.danger}`}
            confirmBtnText={t('Common.Yes')}
            cancelBtnText={t('Common.No')}
            showCancel={showCancel}
          >
            {content}
          </SweetAlert>
        }
      </I18n>
    );
  }
}

ShutdownDiag.propTypes = {
  show: PropTypes.bool.isRequired,
  onConfirm: PropTypes.func,
  onCancel: PropTypes.func

};
ShutdownDiag.defaultProps = {
  onConfirm: () => {
  },
  onCancel: () => {
  }
};

export default withStyles(sweetAlertStyle)(ShutdownDiag);
