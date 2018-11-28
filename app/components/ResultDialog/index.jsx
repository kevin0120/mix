import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
// import { fetchNextWorkOrder } from '../../actions/ongoingWorkOrder';

import withStyles from '@material-ui/core/styles/withStyles';

// import { IOSet } from "../../actions/controllerIO";

import Slide from '@material-ui/core/Slide';

import { I18n } from 'react-i18next';

// import { setShutdown, switchReady } from '../../actions/commonActions';
import customSelectStyle from '../../common/jss/customSelectStyle';
import GridContainer from '../Grid/GridContainer';
import GridItem from '../Grid/GridItem';
import Card from '../Card/Card';
import CardHeader from '../Card/CardHeader';
import CardIcon from '../Card/CardIcon';
import Assignment from '@material-ui/icons/Assignment';
import CardBody from '../Card/CardBody';
import Table from '../Table/Table';
import ReactTable from 'react-table';
import { cardTitle } from '../../common/jss/material-react-pro';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/Select';
import MenuItem from '@material-ui/core/MenuItem';
import DialogActions from '@material-ui/core/DialogActions';
import Button from '../CustomButtons/Button';
import Dialog from '@material-ui/core/Dialog';

const mapStateToProps = (state, ownProps) => ({
  show: state.resultDiag.show,
  results: state.operations.results,
  nextWorkorder: state.ongoingOperation,
  ...ownProps
});

const mapDispatchToProps = {
  // fetchNextWorkOrder
};

/* eslint-disable react/prefer-stateless-function */
class ConnectedResultDialog extends React.Component {
  constructor(props) {
    super(props);
    this.Transition = this.Transition.bind(this);
  }

  shouldComponentUpdate(nextProps,nextState) {
    const {show} = this.props;
    return nextProps.show !== show;
  }

  Transition = props => {
    return <Slide direction="up" {...props} />;
  };

  handleClose = (e) => {

  };

  render() {
    const {
      classes,
      show,
      results,
      nextWorkorder,
    } = this.props;


    const nw = [];

    if (nextWorkorder) {
      nw.push([nextWorkorder.vin, nextWorkorder.model, nextWorkorder.lnr, nextWorkorder.knr, nextWorkorder.long_pin]);
    }

    const localResults = [];
    for (let i = 0; i < results.length; i++) {
      localResults.push([
        results[i].pset,
        results[i].mi,
        results[i].wi,
        results[i].ti,
        results[i].result
      ]);
    }

    return (
      <I18n ns="translations">
        {t => (
          <Dialog
            classes={{
              root: classes.modalRoot,
              paper: classes.modal + ' ' + classes.modalLarge
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
                      <CardHeader color="primary" icon>
                        <CardIcon color="primary">
                          <Assignment />
                        </CardIcon>
                        <h4 className={classes.cardIconTitle}>
                          {t('main.resultQuery')}
                        </h4>
                      </CardHeader>
                      <CardBody>
                        <Table
                          tableHeaderColor="primary"
                          tableHead={['程序号', '扭矩', '角度', '用时', '结果']}
                          tableData={localResults}
                          colorsColls={['primary']}
                        />
                      </CardBody>
                    </Card>
                  </GridItem>
                </GridContainer>
              </div>
              <div>
                <GridContainer className={classes.root}>
                  <GridItem xs={12}>
                    <Card>
                      <CardHeader color="primary" icon>
                        <CardIcon color="primary">
                          <Assignment />
                        </CardIcon>
                        <h4 className={classes.cardIconTitle}>
                          {t('main.nextOrder')}
                        </h4>
                      </CardHeader>
                      <CardBody>
                        <Table
                          tableHeaderColor="primary"
                          tableHead={['Vin', '车型', '车序', 'Knr', 'LongPin']}
                          tableData={nw}
                          colorsColls={['primary']}
                        />
                      </CardBody>
                    </Card>
                  </GridItem>
                </GridContainer>
              </div>
            </DialogContent>
            <DialogActions
              className={classes.modalFooter + ' ' + classes.modalFooterCenter}
            >
              <Button
                onClick={(e) => this.handleClose(e)}
                color="primary"
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
  results: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.number,
  }).isRequired).isRequired,
  // fetchNextWorkOrder: PropTypes.func.isRequired,
};

const ResultDialog = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedResultDialog);

export default withStyles(customSelectStyle)(ResultDialog);
