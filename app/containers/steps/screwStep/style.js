import narBarStyles from '../../../components/NavBar/styles'

const screwImageStyles = (theme) => ({
  ...narBarStyles.switchWorkCenterButton(theme),
  layout: {
    width: '100%',
    height: '100%',
    overflow: 'hidden',
    position: 'relative'
  },
  thumbPaper: {
    position: 'absolute',
    bottom: 10,
    right: 10,
    width: 'auto'
  },
  thumbImage: {
    width: '200px',
    height: '200px'
  }
});

export default screwImageStyles;
