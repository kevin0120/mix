// @flow
import React, { useState, useEffect, useRef, useCallback } from 'react';
import { makeStyles } from '@material-ui/styles';
import styles from './style';
import Point from './point';
import type { tAction } from '../../modules/typeDef';

import { ClsOperationPoint } from '../../modules/step/screwStep/classes/ClsOperationPoint';
import type { CommonLogLvl } from '../../common/utils';

type Props = {
  twinkle: boolean,
  style?: {},
  image: string,
  points: Array<ClsOperationPoint>,
  focus: number,
  enableReWork?: boolean,
  notifyInfo?: ?(variant: CommonLogLvl, message: string, meta: Object)=>tAction<any, any>,
  // activeIndex: number,
  pointScale?: number,
  onPointClick?: ?(ClsOperationPoint)=>boolean
};

export default function ScrewImage({ twinkle, style = {}, image, points, focus, pointScale = 1, notifyInfo, enableReWork, onPointClick }: Props) {
  const classes = makeStyles(styles.image)();
  
  const imageRef = useRef(null);
  const containerRef = useRef(null);
  const [size, setSize] = useState({
    width: 100,
    height: 100
  });
  
  const handleResize = useCallback(() => {
    const img = imageRef.current;
    const container = containerRef.current;
    if (container && img && container?.offsetHeight && container?.offsetWidth) {
      setSize({
        height: (img.offsetHeight / container.offsetHeight) * 100 || 100,
        width: (img.offsetWidth / container.offsetWidth) * 100 || 100
      });
    }
  }, []);
  
  const [focusStyle, setFocusStyle] = useState({
    transform: `translate(${0}%,${0}%) scale(${1},${1})`
  });
  
  useEffect(() => {
    // if (focus && points?.[activeIndex]) {
    //   const transformX = (50 - points?.[activeIndex]?.x || 0) * focus;
    //   const transformY = (50 - points?.[activeIndex]?.y || 0) * focus;
    //   setFocusStyle({
    //     transform: `translate(${transformX || 0}%,${transformY || 0}%) scale(${focus},${focus})`
    //   });
    // } else {
    //   setFocusStyle({
    //     transform: `translate(${0}%,${0}%) scale(${1},${1})`
    //   });
    // }
  }, [focus, points]);
  
  
  useEffect(() => {
    window.addEventListener('resize', () => handleResize && handleResize());
    return () => window.removeEventListener('resize', () => handleResize && handleResize());
  }, [handleResize]);
  
  return (
    <div
      className={classes.container}
      ref={containerRef}
      style={style}
    >
      <img
        ref={imageRef}
        src={image || ''}
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
          ...(focusStyle || {}),
          transition: 'transform 1s'
        }}
      >
        {points && points.map((p) => <Point
          key={`${p.sequence}`}
          x={p.x}
          y={p.y}
          twinkle={twinkle && p.isActive}
          status={p.status}
          label={`${p?.sequence || ''}`}
          scale={pointScale}
          reworkModiBGColor
          onClick={(e) => {
            e.preventDefault();
            e.stopPropagation();
            if (!enableReWork && notifyInfo) {
              notifyInfo('Warn', '当前工作模式无法进行返工作业，请先切换至返工模式!');
              return false;
            }
            if (onPointClick) {
              return onPointClick(p);
            }
            return false;
          }}
        />) || null}
      </div>
    </div>);
}

ScrewImage.defaultProps = {
  style: {},
  pointScale: 1,
  enableReWork: false,
  notifyInfo: null,
  onPointClick: null
};

