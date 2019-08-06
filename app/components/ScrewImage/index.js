import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/styles';
import styles from './style';
import Point from './point';

export default function ScrewImage({ style, image, points, focus, activeIndex, onClick, pointScale }) {
  const classes = makeStyles(styles.image)();

  const [imageRef, setImageRef] = useState(null);
  const [containerRef, setContainerRef] = useState(null);

  const [size, setSize] = useState({
    width: 100,
    height: 100
  });

  const handleResize = () => {
    if (containerRef?.offsetHeight && containerRef?.offsetWidth) {
      setSize({
        height: (imageRef.offsetHeight / containerRef.offsetHeight) * 100 || 100,
        width: (imageRef.offsetWidth / containerRef.offsetWidth) * 100 || 100
      });
    }
  };

  const [focusStyle, setFocusStyle] = useState({
    transform: `translate(${0}%,${0}%) scale(${1},${1})`
  });

  useEffect(() => {
    if (focus && points?.[activeIndex]) {
      const transformX = (50 - points?.[activeIndex]?.x || 0) * focus;
      const transformY = (50 - points?.[activeIndex]?.y || 0) * focus;
      setFocusStyle({
        transform: `translate(${transformX || 0}%,${transformY || 0}%) scale(${focus},${focus})`
      });
    } else {
      setFocusStyle({
        transform: `translate(${0}%,${0}%) scale(${1},${1})`
      });
    }
  }, [activeIndex, focus, points]);


  useEffect(() => {
    window.addEventListener('resize', () => handleResize && handleResize());
    return () => window.removeEventListener('resize', () => handleResize && handleResize());
  }, [handleResize]);

  return (
    <div
      className={classes.container}
      ref={r => setContainerRef(r)}
      style={style}
    >
      <img
        ref={r => setImageRef(r)}
        src={image}
        className={classes.image}
        alt="job"
        onLoad={() => handleResize && handleResize()}
        style={{ ...focusStyle, transition: 'transform 1s' }}
      />
      <div
        style={{
          width: `${size.width || 100}%`,
          height: `${size.height || 100}%`,
          position: 'absolute',
          ...focusStyle,
          transition: 'transform 1s'
        }}
        onClick={() => onClick && onClick()}
      >
        {points?.map((p) => <Point
          key={`${p.x}${p.y}${p.status}`}
          x={p.x}
          y={p.y}
          status={p.status}
          label={p.group_sequence}
          scale={pointScale}
        />) || null}
      </div>
    </div>);
}

