import React from 'react';
import cx from 'classnames';
import PropTypes from 'prop-types';

// @material-ui/core components
import withStyles from '@material-ui/core/styles/withStyles';

// core components
import { Scrollbars } from 'react-custom-scrollbars';
import Badge from '../Badge/Badge';

import timelineStyle from '../../common/jss/components/timelineStyle';

function Timeline({ ...props }) {
  const { classes, stories, simple } = props;
  const timelineClass = `${classes.timeline} ${cx({
    [classes.timelineSimple]: simple
  })}`;
  return (
    <Scrollbars>
      <ul className={timelineClass}>
        {stories.map((prop, key) => {
          const panelClasses = `${classes.timelinePanel} ${cx({
            [classes.timelinePanelInverted]: prop.inverted,
            [classes.timelineSimplePanel]: simple
          })}`;
          const timelineBadgeClasses = `${classes.timelineBadge} ${
            classes[prop.badgeColor]
          } ${cx({
            [classes.timelineSimpleBadge]: simple
          })}`;
          return (
            <li className={classes.item} key={key}>
              {prop.badgeIcon ? (
                <div className={timelineBadgeClasses}>
                  <prop.badgeIcon className={classes.badgeIcon} />
                </div>
              ) : null}
              <div className={panelClasses}>
                {prop.title ? (
                  <div className={classes.timelineHeading}>
                    <Badge color={prop.titleColor}>{prop.title}</Badge>
                  </div>
                ) : null}
                <div className={classes.timelineBody}>{prop.body}</div>
                {prop.footerTitle ? (
                  <h6 className={classes.footerTitle}>{prop.footerTitle}</h6>
                ) : null}
                {prop.footer ? <hr className={classes.footerLine} /> : null}
                {prop.footer ? (
                  <div className={classes.timelineFooter}>{prop.footer}</div>
                ) : null}
              </div>
            </li>
          );
        })}
      </ul>
    </Scrollbars>
  );
}

Timeline.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  stories: PropTypes.array.isRequired,
  simple: PropTypes.bool
};

export default withStyles(timelineStyle)(Timeline);
