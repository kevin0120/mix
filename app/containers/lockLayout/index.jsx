import React from 'react';

import PropTypes from 'prop-types';

// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';


import { connect } from 'react-redux';
// core components
import PagesHeader from '../../components/Header/PagesHeader';
import Footer from '../../components/Footer/Footer';
import styles from './styles';
import bgImage from '../../../resources/imgs/lock.jpeg';
import { loginRequest } from '../../modules/user/action';

class Pages extends React.Component {
  componentDidMount() {
    document.body.style.overflow = 'unset';
  }

  render() {
    const { classes, children, login, ...rest } = this.props;
    const anchor = `${classes.a} ${classes.whiteColor}`;
    return (
      <div>
        <PagesHeader login={() => login('userName', 'password', 'local')} {...rest} />
        <div className={classes.wrapper}>
          <div
            className={classes.fullPage}
            style={{ backgroundImage: `url(${bgImage})` }}
          >
            {children}
            <Footer
              white
              content={
                <p className={classes.right}>
                  &copy; {1900 + new Date().getYear()}{' '}
                  <a href="https://www.liktek.com" className={anchor}>
                    上海途泰信息科技技术有限公司
                  </a>
                </p>
              }
            />
          </div>
        </div>
      </div>
    );
  }
}

Pages.propTypes = {
  classes: PropTypes.shape({}).isRequired
};

const mapState = (state, props) => {
  return {
    ...props

  };
};
const mapDispatch = {
  login: loginRequest
};

export default withStyles(styles)(connect(mapState, mapDispatch)(Pages));
