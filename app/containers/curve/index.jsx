import React from 'react';
import configs from '../../shared/config';

export default class Curve extends React.Component {
  constructor(props) {
    super(props);
    this.webview = null;
  }

  componentDidMount(): void {
    this.webview.addEventListener('new-window', e => {
      e.preventDefault();
      this.webview.setAttribute('src', e.url);
    });
  }

  render() {
    return (
      <div
        style={{
          display: 'flex',
          height: 'calc(100% - 64px)'
        }}
      >
        <webview
          ref={r => {
            this.webview = r;
          }}
          id="cvinetweb"
          title="cvinetweb"
          src={configs.cvinetweb.url}
          style={{
            height: '100%',
            width: '100%'
          }}
          // frameBorder="0"
          // sandbox="allow-forms allow-pointer-lock allow-same-origin allow-scripts allow-top-navigation allow-presentation"
        />
      </div>
    );
  }
}
