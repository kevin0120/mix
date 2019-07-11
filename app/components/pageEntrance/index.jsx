import React from 'react';
import BottomNavigationAction from '@material-ui/core/BottomNavigationAction';
import BottomNavigation from '@material-ui/core/BottomNavigation';
import { I18n } from 'react-i18next';
import Grid from '@material-ui/core/Grid';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardContent from '@material-ui/core/CardContent';
import { withStyles } from '@material-ui/core/styles';
import { cardStyles } from './styles';


const renderCard = withStyles(cardStyles)((props) => {
  const { routes, onItemClick, classes } = props;
  return <I18n ns="translations">
    {t => (<Grid container className={classes.container} justify="center">
      {routes.map(route => route ? (
        <Grid key={route.name} item className={classes.cardGridItem}>
          <Card
            key={route.name}
            className={classes.card}
            style={{ backgroundColor: route.color }}
          >
            <CardActionArea
              onClick={() => onItemClick(route)}
              className={classes.cardActionArea}
            >
              <div
                className={classes.media}
                style={{ backgroundImage: `url(${route.image})`, backgroundSize: 'cover' }}
              />
              <CardContent className={classes.cardContent}>
                <div className={classes.iconWrap}>
                  <route.icon className={classes.icon}/>
                </div>
                <h1 className={classes.title}>{t(route.title)}</h1>
                <p className={classes.subTitle}>
                  {t(route.title, { lng: 'en' })}
                </p>
              </CardContent>
            </CardActionArea>
          </Card>
        </Grid>
      ) : null)}
    </Grid>)}
  </I18n>;
});

const renderNavigation = (props) => {
  const { routes, onItemClick, navigationClassName, ActionClassName, value } = props;
  return <I18n ns="translations">
    {t => (<BottomNavigation
      value={value}
      showLabels
      className={navigationClassName}
    >
      {routes.map(route => route ? (
        <BottomNavigationAction
          key={route.name}
          value={route.url}
          onClick={() => onItemClick(route)}
          label={t(route.title)}
          icon={<route.icon/>}
          className={ActionClassName}
        />
      ) : null)}
    </BottomNavigation>)}
  </I18n>;
};

const styleMap = {
  card: renderCard,
  navigation: renderNavigation
};

export default (props) => {
  const { type } = props;
  const Content = styleMap[type];
  return <Content {...props}/>;
}
