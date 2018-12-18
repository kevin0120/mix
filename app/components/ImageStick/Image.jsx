import React from 'react';
import PropTypes from 'prop-types';

// import popoverStyles from '../../common/jss/popoverStyles';
// import { dangerColor, successColor, warningColor } from '../../common/jss/material-react-pro';
import { withStyles } from '@material-ui/core/styles';

const styles = () => ({
  imgBlock: {
    width: '100%',
    height: '100%',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
    position: 'relative'
  },
  imgSheet: {
    maxHeight: '100%',
    maxWidth: '100%',
    textAlign: 'center',
    objectFit: 'contain',
    resize: 'both'
  }
});

class Image extends React.Component {
  constructor(props) {
    super(props);
    this.imageSize = {
      height: 0,
      width: 0
    };
    this.imageRef = React.createRef();
    this.containerRef = React.createRef();
    this.updateImgSize = this.updateImgSize.bind(this);
    this.handleResize = this.handleResize.bind(this);

  }

  componentDidMount() {
    // this.updateImgSize();
    window.addEventListener('resize', this.handleResize);
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.handleResize);
  }

  handleResize() {
    this.updateImgSize();
    this.forceUpdate();
  }

  updateImgSize() {
    const { focus } = this.props;
    if (this.containerRef.offsetHeight !== 0 && this.containerRef.offsetWidth !== 0) {
      this.imageSize = {
        height: (this.imageRef.offsetHeight || 0) / this.containerRef.offsetHeight * 100,
        width: (this.imageRef.offsetWidth || 0) / this.containerRef.offsetWidth * 100
      };
    }
  }

  render() {
    const { src, alt, children, classes, focus } = this.props;
    let style={};
    if (focus && focus.do === true) {
      style = {
        transform: `translate(${focus.transform.x || 0}%,${focus.transform.y || 0
          }%) scale(${focus.scale},${focus.scale})`
      };
    }

    return (
      <div
        ref={r => {
          this.containerRef = r;
        }}
        className={classes.imgBlock}

      >
        <img
          ref={r => {
            this.imageRef = r;
          }}
          src={src}
          className={classes.imgSheet}
          alt={alt}
          onLoad={() => {
            this.handleResize();
          }}
          style={style}
        />
        <div
          style={{
            width: `${this.imageSize.width || 0}%`,
            height: `${this.imageSize.height || 0}%`,
            position: 'absolute',
            ...style
          }}>
          {children}
        </div>
      </div>
    );
  }
}

Image.propTypes = {
  classes: PropTypes.shape({}).isRequired,
  className: PropTypes.string

};

export default withStyles(styles)(Image);
