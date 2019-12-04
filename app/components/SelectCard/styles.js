const styles = (height = 130) => (theme) => ({
  orderCard: {
    width: '100%',
    height,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'flex-start'
  },
  bgOdd: {
    overflow: 'auto',
    backgroundColor: '#747474'
  },
  bgEven: {
    overflow: 'auto',
    backgroundColor: '#7c7c7c'
  },
  image: {
    height,
    width: '30%'
  },
  info: {
    display: 'flex',
    flexDirection: 'column',
    padding: '10px 5px',
    height,
    flex: 1,
    textAlign: 'center',
    alignItems: 'flex-start',
    justifyContent: 'center',
    alignContent: 'center'
  },
  statusIndicator: {
    display: 'block',
    width: '5px',
    height
  }
});
export default styles;
