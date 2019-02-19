import React from 'react';
import SweetAlert from 'react-bootstrap-sweetalert';

export default class Alert extends React.Component {
  render() {
    const { show, style, children, ...restProps } = this.props;
    return (
      show ? (<div
          style={{
            position: 'fixed',
            height: '100vh',
            width: '100vw',
            top: 0,
            left: 0,
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            alignContents: 'center',
            zIndex: 99999
          }}
        >
          <SweetAlert
            show={show}
            style={{
              display: 'block',
              position: 'relative',
              marginTop: 0,
              ...style
            }}
            {...restProps}
          >
            {children}
          </SweetAlert>
        </div>)
        : null
    );
  }
}
