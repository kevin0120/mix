import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import SweetAlert from 'react-bootstrap-sweetalert';

import withStyles from '@material-ui/core/styles/withStyles';

import Slide from '@material-ui/core/Slide';
import { I18n } from 'react-i18next';
import sweetAlertStyle from '../../common/jss/views/sweetAlertStyle';

// import Dialog from '@material-ui/core/Dialog';
// import DialogActions from '@material-ui/core/DialogActions';
// import DialogContent from '@material-ui/core/DialogContent';
// import DialogContentText from '@material-ui/core/DialogContentText';
// import DialogTitle from '@material-ui/core/DialogTitle';

import * as shutDownDiagActions from '../../actions/shutDownDiag';

const mapStateToProps = (state, ownProps) => ({
  show: state.shutDownDiag.show,
  type: state.shutDownDiag.type,
  msg: state.shutDownDiag.msg,
  ...ownProps
});

const mapDispatchToProps = dispatch =>
  bindActionCreators(shutDownDiagActions, dispatch);

/* eslint-disable react/prefer-stateless-function */
class ConnectedShutdownDiag extends React.Component {
  constructor(props) {
    super(props);
    this.handleClose = this.handleClose.bind(this);
    this.handleConfirm = this.handleConfirm.bind(this);
    this.Transition = this.Transition.bind(this);
  }

  shouldComponentUpdate(nextProps, nextState) {
    const { show, type } = this.props;
    return show !== nextProps.show || type !== nextProps.type;
  }

  handleClose = () => {
    const { type, closeShutDownDiag } = this.props;
    closeShutDownDiag(type); // 關閉窗口
  };

  handleConfirm = () => {
    const { type, confirmShutDownDiag } = this.props;
    confirmShutDownDiag(type); // 關閉窗口
  };

  Transition = props => <Slide direction="up" {...props} />;

  handleReset = () => {};

  render() {
    const { classes, show, type, msg } = this.props;

    let isShow = false;

    if (
      (type === 'bypass' || type === 'shutdown' || type === 'verify') &&
      show
    ) {
      // bypass 车辆放行, shutdown关机, verify验证车辆序号是否准确
      isShow = true;
    }

    const typeByPass = type === 'bypass';

    let title = 'Common.Shutdown';
    let question = 'Common.QuestShutdown';

    switch (type) {
      case 'bypass':
        title = 'Common.Bypass';
        question = 'Common.QuestBypass';
        break;
      case 'shutdown':
        title = 'Common.Shutdown';
        question = 'Common.QuestShutdown';
        break;
      case 'verify':
        title = 'Common.VerifyCar';
        question = 'Common.QuestConfirmCarInfo';
        setTimeout(this.handleClose, 5000);
        break;
      default:
        break;
    }

    return (
      <I18n ns="translations">
        {t =>
          typeByPass ? (
            <SweetAlert
              warning
              show={isShow}
              style={{ display: 'block', marginTop: '-100px' }}
              title={t(title)}
              onConfirm={() => this.handleConfirm()}
              onCancel={() => this.handleClose()}
              confirmBtnCssClass={`${classes.button} ${classes.success}`}
              cancelBtnCssClass={`${classes.button} ${classes.danger}`}
              confirmBtnText={t('Common.Yes')}
              cancelBtnText={t('Common.No')}
              showCancel
            >
              <div>
                <h3>{t(question)} </h3>
                <h4 style={{ textAlign: 'center' }}> OR </h4>
                {/* <Button onClick={() => this.handleReset} color="danger" round> */}
                {/* {t('Common.Reset')} */}
                {/* </Button> */}
              </div>
            </SweetAlert>
          ) : (
            <SweetAlert
              warning
              show={isShow}
              style={{ display: 'block', marginTop: '-100px' }}
              title={t(title)}
              onConfirm={() => this.handleConfirm(type)}
              onCancel={() => this.handleClose()}
              confirmBtnCssClass={`${classes.button} ${classes.success}`}
              cancelBtnCssClass={`${classes.button} ${classes.danger}`}
              confirmBtnText={t('Common.Yes')}
              cancelBtnText={t('Common.No')}
              showCancel
            >
              {t(question)} {msg}
            </SweetAlert>
          )
        }
      </I18n>
    );
  }
}

ConnectedShutdownDiag.propTypes = {
  show: PropTypes.bool.isRequired,
  msg: PropTypes.string.isRequired
};

const ShutdownDiag = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedShutdownDiag);

export default withStyles(sweetAlertStyle)(ShutdownDiag);
