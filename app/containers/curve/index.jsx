import React from 'react';
import configs from '../../shared/config';

export default class Curve extends React.Component {
  render() {
    return (
      <div
        style={{
          display: 'flex',
          height: 'calc(100% - 64px)'
        }}
      >
        <iframe
          title="cvinetweb"
          src={configs.cvinetweb.url}
          frameBorder={0}
          style={{
            display: 'flex',
            flex: 1
          }}
          seamless
        />
      </div>
    );
  }
}
