const styles = (theme) => ({
  root: {
    width: '100%',
    height: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center'
  },
  card: {
    width:'fit-content'
  },
  cardIconTitle: {
    ...theme.title.card,
    marginTop: '15px',
    marginBottom: '0px'
  },
  cardContent: {
    width:'fit-content',
    padding: 60
  },
  row: {
    // display: 'flex',
    // flexDirection: 'row',
    // justifyContent: 'center',
    // alignItems: 'center',
    // marginTop: 10,
    // marginBottom: 10
  },
  label: {
    // marginRight: 20
  },
});

export default styles;
