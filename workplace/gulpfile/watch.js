import createDebug from 'debug';
import pkg from 'gulp';

import {  
  lintJsTask
} from './tasks.js';

const debug = createDebug('workplace:watch');
const { watch, series } = pkg;

export function watchLintJs(files) { 
  debug('start watch lint js', files);  
  return watch(
    files,
    series(lintJsTask)
  );
}
