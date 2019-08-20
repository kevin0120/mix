import React from 'react';
import { GridLoader } from 'react-spinners';
// import {connect} from 're'
function Loading({ loading }) {
  return <GridLoader
    sizeUnit="px"
    size={50}
    color="#36D7B7"
    loading={loading}
  />;
}

const mapDispatch={

};

const mapState=(state,props)=>({

});

export default connect()(Loading)
