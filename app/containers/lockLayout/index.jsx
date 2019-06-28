import React from 'react';

import PropTypes from 'prop-types';

// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';

// core components
import PagesHeader from '../../components/Header/PagesHeader';
import Footer from '../../components/Footer/Footer';

import pagesStyle from '../../common/jss/layouts/pagesStyle';

import bgImage from '../../../resources/imgs/lock.jpeg';

// var ps;

class Pages extends React.Component {
  componentDidMount() {
    document.body.style.overflow = 'unset';
  }

  render() {
    const { classes, children, ...rest } = this.props;
    console.log(children);
    const anchor = `${classes.a} ${classes.whiteColor}`;
    return (
      <div>
        <PagesHeader {...rest} />
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

export default withStyles(pagesStyle)(Pages);
