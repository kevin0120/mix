import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import withStyles from '@material-ui/core/styles/withStyles';

import Slide from '@material-ui/core/Slide';

import { I18n } from 'react-i18next';

import Assignment from '@material-ui/icons/Assignment';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import CardActionArea from '@material-ui/core/CardActionArea';

import DialogActions from '@material-ui/core/DialogActions';
import Dialog from '@material-ui/core/Dialog';
import { bindActionCreators } from 'redux';
import GridContainer from '../Grid/GridContainer';
import GridItem from '../Grid/GridItem';
import Card from '../Card/Card';
import CardHeader from '../Card/CardHeader';
import CardIcon from '../Card/CardIcon';
import CardBody from '../Card/CardBody';
import Table from '../Table/Table';
import Button from '../CustomButtons/Button';

import { setResultDiagShow } from '../../actions/resultDiag';
import { NewCar } from '../../actions/scannerDevice';

import resultDiagStyles from './styles';
import configs from '../../shared/config';

import { OPERATION_STATUS, OPERATION_SOURCE } from '../../reducers/operations';

const lodash = require('lodash');

const mapStateToProps = (state, ownProps) => ({
  show: state.resultDiag.show,
  results: state.operations.results,
  jobID: state.operations.jobID,
  nextWorkorder: state.ongoingOperation,
  ...ownProps
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ setResultDiagShow, NewCar }, dispatch);

/* eslint-disable react/prefer-stateless-function */
class ConnectedResultDialog extends React.Component {
  constructor(props) {
    super(props);
    this.Transition = this.Transition.bind(this);
  }

  shouldComponentUpdate(nextProps, nextState) {
    const { show } = this.props;
    return nextProps.show !== show;
  }

  Transition = props => <Slide direction="up" {...props} />;

  handleClose = e => {
    e.preventDefault();
    const { setResultDiagShow } = this.props;
    setResultDiagShow(false);
  };

  handleClickNewWorkOrder = e => {
    e.preventDefault();
    const { NewCar } = this.props;
    const { vin } = this.state.nextWorkorder;
    if (!lodash.isNil(vin) || vin !== '') {
      NewCar(vin, OPERATION_SOURCE.MANUAL);
    }
  };

  render() {
    const { classes, show, results, nextWorkorder, jobID } = this.props;

    const showNextVehicle = configs.operationSettings.opMode === 'order';

    const nw = [];

    if (nextWorkorder) {
      nw.push([
        nextWorkorder.vin,
        nextWorkorder.model,
        nextWorkorder.lnr,
        nextWorkorder.knr,
        nextWorkorder.long_pin
      ]);
    }

    const localResults = [];
    for (let i = 0; i < results.length; i++) {
      if (showNextVehicle) {
        localResults.push([
          results[i].pset,
          results[i].mi,
          results[i].wi,
          results[i].ti,
          `${results[i].group_sequence}/${results[results.length - 1].group_sequence}`,
          results[i].result
        ]);
      } else {
        localResults.push([
          jobID,
          results[i].mi,
          results[i].wi,
          results[i].ti,
          `${results[i].group_sequence}/${results[results.length - 1].group_sequence}`,
          results[i].result
        ]);
      }
    }

    return (
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
            <DialogTitle id="form-dialog-title" className={classes.modalHeader}>
              {t('Common.Result')}
            </DialogTitle>
            <DialogContent className={classes.modalBody}>
              <div>
                <GridContainer className={classes.root}>
                  <GridItem xs={12}>
                    <Card>
                      <CardHeader color="info" icon>
                        <CardIcon color="info">
                          <Assignment />
                        </CardIcon>
                        <h4 style={{ color: '#000' }}>
                          {t('main.currentOrder')}
                        </h4>
                      </CardHeader>
                      <CardBody>
                        <Table
                          tableHeaderColor="info"
                          tableHead={['程序号', '扭矩(N·M)', '角度(Deg)', '用时(ms)', '批次', '结果']}
                          tableData={localResults}
                          colorsColls={['info']}
                        />
                      </CardBody>
                    </Card>
                  </GridItem>
                </GridContainer>
              </div>
              {showNextVehicle ? (
                <div>
                  <GridContainer className={classes.root}>
                    <GridItem xs={12}>
                      <Card>
                        <CardHeader color="info" icon>
                          <CardIcon color="info">
                            <Assignment />
                          </CardIcon>
                          <h4 style={{ color: '#000' }}>
                            {t('main.nextOrder')}
                          </h4>
                        </CardHeader>
                        <CardActionArea
                          component={Button}
                          onClick={e => this.handleClickNewWorkOrder(e)}
                          className={classes.cardActionArea}
                        >
                          <CardBody>
                            <Table
                              tableHeaderColor="info"
                              tableHead={[
                                'Vin',
                                '车型',
                                '车序',
                                'Knr',
                                'LongPin'
                              ]}
                              tableData={nw}
                              colorsColls={['info']}
                            />
                          </CardBody>
                        </CardActionArea>
                      </Card>
                    </GridItem>
                  </GridContainer>
                </div>
              ) : null}
            </DialogContent>
            <DialogActions
              className={`${classes.modalFooter} ${classes.modalFooterCenter}`}
            >
              <Button
                onClick={e => this.handleClose(e)}
                color="info"
                autoFocus
                round
              >
                {t('Common.Close')}
              </Button>
            </DialogActions>
          </Dialog>
        )}
      </I18n>
    );
  }
}

ConnectedResultDialog.propTypes = {
  show: PropTypes.bool.isRequired,
  nextWorkorder: PropTypes.shape({}).isRequired,
  jobID: PropTypes.number.isRequired,
  results: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.number
    }).isRequired
  ).isRequired,
  setResultDiagShow: PropTypes.func.isRequired,
  NewCar: PropTypes.func.isRequired
};

const ResultDialog = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedResultDialog);

export default withStyles(resultDiagStyles)(ResultDialog);
