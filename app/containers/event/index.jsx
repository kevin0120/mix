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
import Alert from '../../components/Alert';

// core components
import Input from '@material-ui/core/Input';
import GridContainer from '../../components/Grid/GridContainer';
import GridItem from '../../components/Grid/GridItem';
import Button from '../../components/CustomButtons/Button';
import Card from '../../components/Card/Card';
import CardBody from '../../components/Card/CardBody';
import CardIcon from '../../components/Card/CardIcon';
import CardHeader from '../../components/Card/CardHeader';
import CustomReactTable from '../../components/CustomReactTable';

import { Query, CreateDailyLogger, Warn } from '../../logger';

import sweetAlertStyle from '../../common/jss/views/sweetAlertStyle';
import withKeyboard from '../../components/Keyboard';

const lodash = require('lodash');
const dayjs = require('dayjs');

const styles = theme => ({
  ...sweetAlertStyle(theme),
  root: {
    flexGrow: 1,
    zIndex: 1,
    height: '100%',
    overflowY: 'auto',
    position: 'relative',
    display: 'flex',
    width: '100%',
    flexDirection: 'column'
  },
  page: {
    flex:1,
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
  infoTr: {
    backgroundColor: theme.palette.info
  },
  warnTr: {
    backgroundColor: theme.palette.warning
  },
  maintenanceTr: {
    backgroundColor: theme.palette.warning
  },
  errorTr: {
    backgroundColor: theme.palette.danger
  }
});

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
      fields: ['timestamp', 'level', 'message', 'meta']
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

  fetchData() {
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
          meta: item.meta,
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
                color="rose"
                className="edit"
              >
                <Dvr/>
              </Button>
            </div>
          )
        }))
      });
    }).catch((e) => console.error(e));
  }

  handleClose = () => {
    this.setState({
      isShow: false
    });
  };

  renderMsg = (selectObj) => {
    if (!selectObj) {
      return ' ';
    }

    let levelText = 'Event.Info';
    switch (selectObj.level) {
      case 'info':
        levelText = 'Event.Info';
        break;
      case 'maintenance':
        levelText = 'Event.Maintenance';
        break;
      default:
        levelText = 'Event.Alert';
        break;
    }
    const keys = Object.keys(selectObj.meta || {});
    return (
      <I18n ns="translations">
        {t => (
          <div>
            <List>
              <ListItem>
                <ListItemText primary={`${t('Event.Time')}:   ${selectObj.timestamp}`}/>
              </ListItem>
              <li>
                <Divider/>
              </li>
              <ListItem>
                <ListItemText
                  primary={`${t('Event.Level')}: ${t(levelText)}`}
                />
              </ListItem>
              <Divider component="li"/>
              <ListItem>
                <ListItemText primary={`${t('Event.Message')}: ${selectObj.message}`}/>
              </ListItem>
              {
                keys.map((key) => (

                  <React.Fragment>
                    <Divider component="li"/>
                    <ListItem>
                      <ListItemText primary={`${t(key)}: ${selectObj.meta[key]}`}/>
                    </ListItem>
                  </React.Fragment>
                ))
              }
            </List>
          </div>
        )}
      </I18n>
    );
  };

  render() {
    const { classes, keyboardInput } = this.props;
    const { data, isShow, selectObj, loading, messageFilter } = this.state;

    return (
      <I18n ns="translations">
        {t => (
          <React.Fragment>
            <div className={classes.root}>
              <div className={classes.page}>
                <Card>
                  <CardHeader color="info" icon>
                    <CardIcon color="info">
                      <Assignment/>
                    </CardIcon>
                    <h4 className={classes.cardIconTitle}>{t('main.event')}</h4>
                  </CardHeader>
                  <CardBody>
                    <CustomReactTable
                      translate={t}
                      showPageJump={false}
                      loading={loading}
                      data={data}
                      getTrProps={(state, rowInfo) => {
                        if (rowInfo) {
                          let className = classes.infoTr;
                          switch (rowInfo.row.level) {
                            case 'warn':
                              className = classes.warnTr;
                              break;
                            case 'maintenance':
                              className = classes.maintenanceTr;
                              break;
                            case 'error':
                              className = classes.errorTr;
                              break;
                            default:
                              break;
                          }
                          return {
                            className
                          };
                        }
                        return {};
                      }}
                      filterable
                      columns={[
                        {
                          Header: t('Event.Time'),
                          accessor: 'timestamp',
                          filterable: false,
                          filterMethod: (filter, row) =>
                            lodash.includes(
                              lodash.toUpper(row[filter.id]),
                              lodash.toUpper(filter.value)
                            )
                        },
                        {
                          Header: t('Event.Level'),
                          accessor: 'level',
                          filterMethod: (filter, row) => {
                            if (filter.value === 'all') {
                              return true;
                            }
                            if (filter.value === 'info') {
                              return row[filter.id] === 'info';
                            }
                            if (filter.value === 'warn') {
                              return row[filter.id] === 'warn';
                            }
                            if (filter.value === 'maintenance') {
                              return row[filter.id] === 'maintenance';
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
                              <option value="all">{t('Event.All')}</option>
                              <option value="info">{t('Event.Info')}</option>
                              <option value="maintenance">{t('Event.Maintenance')}</option>
                              <option value="warn">{t('Event.Warn')}</option>
                              <option value="error">{t('Event.Error')}</option>
                            </select>
                          )
                        },
                        {
                          Header: t('Event.Message'),
                          accessor: 'message',
                          sortable: false,
                          filterable: false
                          // filterMethod: (filter, row) => lodash.includes(
                          //     lodash.toUpper(row[filter.id]),
                          //     lodash.toUpper(messageFilter || '')
                          //   ),
                          // Filter: ({ onChange }) => (
                          //   <Input
                          //     onClick={() => {
                          //       keyboardInput({
                          //         onSubmit: text => {
                          //           this.setState(
                          //             { messageFilter: text },
                          //             () => {
                          //               onChange(messageFilter);
                          //             }
                          //           );
                          //         },
                          //         text: messageFilter,
                          //         title: 'Message',
                          //         label: 'Message'
                          //       });
                          //     }}
                          //     classes={{
                          //       root: classes.InputRoot,
                          //       input: classes.InputInput
                          //     }}
                          //     // style={{ width: "100%" ,height:'36px'}}
                          //     value={messageFilter || ''}
                          //   />
                          // )
                        },
                        {
                          Header: t('Event.Actions'),
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
              </div>
            </div>
            <Alert
              warning
              show={isShow}
              title={t('Event.Detail')}
              onConfirm={this.handleClose}
              onCancel={this.handleClose}
              confirmBtnCssClass={`${classes.button} ${
                classes.successWarn
                }`}
              cancelBtnCssClass={`${classes.button} ${
                classes.danger
                }`}
              confirmBtnText={t('Common.Yes')}
              cancelBtnText={t('Common.No')}
              showCancel
            >
              {this.renderMsg(selectObj)}
            </Alert>
          </React.Fragment>
        )}
      </I18n>
    );
  }
}

Event.propTypes = {
  classes: PropTypes.shape({}).isRequired
};

export default withKeyboard(withStyles(styles)(Event));
