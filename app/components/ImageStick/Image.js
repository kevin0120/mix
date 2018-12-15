import React from 'react';
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
    position: 'relative',
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
    this.updateImgSize();
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
    this.imageSize = {
      height: (this.imageRef.offsetHeight || 0) / this.containerRef.offsetHeight * 100,
      width: (this.imageRef.offsetWidth || 0) / this.containerRef.offsetWidth * 100
    };
  }



  render() {
    const { style, src, alt, children, classes} = this.props;
    this.updateImgSize();

    return (
      <div
        ref={r => {
          this.containerRef = r;
        }}
        className={classes.imgBlock}
        style={style}
      >
        <img
          ref={r => {
            this.imageRef = r;
          }}
          src={src}
          className={classes.imgSheet}
          alt={alt}
        />
        <div style={{
          width: `${this.imageSize.width || 0}%`,
          height: `${this.imageSize.height || 0}%`,
          position: 'absolute'
        }}>
          {children}
        </div>
      </div>
    );
  }
}
export default withStyles(styles)(Image);
