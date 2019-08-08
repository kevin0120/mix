import React from 'react';
import {connect} from 'react-redux';

function materialStep(){


  return(
    <div>

    </div>
  )
}

const mapState = (state, props) => ({
  ...props,

});

const mapDispatch = {

};

export default connect(mapState, mapDispatch)(materialStep);
