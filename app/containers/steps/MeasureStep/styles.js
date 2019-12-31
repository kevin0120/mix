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
    width: 'fit-content'
  },
  cardIconTitle: {
    ...theme.title.card,
    marginTop: '15px',
    marginBottom: '0px'
  },
  cardContent: {
    width: 'fit-content',
    padding: 60
  },
  rowContent: {
    flex: 1,
    display: 'flex',
    flexDirection: 'row',
    justifyContent: 'center',
    alignItems: 'center',
    textAlign: 'center',
    minWidth: '180px',
    height: '100%',
    paddingTop: 5,
    paddingBottom: 5,
  },
  row: {
    height: 60,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center'
  },
  inputContainer: {
    width: '100%',
    height: '100%',
    border: 'solid 1px',
    borderRadius: 3,
    borderColor: theme.palette.primary.main,
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent:'center'
  },
  inputContainerDisabled: {
    border: 'dashed 1px',
    borderColor: theme.palette.gray[500],
  },
  inputText: {
    color: theme.palette.gray[900]
  },
  inputTextDisabled: {
    color: theme.palette.gray[500]
  },
  rowText: {
    color: theme.palette.gray[900]
  }
});

export default styles;
