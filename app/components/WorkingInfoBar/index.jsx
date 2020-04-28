import React from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';

import { withStyles } from '@material-ui/core/styles';

import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';

import styles from './styles';

const mapStateToProps = (state, ownProps) => ({
  infos: ownProps.infos
});

const mapDispatchToProps = {};

function ConnectedWorkingInfoBar(props) {
  const { classes, styleOptions, infos } = props;

  return infos.map(info => {
    const { key, value, displayTitle } = info;
    if (key === 'taotong') {
      return;
    }
    return (
      <ListItem
        key={key}
        disableGutters={styleOptions.disableGutters}
        className={classes.infoItem}
      >
        <ListItemText className={classes.itemIitle} primary={displayTitle} />

        <span className={classes.infoText}>{value}</span>
      </ListItem>
    );
  });
}

ConnectedWorkingInfoBar.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  styleOptions: PropTypes.shape({
    disableGutters: PropTypes.bool
  }),
  options: PropTypes.shape({
    enableDebugLog: PropTypes.bool
  }),
  infos: PropTypes.arrayOf(
    PropTypes.shape({
      key: PropTypes.string.isRequired,
      value: PropTypes.all,
      displayTitle: PropTypes.string.isRequired
    })
  ).isRequired
};

ConnectedWorkingInfoBar.defaultProps = {
  styleOptions: {
    disableGutters: false
  },
  options: {
    enableDebugLog: false
  }
};

const WorkingInfoBar = connect(
  mapStateToProps,
  mapDispatchToProps
)(ConnectedWorkingInfoBar);

export default withStyles(styles)(WorkingInfoBar);
