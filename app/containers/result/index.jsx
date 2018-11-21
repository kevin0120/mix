import React from 'react';
import { connect } from 'react-redux';
import ReactTable from 'react-table';

import Divider from '@material-ui/core/Divider';

import { I18n } from 'react-i18next';

import SweetAlert from 'react-bootstrap-sweetalert';
import { cardTitle } from '../../common/jss/material-react-pro';
import withLayout from '../../components/Layout/layout';

import sweetAlertStyle from '../../common/jss/views/sweetAlertStyle';

import axios from 'axios';
import axiosRetry from 'axios-retry';

// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';
// @material-ui/icons
import Assignment from '@material-ui/icons/Assignment';
import Dvr from '@material-ui/icons/Dvr';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
// core components
import GridContainer from '../../components/Grid/GridContainer';
import GridItem from '../../components/Grid/GridItem';
import Button from '../../components/CustomButtons/Button';
import Card from '../../components/Card/Card';
import CardBody from '../../components/Card/CardBody';
import CardIcon from '../../components/Card/CardIcon';
import CardHeader from '../../components/Card/CardHeader';

const lodash = require('lodash');
const dayjs = require('dayjs');

const defaultInstance = axios.create({
  timeout: 3000,
  headers: { 'Content-Type': 'application/json' }
});

axiosRetry(defaultInstance, {
  retries: 2,
  retryDelay: axiosRetry.exponentialDelay,
  retryCondition: err => {
    if (err.message.indexOf('200') !== -1) {
      return false;
    }

    return true;
  }
});

const styles = {
  ...sweetAlertStyle,
  root: {
    flexGrow: 1,
    zIndex: 1,
    height: 'calc(100% - 64px)',
    overflow: 'hidden',
    position: 'relative',
    display: 'flex',
    width: '100%'
  },
  cardIconTitle: {
    ...cardTitle,
    marginTop: '15px',
    marginBottom: '0px'
  }
};

const mapStateToProps = (state, ownProps) => ({
  masterpcUrl: state.connections.masterpc,
  hmiSn: state.setting.page.odooConnection.hmiSn.value,
  ...ownProps
});

const mapDispatchToProps = {};

function requestData(masterpcUrl, hmiSN) {
  const url = `${masterpcUrl}/rush/v1/local-results`;
  return defaultInstance.get(url, {
    params: {
      hmi_sn: hmiSN,
      filters: 'vin,job_id,batch,torque,angle,timestamp,vehicle_type',
      limit: 500
    }
  });
}

class Result extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      isShow: false,
      selectObj: null
    };
    this.fetchData = this.fetchData.bind(this);
  }

  componentDidMount() {
    this.fetchData();
  }

  fetchData() {
    const { masterpcUrl, hmiSn } = this.props;
    requestData(masterpcUrl, hmiSn)
      .then(res => {
        const statusCode = res.status;
        if (statusCode === 200) {
          this.setState({
            data: res.data.map((item, key) => {
              return {
                id: key,
                timestamp: dayjs(item.timestamp).format('YYYY MM-DD HH:mm:ss'),
                vin: item.vin,
                torque: item.torque,
                angle: item.angle,
                job_id: item.job_id,
                batch: item.batch,
                vehicle_type: item.vehicle_type,
                actions: (
                  // we've added some custom button actions
                  <div className="actions-right">
                    {/* use this button to add a like kind of action */}
                    <Button
                      justIcon
                      round
                      simple
                      onClick={() => {
                        let obj = this.state.data.find(o => o.id === key);
                        this.setState({
                          isShow: true,
                          selectObj: obj
                        });
                      }}
                      color="warning"
                      className="edit"
                    >
                      <Dvr />
                    </Button>{' '}
                  </div>
                )
              };
            })
          });
        }
      })
      .catch(error => {
        console.log('get error' + error.toString());
      });
  }

  handleClose = () => {
    this.setState({
      isShow: false
    });
  };

  render() {
    const { classes } = this.props;
    const { data, isShow, selectObj } = this.state;

    const Msg = selectObj ? (
      <div>
        <List>
          <ListItem>
            <ListItemText primary={`VIN:   ${selectObj.vin}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`车型:   ${selectObj.vehicle_type}`} />
          </ListItem>
          <li>
            <Divider inset />
          </li>
          <ListItem>
            <ListItemText primary={`扭矩: ${selectObj.torque}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`角度: ${selectObj.angle}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`批次:   ${selectObj.batch}`} />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`拧紧时间:   ${selectObj.timestamp}`} />
          </ListItem>
        </List>
      </div>
    ) : (
      ' '
    );

    return (
      <I18n ns="translations">
        {t => (
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
                    <ReactTable
                      data={data}
                      filterable
                      columns={[
                        {
                          Header: 'VIN',
                          accessor: 'vin',
                          filterMethod: (filter, row) =>
                            lodash.includes(
                              lodash.toUpper(row[filter.id]),
                              lodash.toUpper(filter.value)
                            )
                        },
                        {
                          Header: '车型',
                          accessor: 'vehicle_type',
                          filterMethod: (filter, row) =>
                            lodash.includes(
                              lodash.toUpper(row[filter.id]),
                              lodash.toUpper(filter.value)
                            )
                        },
                        {
                          Header: '程序号',
                          accessor: 'job_id',
                          sortable: false
                        },
                        {
                          Header: '扭矩',
                          accessor: 'torque',
                          sortable: false,
                          filterable: false
                        },
                        {
                          Header: '角度',
                          accessor: 'angle',
                          sortable: false,
                          filterable: false
                        },
                        {
                          Header: '批次',
                          accessor: 'batch',
                          sortable: false,
                          filterable: false
                        },
                        {
                          Header: '拧紧时间',
                          accessor: 'timestamp',
                          filterable: false,
                          filterMethod: (filter, row) =>
                            lodash.includes(
                              lodash.toUpper(row[filter.id]),
                              lodash.toUpper(filter.value)
                            )
                        },
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
            {isShow ? (
              <SweetAlert
                info
                show={isShow}
                style={{ display: 'block', marginTop: '-100px' }}
                title="事件详情"
                onConfirm={this.handleClose}
                onCancel={this.handleClose}
                confirmBtnCssClass={
                  this.props.classes.button + ' ' + this.props.classes.success
                }
                cancelBtnCssClass={
                  this.props.classes.button + ' ' + this.props.classes.danger
                }
                confirmBtnText={t('Common.Yes')}
                cancelBtnText={t('Common.No')}
                showCancel
              >
                {Msg}
              </SweetAlert>
            ) : null}
          </div>
        )}
      </I18n>
    );
  }
}

const ConnResult = connect(
  mapStateToProps,
  mapDispatchToProps
)(Result);

export default withLayout(withStyles(styles)(ConnResult), false);
