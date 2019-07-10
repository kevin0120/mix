import React from 'react';
import { connect } from 'react-redux';
import { withStyles } from '@material-ui/core/styles';
import styles from './styles';

type Props = {
  classes: {}
};

class scannerStep extends React.Component<Props> {

// eslint-disable-next-line flowtype/no-weak-types

  render(): React.ReactElement<any> | string | number | {} | React.ReactNodeArray | React.ReactPortal | boolean | null | undefined {
    const { classes } = this.props;
    return <div className={classes.root} >
      {'code icon here'}
    </div>;
  }
}

const mapState = (state,props) => ({
  ...props
});

const mapDispatch = {};

export default withStyles(styles)(connect(mapState, mapDispatch)(scannerStep));
