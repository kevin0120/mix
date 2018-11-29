import React from 'react';
// react component for creating dynamic tables
import ReactTable from 'react-table';
import PropTypes from 'prop-types';

// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';
// @material-ui/icons
import Assignment from '@material-ui/icons/Assignment';
import Dvr from '@material-ui/icons/Dvr';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

import Divider from '@material-ui/core/Divider';
import { I18n } from 'react-i18next';
import SweetAlert from 'react-bootstrap-sweetalert';

// core components
import GridContainer from '../../components/Grid/GridContainer';
import GridItem from '../../components/Grid/GridItem';
import Button from '../../components/CustomButtons/Button';
import Card from '../../components/Card/Card';
import CardBody from '../../components/Card/CardBody';
import CardIcon from '../../components/Card/CardIcon';
import CardHeader from '../../components/Card/CardHeader';
import Input from "@material-ui/core/Input";

import { Query, CreateDailyLogger, Warn } from '../../logger';

import {
  cardTitle,
  infoColor,
  warningColor,
  dangerColor
} from '../../common/jss/material-react-pro';
import withLayout from '../../components/Layout/layout';

import sweetAlertStyle from '../../common/jss/views/sweetAlertStyle';
import withKeyboard from '../../components/Keyboard';

const lodash = require('lodash');
const dayjs = require('dayjs');

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
  },
  InputRoot:{
    width: "100%" ,height:'36px',
    overflow:'hidden'
  },
  InputInput:{
    width: "100%" ,height:'100%',
  }
};

//
// Find items logged between today and yesterday.
//

const requestData = () =>
  new Promise((resolve, reject) => {
    // You can retrieve your data however you want, in this case, we will just use some local data.

    const options = {
      from: new Date() - 24 * 60 * 60 * 1000,
      until: new Date(),
      limit: 600,
      start: 0,
      order: 'desc',
      fields: ['timestamp', 'level', 'message']
    };

    Query(options, (err, results) => {
      if (err) {
        /* TODO: handle me */
        console.log(err.toString());
      }

      resolve(results);
    });
  });

class Event extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      data: [],
      loading: true,
      isShow: false,
      selectObj: null
    };
    this.fetchData = this.fetchData.bind(this);
  }

  componentDidMount() {
    CreateDailyLogger();
    this.fetchData();
  }

  fetchData(state, instance) {
    // Whenever the table model changes, or the user sorts or changes pages, this method gets called and passed the current table model.
    // You can set the `loading` prop of the table to true to use the built-in one or show you're own loading bar if you want.
    // Request the data however you want.  Here, we'll use our mocked service we created earlier
    requestData().then(res => {
      // Now just get the rows of data to your React Table (and update anything else like total pages or loading)
      this.setState({
        loading: false,
        data: res.info.map((item, key) => ({
          id: key,
          timestamp: dayjs(item.timestamp).format('YYYY MM-DD HH:mm:ss'),
          level: item.level,
          message: item.message,
          actions: (
            // we've added some custom button actions
            <div className="actions-right">
              {/* use this button to add a like kind of action */}
              <Button
                justIcon
                round
                simple
                onClick={() => {
                  const obj = this.state.data.find(o => o.id === key);
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
        }))
      });
    });
  }

  handleClose = () => {
    this.setState({
      isShow: false
    });
  };

  render() {
    const { classes } = this.props;
    const { data, isShow, selectObj, loading } = this.state;

    const Msg = selectObj ? (
      <div>
        <List>
          <ListItem>
            <ListItemText primary={`时间:   ${selectObj.timestamp}`} />
          </ListItem>
          <li>
            <Divider inset />
          </li>
          <ListItem>
            <ListItemText
              primary={`事件等级: ${
                selectObj.level === 'info' ? '日常信息' : '报警事件'
              }`}
            />
          </ListItem>
          <Divider inset component="li" />
          <ListItem>
            <ListItemText primary={`消息: ${selectObj.message}`} />
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
                  <CardHeader color="info" icon>
                    <CardIcon color="info">
                      <Assignment />
                    </CardIcon>
                    <h4 className={classes.cardIconTitle}>{t('main.event')}</h4>
                  </CardHeader>
                  <CardBody>
                    <ReactTable
                      loading={loading}
                      data={data}
                      getTrProps={(state, rowInfo) => {
                        if (rowInfo) {
                          let color = infoColor;
                          switch (rowInfo.row.level) {
                            case 'warn':
                              color = warningColor;
                              break;

                            case 'error':
                              color = dangerColor;
                              break;
                            default:
                              break;
                          }

                          return {
                            style: {
                              background: color
                            }
                          };
                        }

                        return {};
                      }}
                      filterable
                      columns={[
                        {
                          Header: 'Time',
                          accessor: 'timestamp',
                          filterable: false,
                          filterMethod: (filter, row) =>
                            lodash.includes(
                              lodash.toUpper(row[filter.id]),
                              lodash.toUpper(filter.value)
                            )
                        },
                        {
                          Header: 'Level',
                          accessor: 'level',
                          filterMethod: (filter, row) => {
                            if (filter.value === 'all') {
                              return true;
                            }
                            if (filter.value === 'info') {
                              return row[filter.id] === 'info';
                            }
                            if (filter.value === 'alarm') {
                              return row[filter.id] === 'warn';
                            }
                            if (filter.value === 'error') {
                              return row[filter.id] === 'error';
                            }
                            return true;
                          },
                          Filter: ({ filter, onChange }) => (
                            <select
                              onChange={event => onChange(event.target.value)}
                              style={{ width: '100%' }}
                              value={filter ? filter.value : 'all'}
                            >
                              <option value="all">All</option>
                              <option value="info">Info</option>
                              <option value="alarm">Alarm</option>
                              <option value="error">Error</option>
                            </select>
                          )
                        },
                        {
                          Header: 'Message',
                          accessor: 'message',
                          sortable: false,
                          filterMethod: (filter, row) => {
                            return lodash.includes(lodash.toUpper(row[filter.id]), lodash.toUpper(this.state.messageFilter||''));
                          },
                          Filter: ({ filter, onChange }) =>
                            <Input
                              onClick={()=>{
                                this.props.keyboardInput({
                                  onSubmit:(text)=>{
                                    this.setState({messageFilter:text},()=>{
                                      onChange(this.state.messageFilter);
                                    });
                                  },
                                  text:this.state.messageFilter,
                                  title:'Message',
                                  label:'Message'
                                });
                              }}
                              classes={{
                                root:classes.InputRoot,
                                input:classes.InputInput,
                              }}
                              // style={{ width: "100%" ,height:'36px'}}
                              value={this.state.messageFilter || ""}
                            />
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
                warning
                show={isShow}
                style={{ display: 'block', marginTop: '-100px' }}
                title="事件详情"
                onConfirm={this.handleClose}
                onCancel={this.handleClose}
                confirmBtnCssClass={`${this.props.classes.button} ${
                  this.props.classes.success
                }`}
                cancelBtnCssClass={`${this.props.classes.button} ${
                  this.props.classes.danger
                }`}
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

Event.propTypes = {
  classes: PropTypes.shape({}).isRequired
};

export default withLayout(withKeyboard(withStyles(styles)(Event)), false);
