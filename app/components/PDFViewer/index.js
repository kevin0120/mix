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
      file={file}
      onLoadSuccess={(pdf) => setPageNum(pdf.numPages)}
      onLoadError={e => console.log(e)}
      loading="正在加载..."
      noData="没有找到文件"
      error="打开文件失败"
    >
      {[...new Array(pageNum)].map((p, idx) => idx).map(p => <Page
        className={classes.page}
        scale={scale}
        loading="正在加载..."
        noData="页面为空"
        error="打开页面失败"
        pageIndex={pageNum > p ? p : pageNum}
        renderAnnotationLayer={false}
      />)}
    </Document>
  </div>;
}

