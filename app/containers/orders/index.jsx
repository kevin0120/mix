/*
 * Copyright (c) 2018. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

import React from 'react';
import { connect } from 'react-redux';
import ReactTable from 'react-table';
import PropTypes from 'prop-types';
import { push } from 'connected-react-router';
import Divider from '@material-ui/core/Divider';

import { I18n } from 'react-i18next';

import SweetAlert from 'react-bootstrap-sweetalert';
import withStyles from '@material-ui/core/styles/withStyles';
import Assignment from '@material-ui/icons/Assignment';
import Dvr from '@material-ui/icons/Dvr';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import isURL from 'validator/lib/isURL';
import Input from '@material-ui/core/Input';
import sweetAlertStyle from '../../common/jss/views/sweetAlertStyle';

// @material-ui/core components
// @material-ui/icons
// core components
import GridContainer from '../../components/Grid/GridContainer';
import GridItem from '../../components/Grid/GridItem';
import Button from '../../components/CustomButtons/Button';
import Card from '../../components/Card/Card';
import CardBody from '../../components/Card/CardBody';
import CardIcon from '../../components/Card/CardIcon';
import CardHeader from '../../components/Card/CardHeader';

import { defaultClient } from '../../common/utils';
import withKeyboard from '../../components/Keyboard';

const lodash = require('lodash');

const styles = theme=>({
  ...sweetAlertStyle(theme),
  root: {
    flexGrow: 1,
    zIndex: 1,
    height: '100%',
    overflow: 'hidden',
    position: 'relative',
    display: 'flex',
    width: '100%'
  },
  cardIconTitle: {
    ...theme.title.card,
    marginTop: '15px',
    marginBottom: '0px'
  },
  InputRoot: {
    width: '100%',
    height: '36px',
    overflow: 'hidden'
  },
  InputInput: {
    width: '100%',
    height: '100%'
  },
  alertPosition: {
    top: '25%',
    display: 'block',
    marginTop: '-100px'
  }
});

const mapStateToProps = (state, ownProps) => ({
  masterpcUrl: state.connections.masterpc,
  workcenterCode: state.connections.workcenterCode,
  ...ownProps
});

const mapDispatchToProps = {
  // dispatchNewCar: ScannerNewData,
  dispatchPush: push
};

function requestData(masterpcUrl, workcenterCode) {
  const url = `${masterpcUrl}/rush/v1/workorders`;
  if (!isURL(url, { require_protocol: true })) {
    return new Promise(() => {
      throw new Error('conn is Error!');
    });
  }
  return defaultClient.get(url, {
    params: {
      workcenter_code: workcenterCode,
      status: 'ready',
      limit: 100
    }
  });
}

class WorkOrder extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      selectObj: null,
      confirmInfo: null,
      alert: 'none' // '车辆详情',''
    };
    this.fetchData = this.fetchData.bind(this);
  }

  componentDidMount() {
    this.fetchData();
  }

  fetchData() {
    const { masterpcUrl, workcenterCode } = this.props;
    requestData(masterpcUrl, workcenterCode)
      .then(res => {
        const statusCode = res.status;
        if (statusCode === 200) {
          this.setState({
            data: res.data.map((item, key) => ({
              id: key,
              // timestamp: dayjs(item.timestamp).format('YYYY MM-DD HH:mm:ss'),
              vin: item.vin,
              model: item.model,
              long_pin: item.long_pin,
              knr: item.knr,
              lnr: item.lnr,
              actions: (
                // we've added some custom button actions
                <div className="actions-right">
                  {/* use this button to add a like kind of action */}
                  <Button
                    justIcon
                    round
                    simple
                    onClick={e => {
                      e.preventDefault();
                      e.stopPropagation();
                      const obj = this.state.data.find(o => o.id === key);
                      this.setState({
                        alert: '车辆详情',
                        selectObj: obj
                      });
                    }}
                    style={{
                      width: 100,
                      height: '100%'
                    }}
                    color="warning"
                    className="edit"
                  >
                    <Dvr />
                  </Button>{' '}
                </div>
              )
            }))
          });
        }
      })
      .catch(error => {
        console.log(`get error: ${error.toString()}`);
      });
  }

  handleClose = () => {
    this.setState({
      alert: 'none'
    });
  };

  handleRowClick = rowInfo => {
    this.setState({
      alert: '确认选择',
      confirmInfo: rowInfo.original
    });
  };

  handleRowConfirm = () => {
    const { dispatchNewCar, dispatchPush } = this.props;
    const { confirmInfo } = this.state;
    dispatchPush('/working');
    dispatchNewCar(confirmInfo.long_pin);
  };

  Msg = selectObj =>
    selectObj ? (
      <div>
        <List>
          <ListItem>
            <ListItemText primary={`车序: ${selectObj.lnr}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`VIN:   ${selectObj.vin}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`车型:   ${selectObj.model}`} />
          </ListItem>
          <li>
            <Divider inset />
          </li>
          <ListItem>
            <ListItemText primary={`KNR: ${selectObj.knr}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`LongPIN:   ${selectObj.long_pin}`} />
          </ListItem>
        </List>
      </div>
    ) : (
      ' '
    );

  ConfirmInfo = info =>
    info ? (
      <div>
        <List>
          <ListItem>
            <ListItemText primary={`VIN:   ${info.vin}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`车型:   ${info.model}`} />
          </ListItem>
          <li>
            <Divider inset />
          </li>
          <ListItem>
            <ListItemText primary={`KNR: ${info.knr}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`车序: ${info.lnr}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`LongPIN:   ${info.long_pin}`} />
          </ListItem>
        </List>
      </div>
    ) : (
      ' '
    );

  Alert = (type, t) => {
    if (type === 'none') {
      return null;
    }
    const { classes } = this.props;
    const { selectObj, confirmInfo } = this.state;
    let AlertSettings = {};
    let content = null;
    switch (type) {
      case '车辆详情':
        AlertSettings = {
          info: true,
          show: true,
          title: '车辆详情',
          onConfirm: this.handleClose,
          onCancel: this.handleClose,
          confirmBtnCssClass: `${classes.button} ${classes.success}`,
          cancelBtnCssClass: `${classes.button} ${classes.danger}`,
          confirmBtnText: t('Common.Yes'),
          cancelBtnText: t('Common.No'),
          showCancel: false
        };
        content = this.Msg(selectObj);
        break;
      case '确认选择':
        AlertSettings = {
          info: true,
          show: true,
          title: '确认选择',
          onConfirm: this.handleRowConfirm,
          onCancel: this.handleClose,
          confirmBtnCssClass: `${classes.button} ${classes.success}`,
          cancelBtnCssClass: `${classes.button} ${classes.danger}`,
          confirmBtnText: t('Common.Yes'),
          cancelBtnText: t('Common.No'),
          showCancel: true
        };
        content = this.ConfirmInfo(confirmInfo);
        break;
      default:
        return null;
    }
    return (
      <SweetAlert {...AlertSettings} style={styles.alertPosition}>
        {content}
      </SweetAlert>
    );
  };

  render() {
    const { classes } = this.props;
    const { data, alert } = this.state;

    return (
      <I18n ns="translations">
        {t => (
          <div>
            <GridContainer className={classes.root}>
              <GridItem xs={12}>
                <Card>
                  <CardHeader color="info" icon>
                    <CardIcon color="info">
                      <Assignment />
                    </CardIcon>
                    <h4 className={classes.cardIconTitle}>
                      {t('main.orders')}
                    </h4>
                  </CardHeader>
                  <CardBody>
                    <ReactTable
                      getTrProps={(state, rowInfo, column, instance) => {
                        return {
                          onClick: (e, handleOriginal) => {
                            e.preventDefault();
                            e.stopPropagation();

                            this.handleRowClick(rowInfo);
                            // IMPORTANT! React-Table uses onClick internally to trigger
                            // events like expanding SubComponents and pivots.
                            // By default a custom 'onClick' handler will override this functionality.
                            // If you want to fire the original onClick handler, call the
                            // 'handleOriginal' function.
                            if (handleOriginal) {
                              handleOriginal();
                            }
                          }
                        };
                      }}
                      data={data}
                      filterable
                      columns={[
                        {
                          Header: '车序',
                          accessor: 'lnr',
                          sortable: true,
                          filterable: false
                        },
                        {
                          Header: 'VIN',
                          accessor: 'vin',
                          filterMethod: (filter, row) => {
                            return lodash.includes(
                              lodash.toUpper(row[filter.id]),
                              lodash.toUpper(this.state.vinFilter || '')
                            );
                          },
                          Filter: ({ filter, onChange }) => (
                            <Input
                              onClick={() => {
                                this.props.keyboardInput({
                                  onSubmit: text => {
                                    this.setState({ vinFilter: text }, () => {
                                      onChange(this.state.vinFilter);
                                    });
                                  },
                                  text: this.state.vinFilter,
                                  title: 'VIN',
                                  label: 'VIN'
                                });
                              }}
                              classes={{
                                root: classes.InputRoot,
                                input: classes.InputInput
                              }}
                              // style={{ width: "100%" ,height:'36px'}}
                              value={this.state.vinFilter || ''}
                            />
                          )
                        },
                        {
                          Header: '车型',
                          accessor: 'model',
                          filterMethod: (filter, row) => {
                            return lodash.includes(
                              lodash.toUpper(row[filter.id]),
                              lodash.toUpper(this.state.vehicleTypeFilter || '')
                            );
                          },
                          Filter: ({ filter, onChange }) => (
                            <Input
                              onClick={() => {
                                this.props.keyboardInput({
                                  onSubmit: text => {
                                    this.setState(
                                      { vehicleTypeFilter: text },
                                      () => {
                                        onChange(this.state.vehicleTypeFilter);
                                      }
                                    );
                                  },
                                  text: this.state.vehicleTypeFilter,
                                  title: '车型',
                                  label: '车型'
                                });
                              }}
                              classes={{
                                root: classes.InputRoot,
                                input: classes.InputInput
                              }}
                              // style={{ width: "100%" ,height:'36px'}}
                              value={this.state.vehicleTypeFilter || ''}
                            />
                          )
                        },
                        {
                          Header: 'KNR',
                          accessor: 'knr',
                          sortable: false,
                          filterable: false
                        },
                        {
                          Header: 'LongPIN',
                          accessor: 'long_pin',
                          sortable: false,
                          filterable: false
                        },
                        // {
                        //   Header: '拧紧时间',
                        //   accessor: 'timestamp',
                        //   filterable: false,
                        //   filterMethod: (filter, row) =>
                        //     lodash.includes(
                        //       lodash.toUpper(row[filter.id]),
                        //       lodash.toUpper(filter.value)
                        //     )
                        // },
                        {
                          Header: 'Actions',
                          accessor: 'actions',
                          sortable: false,
                          filterable: false
                        }
                      ]}
                      defaultPageSize={10}
                      showPaginationTop
                      showPaginationBottom={false}
                      className="-striped -highlight"
                    />
                  </CardBody>
                </Card>
              </GridItem>
            </GridContainer>
            {this.Alert(alert, t)}
          </div>
        )}
      </I18n>
    );
  }
}

WorkOrder.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  masterpcUrl: PropTypes.string.isRequired,
  workcenterCode: PropTypes.string.isRequired
};

const ConnWorkOrders = connect(
  mapStateToProps,
  mapDispatchToProps
)(WorkOrder);

export default withKeyboard(withStyles(styles)(ConnWorkOrders));
