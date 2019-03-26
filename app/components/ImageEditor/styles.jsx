import CONFIG from '../../shared/config';

const circleWidth = CONFIG.systemSettings.psetPointDiameter;

export default () => ({
  imageWrap: {
    height: '100%',
    width:'100%',
    marginTop: 0,
    position: 'relative',
    overflow: 'hidden',
    textAlign: 'center',
    backgroundSize: 'contain',
    backgroundRepeat: 'no-repeat',
    backgroundPosition: 'center',
  },
  img: {
    width: '100%',
  },
  circleText:{

  },
  circleWrap: {
  },
});
