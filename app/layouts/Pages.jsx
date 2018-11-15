import React from 'react';
import { Switch, Route, Redirect } from 'react-router-dom';

import PropTypes from 'prop-types';

// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';

// core components
import PagesHeader from '../components/Header/PagesHeader';
import Footer from '../components/Footer/Footer';

import pagesRoutes from '../routes/pages';

import pagesStyle from '../common/jss/layouts/pagesStyle';

import bgImage from '../../resources/imgs/lock.jpeg';

// var ps;

class Pages extends React.Component {
  componentDidMount() {
    document.body.style.overflow = 'unset';
  }

  render() {
    const { classes, ...rest } = this.props;

    const anchor = `${classes.a} ${classes.whiteColor}`;
    return (
      <div>
        <PagesHeader {...rest} />
        <div className={classes.wrapper}>
          <div
            className={classes.fullPage}
            style={{ backgroundImage: `url(${bgImage})` }}
          >
            <Switch>
              {pagesRoutes.map(prop => {
                if (prop.collapse) {
                  return null;
                }
                if (prop.redirect) {
                  return (
                    <Redirect
                      from={prop.path}
                      to={prop.pathTo}
                      key={prop.path}
                    />
                  );
                }
                return (
                  <Route
                    path={prop.path}
                    component={prop.component}
                    key={prop.path}
                  />
                );
              })}
            </Switch>
            <Footer
              white
              content={
                <p className={classes.right}>
                  &copy; {1900 + new Date().getYear()}{' '}
                  <a href="https://www.liktek.com" className={anchor}>
                    上海砺星信息科技技术有限公司
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
  classes: PropTypes.object.isRequired
};

export default withStyles(pagesStyle)(Pages);
