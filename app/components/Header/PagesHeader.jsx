import React from 'react';
import cx from 'classnames';
import PropTypes from 'prop-types';
import { NavLink } from 'react-router-dom';

// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Hidden from '@material-ui/core/Hidden';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';

// @material-ui/icons
import Dashboard from '@material-ui/icons/Dashboard';
import Signal0 from '@material-ui/icons/SignalWifi0Bar';
import Signal1 from '@material-ui/icons/SignalWifi1Bar';
import Signal2 from '@material-ui/icons/SignalWifi2Bar';
import Signal3 from '@material-ui/icons/SignalWifi3Bar';
import Signal4 from '@material-ui/icons/SignalWifi4Bar';

import Battery20 from '@material-ui/icons/Battery20';
import Battery30 from '@material-ui/icons/Battery30';
import Battery50 from '@material-ui/icons/Battery50';
import Battery60 from '@material-ui/icons/Battery60';
import Battery80 from '@material-ui/icons/Battery80';
import Battery90 from '@material-ui/icons/Battery90';
import BatteryAlert from '@material-ui/icons/BatteryAlert';
import BatteryFull from '@material-ui/icons/BatteryFull';
import BatteryUnknown from '@material-ui/icons/BatteryUnknown';
import Menu from '@material-ui/icons/Menu';

// core components
import Button from '../CustomButtons/Button';

import pagesHeaderStyle from '../../common/jss/components/pagesHeaderStyle';
import connect from 'react-redux/es/connect/connect';
import { networkCheck, networkSignal } from '../../modules/network/action';
import { batteryCheck } from '../../modules/battery/action';

const signalLevel = (signal) => {
  const size = 'large';
  if (signal > 80)
    return <Signal4 fontSize={size}/>;
  else if (signal > 60)
    return <Signal3 fontSize={size}/>;
  else if (signal > 40)
    return <Signal2 fontSize={size}/>;
  else if (signal > 20)
    return <Signal1 fontSize={size}/>;
  else
    return <Signal0 fontSize={size}/>;
};

const batteryLevel = (percentage) => {
  const size = 'large';
  if (percentage > 99)
    return <BatteryFull fontSize={size}/>;
  else if (percentage >= 90)
    return <Battery90 fontSize={size}/>;
  else if (percentage >= 80)
    return <Battery80 fontSize={size}/>;
  else if (percentage >= 60)
    return <Battery60 fontSize={size}/>;
  else if (percentage >= 50)
    return <Battery50 fontSize={size}/>;
  else if (percentage >= 30)
    return <Battery30 fontSize={size}/>;
  else if (percentage >= 20)
    return <Battery20 fontSize={size}/>;
  else if (percentage >= 0)
    return <BatteryAlert fontSize={size}/>;
  else
    return <BatteryUnknown fontSize={size}/>;
};

class PagesHeader extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      open: false
    };
  }

  componentDidMount(): void {
    // check wifi signal
    this.handleStatusCheck();
  }

  handleStatusCheck=()=>{
    const { doNetworkSignal, doBatteryCheck } = this.props;
    doNetworkSignal();
    doBatteryCheck();
  };

  handleDrawerToggle = () => {
    const { open } = this.state;
    this.setState({ open: !open });
  };

  // verifies if routeName is the one active (in browser input)
  activeRoute(routeName) {
    const { location } = this.props;
    return location.pathname.indexOf(routeName) > -1;
  }

  // componentDidUpdate(e) {
  //   if (e.history.location.pathname !== e.location.pathname) {
  //     this.setState({open: false});
  //   }
  // }

  render() {
    const { classes, color, ssid, signal, batteryPercentage } = this.props;
    const { open } = this.state;
    const appBarClasses = cx({
      [` ${classes[color]}`]: color
    });
    const list = (
      <List className={classes.list}>
        <ListItem className={classes.listItem}>
          <NavLink to="/app" className={classes.navLink}>
            <ListItemIcon className={classes.listItemIcon}>
              <Dashboard/>
            </ListItemIcon>
            <ListItemText
              primary="Welcome"
              disableTypography
              className={classes.listItemText}
            />
          </NavLink>
        </ListItem>
      </List>
    );
    return (
      <AppBar position="static" className={classes.appBar + appBarClasses}>
        <Toolbar className={classes.container}>
          <div style={{
            display: 'flex',
            flexDirection: 'row',
            justifyContent: 'space-between',
            alignItems: 'center',
            padding: 0,
            flex: 1
          }}>
            <div style={{
              display: 'flex',
              flexDirection: 'row',
              alignItems: 'center'
            }}>
              <Hidden smDown implementation="css">
                <div className={classes.flex}>
                  <Button className={classes.title} color="transparent">
                    Control Panel
                  </Button>
                </div>
              </Hidden>
              {/* <Hidden mdUp> */}
              {/* <div className={classes.flex}> */}
              {/* <Button className={classes.title} color="transparent"> */}
              {/* For Smart Assembly */}
              {/* </Button> */}
              {/* </div> */}
              {/* </Hidden> */}
              <Hidden smDown implementation="css">
                {list}
              </Hidden>
            </div>
            <Button className={classes.indicator} color="transparent" onClick={this.handleStatusCheck}>
              {signalLevel(signal)}
              <span style={{marginRight:'7px'}}>{`${ssid || '无连接'}`}</span>
              {batteryLevel(batteryPercentage)}
              <span>{batteryPercentage>=0?`${batteryPercentage}%` : '电池检测中'}</span>
            </Button>
          </div>
          {/* <Hidden mdUp> */}
          {/* <Button */}
          {/* className={classes.sidebarButton} */}
          {/* color="transparent" */}
          {/* justIcon */}
          {/* aria-label="open drawer" */}
          {/* onClick={this.handleDrawerToggle} */}
          {/* > */}
          {/* <Menu/> */}
          {/* </Button> */}
          {/* </Hidden> */}
          {/* <Hidden mdUp implementation="css"> */}
          {/* <Hidden mdUp> */}
          {/* <Drawer */}
          {/* variant="temporary" */}
          {/* anchor="right" */}
          {/* open={open} */}
          {/* classes={{ */}
          {/* paper: classes.drawerPaper */}
          {/* }} */}
          {/* onClose={this.handleDrawerToggle} */}
          {/* ModalProps={{ */}
          {/* keepMounted: true // Better open performance on mobile. */}
          {/* }} */}
          {/* > */}
          {/* {list} */}
          {/* </Drawer> */}
          {/* </Hidden> */}
          {/* </Hidden> */}
        </Toolbar>
      </AppBar>
    );
  }
}

PagesHeader.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  color: PropTypes.oneOf(['primary', 'info', 'success', 'warning', 'danger'])
};

const mapStateToProps = (state) => ({
  ssid: state.network.ssid,
  signal: state.network.signal,
  batteryPercentage: state.battery.percentage
});

const mapDispatchToProps = {
  doNetworkCheck: networkCheck,
  doNetworkSignal: networkSignal,
  doBatteryCheck: batteryCheck
};

const ConnectedPagesHeader = connect(
  mapStateToProps,
  mapDispatchToProps
)(PagesHeader);

export default withStyles(pagesHeaderStyle)(ConnectedPagesHeader);

// WEBPACK FOOTER //
// ./src/components/Header/PagesHeader.jsx
