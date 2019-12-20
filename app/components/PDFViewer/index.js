import { Document, Page } from 'react-pdf/dist/entry.webpack';
import React, { useState } from 'react';
import makeStyles from '@material-ui/core/styles/makeStyles';
import styles from './styles';

export default function PDFViewer({ file = '', page, scale = 1.3, bindNext, bindPrevious }) {
  const classes = makeStyles(styles)();
  const [pageNum, setPageNum] = useState(0);
  // TODO switch pages
  return <div className={classes.container}>
    <Document
      className={classes.document}
      file={{
        url: file
      }}
      onLoadSuccess={() => setPageNum(page)}
      onLoadError={e => console.log(e)}
    >
      <Page
        className={classes.page}
        scale={scale}
        pageIndex={pageNum > page ? page : pageNum}
        renderAnnotationLayer={false}
      />
    </Document>
  </div>;
}

