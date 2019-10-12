// @flow

import ReactTable from 'react-table';
import React from 'react';
import PropTypes from 'prop-types';

type Props = {
  translate: (string)=>string
};
export default class CustomReactTable extends React.Component<Props> {

  render() {
    const { translate, ...restProps } = this.props;
    return (
      <ReactTable
        previousText={translate('Table.Previous')}
        nextText={translate('Table.Next')}
        loadingText={translate('Table.Loading')}
        noDataText={translate('Table.NoRowsFound')}
        pageText={translate('Table.Page')}
        ofText={translate('Table.of')}
        rowsText={translate('Table.rows')}
        showPageJump={false}
        {...restProps}
      />

    );
  }
}

CustomReactTable.propTypes = {
  translate: PropTypes.func.isRequired
};
