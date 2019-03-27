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
import Dialog from '@material-ui/core/Dialog';

// react-spinners
import { GridLoader } from 'react-spinners';

// react-emotion
import { css } from 'react-emotion';

// I18n
import { I18n } from 'react-i18next';

// app
import connect from 'react-redux/es/connect/connect';
import AppBarBack from '../../components/AppBarBack';
import LeftMenuWithAvatar from '../../components/LeftMenuWithAvatar';
import styles from './styles';
import {
  fetchOperationListStart,
  fetchOperationDetailStart,
  editOperation
} from '../../actions/operationViewer';
import ImageEditor from '../../components/ImageEditor';
import { get } from 'lodash';

const override = css`
    display: block;
    margin: auto;
    border-color: red;
`;

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
    this.setContentsArray();
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

  renderImage = () => {
    const { contents, data, dispatchEditOperation } = this.props;
    if (!contents.image) {
      return null;
    }
    if (contents.image === 'editable') {
      return (
        < ImageEditor
          doEdit={({ points, img }) => {
            dispatchEditOperation(data.detail.id, points, img);
          }}
          points={data.detail ? data.detail.points || [] : []}
          img={data.detail.img}
        />
      );
    }
    return (
      <img
        alt="operation"
        src={data.detail.img}
        style={{
          height: '100%',
          width: '100%',
          objectFit: 'contain'
        }}
      />
    );
  };

  setContentsArray = () => {
    const { contents } = this.props;
    const contentsArray = [];
    if (contents.image) {
      contentsArray.push({
        label: 'Viewer.Image',
        content: this.renderImage
      });
    }
    if (contents.file) {
      contentsArray.push({
        label: 'Viewer.File',
        content: this.renderFile
      });
    }
    this.setState({
      contentsArray
    });
  };

  renderFile = () => {
    const { data, odooUrl } = this.props;
    const odooHost = new URL(odooUrl.value).host;
    return (
      <iframe
        title="worksheet"
        src={`http://${odooHost}/web/static/lib/pdfjs/web/viewer.html?file=%2Fapi%2Fv1%2Fworksheet%3Fmodel%3Dmrp.routing.workcenter%26field%3Dworksheet%26id%3D${
          data.detail.id}`
        }
        style={{
          width: '100%',
          height: '100%'
        }}
      />
    );
  };

  renderTabs = (t) => {
    const { classes } = this.props;
    const { currentTab, contentsArray } = this.state;
    return (
      <div className={classes.content}>
        <Tabs
          value={currentTab}
          onChange={this.handleChangeTab}
          style={{
            display: 'flex',
            flex: 1
          }}
        >
          {contentsArray.map((item,key) =>
            <Tab key={key} label={t(item.label)}/>
          )}
        </Tabs>
        <div style={{ display: 'flex', flex: 1, height: 'calc(100% - 45px)' }}>
          {contentsArray[currentTab].content()}
        </div>
      </div>
    );
  };

  render() {
    const { classes, data } = this.props;
    const { currentMenuItem } = this.state;
    return (
      <I18n ns="translations">
        {t => (
          <div className={classes.root}>
            <AppBarBack/>
            <LeftMenuWithAvatar>
              <MenuList>{this.genMenuList(t)}</MenuList>
            </LeftMenuWithAvatar>
            {currentMenuItem !== -1 ? this.renderTabs(t) : null}
            <Dialog fullScreen
                    classes={{
                      root: classes.loadModal
                    }}
                    open={data.loading}
                    style={{ opacity: 0.7 }}
                    TransitionComponent={this.Transition}
            >
              <GridLoader
                className={override}
                sizeUnit={'px'}
                size={50}
                color={'#36D7B7'}
                loading={data.loading}
              />
            </Dialog>
          </div>
        )}
      </I18n>
    );
  }
}

const mapStateToProps = (state, ownProps) => ({
  ...ownProps,
  data: state.operationViewer,
  odooUrl: state.setting.page.odooConnection.odooUrl,
  contents: state.setting.systemSettings.viewer
});

const mapDispatchToProps = {
  operationList: fetchOperationListStart,
  operationDetail: fetchOperationDetailStart,
  dispatchEditOperation: editOperation
};

const ConnectedViewer = connect(
  mapStateToProps,
  mapDispatchToProps
)(Viewer);

export default withStyles(styles)(ConnectedViewer);
