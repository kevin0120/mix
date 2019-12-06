import React from 'react';
import CardActionArea from '@material-ui/core/CardActionArea';
import clsx from 'clsx';
import CardMedia from '@material-ui/core/CardMedia';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import { makeStyles } from '@material-ui/core';
import { connect } from 'react-redux';
import settingImg from '../../../resources/imgs/setting.png';
import styles from './styles';

const mapState = (state, ownProps) => ownProps;

const mapDispatch = (selectItem) => {
  if (selectItem) {
    return { selectItem };
  }
  return {};
};


function SelectCard(props) {

  const { height, status, selectItem, onClick, item, statusClass, image, name, infoArr } = props;
  const classes = makeStyles(styles(height))();
  return <Paper square className={classes.orderCardContainer}>
    <CardActionArea
      className={classes.orderCard}
      onClick={() => {
        if (onClick) {
          onClick(item);
        }
        if (selectItem) {
          selectItem(item);
        }
      }}
    >
      <div
        className={clsx(
          statusClass,
          classes.statusIndicator
        )}
      />
      <CardMedia
        className={classes.image}
        src={image || settingImg}
        component="img"
        style={{ overflow: 'hidden' }}
      />
      <CardContent className={classes.info}>
        <Typography
          variant="body1"
          align="left"
          className={classes.orderNameText}
        >
          {name}
        </Typography>
        {infoArr && infoArr.map((i, idx) => (
          <Typography
            key={`${i}-${idx}`}
            variant="body2"
            color="textSecondary"
            align="left"
            className={classes.orderInfoText}
          >
            {i}
          </Typography>
        ))}
        <Typography
          variant="body2"
          color="textSecondary"
          align="left"
          className={classes.orderStatusText}
        >
          {status}
        </Typography>
      </CardContent>
    </CardActionArea>
  </Paper>;
}

export default (selectItem) => connect(mapState, mapDispatch(selectItem))(SelectCard);

