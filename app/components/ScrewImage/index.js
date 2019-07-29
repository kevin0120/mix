import React, { useState, useEffect } from 'react';
import { makeStyles } from '@material-ui/styles';
import styles from './style';
import Point from './point';


export default function ScrewImage({ image, points, focus, activePoint }) {
  const classes = makeStyles(styles.image)();

  const [imageRef, setImageRef] = useState(null);
  const [containerRef, setContainerRef] = useState(null);

  const [size, setSize] = useState({
    width: 100,
    height: 100
  });

  function handleResize() {
    if (containerRef?.offsetHeight && containerRef?.offsetWidth) {
      setSize({
        height: (imageRef.offsetHeight / containerRef.offsetHeight) * 100 || 100,
        width: (imageRef.offsetWidth / containerRef.offsetWidth) * 100 || 100
      });
    }
  }

  const [focusStyle, setFocusStyle] = useState({
    transform: `translate(${0}%,${0}%) scale(${1},${1})`
  });

  useEffect(() => {
    if (focus && points?.[activePoint]) {
      const transformX = (50 - points?.[activePoint]?.x || 0) * 2;
      const transformY = (50 - points?.[activePoint]?.y || 0) * 2;
      setFocusStyle({
        transform: `translate(${transformX || 0}%,${transformY || 0}%) scale(${focus},${focus})`
      });
    } else {
      setFocusStyle({
        transform: `translate(${0}%,${0}%) scale(${1},${1})`
      });
    }
  }, [activePoint, focus, points]);


  useEffect(() => {
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, [handleResize]);

  return (
    <div
      className={classes.container}
      ref={r => setContainerRef(r)}
    >
      <img
        ref={r => setImageRef(r)}
        src={image}
        className={classes.image}
        alt="job"
        onLoad={() => handleResize()}
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
      >
        {points?.map((p, key) => <Point
          key={`${p.x}${p.y}${p.status}`}
          x={p.x}
          y={p.y}
          status={p.status}
          label={key + 1}
        />) || null}
      </div>
    </div>);
}

