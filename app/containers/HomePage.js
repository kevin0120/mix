// @flow
import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import Home from '../components/Home';

// function mapStateToProps(state) {
//   return {
//     counter: state.counter
//   };
// }
//
// function mapDispatchToProps(dispatch) {
//   return bindActionCreators(ScannerActions, dispatch);
// }

export default class HomePage extends Component<Props> {
  props: Props;

  render() {
    return <Home />;
  }
}

// export default connect(mapStateToProps, mapDispatchToProps)(Home)
