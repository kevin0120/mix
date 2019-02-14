// React
import React from 'react';

// material-ui
import { withStyles } from '@material-ui/core';
import MenuList from '@material-ui/core/MenuList';
import MenuItem from '@material-ui/core/MenuItem';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';

// I18n
import { I18n } from 'react-i18next';

// app
import connect from 'react-redux/es/connect/connect';
import AppBarBack from '../../components/AppBarBack';
import LeftMenuWithAvatar from '../../components/LeftMenuWithAvatar';
import styles from './styles';
import {
  fetchOperationListStart,
  fetchOperationDetailStart
} from '../../actions/operationViewer';

class Viewer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      currentMenuItem: -1,
      currentTab: 0
    };
  }

  componentDidMount(): void {
    const { operationList } = this.props;
    operationList();
  }

  handleChangeMenu = (index, operationID) => {
    const { operationDetail } = this.props;
    this.setState({ currentMenuItem: index });
    operationDetail(operationID);
  };

  handleChangeTab = (event, value) => {
    this.setState({ currentTab: value });
  };

  genMenuList = t => {
    const { classes, data } = this.props;
    const { currentMenuItem } = this.state;
    const menus = data.list.map((item, index) => {
      const selected = currentMenuItem === index;
      return (
        <MenuItem
          selected={selected}
          key={item.id}
          className={classes.menuItem}
          component={() => (
            <Card
              onClick={() => this.handleChangeMenu(index, item.id)}
              className={selected ? classes.menuItemSelected : classes.menuItem}
            >
              <CardActionArea
                classes={{
                  root: selected
                    ? classes.cardActionAreaSelected
                    : classes.cardActionArea
                }}
              >
                <span className={classes.itemText}>{t(item.name)}</span>
              </CardActionArea>
            </Card>
          )}
        />
      );
    });
    return <div>{menus}</div>;
  };

  render() {
    const { classes, data, odooUrl } = this.props;
    const { currentMenuItem, currentTab } = this.state;
    const odooHost = new URL(odooUrl.value).host;
    return (
      <I18n ns="translations">
        {t => (
          <div className={classes.root}>
            <AppBarBack />
            <LeftMenuWithAvatar>
              <MenuList>{this.genMenuList(t)}</MenuList>
            </LeftMenuWithAvatar>
            {currentMenuItem !== -1 ? (
              <div className={classes.content}>
                <Tabs
                  value={currentTab}
                  onChange={this.handleChangeTab}
                  style={{
                    display: 'flex',
                    flex: 1,
                  }}
                >
                  <Tab label="Image" />
                  <Tab label="File" />
                </Tabs>
                <div style={{display:'flex',flex:1,height:'calc(100% - 45px)'}}>
                {currentTab === 0 && (
                  <img
                    alt="operation image"
                    src={data.detail.img}
                    style={{
                      display: 'flex',
                      flex: 1,
                      maxHeight: '100%',
                      maxWidth: '100%'
                    }}
                  />
                )}
                {currentTab === 1 && (
                  <iframe
                    title="worksheet"
                    src={`http://${odooHost}/web/static/lib/pdfjs/web/viewer.html?file=%2Fapi%2Fv1%2Fworksheet%3Fmodel%3Dmrp.routing.workcenter%26field%3Dworksheet%26id%3D${
                      data.detail.id
                    }`}
                    // file={pdf}
                    // src={`/home/cosine/文档/001.课程/大四上/储能技术/储能技术的核心问题.pdf`}
                    style={{
                      width: '100%',
                      height: '100%'
                    }}
                    // plugins
                  />
                )}
                </div>
              </div>
            ) : null}
          </div>
        )}
      </I18n>
    );
  }
}

const mapStateToProps = (state, ownProps) => ({
  ...ownProps,
  data: state.operationViewer,
  odooUrl: state.setting.page.odooConnection.odooUrl
});

const mapDispatchToProps = {
  operationList: fetchOperationListStart,
  operationDetail: fetchOperationDetailStart
};

const ConnectedViewer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Viewer);

export default withStyles(styles)(ConnectedViewer);
