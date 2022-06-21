import pkg from 'gulp';
import createDebug from 'debug';
import { styleTask, lintStyleTask }  from './index.js';

const { watch, series } = pkg;
const debug = createDebug('workplace:watch');

export function watchStyles(files) {
  debug('start watch styles');
  return watch(
    files,
    series(styleTask)
  );
}

export function watchLintStyles(files) {
  debug('start watch lint styles');
  return watch(
    files,
    series(lintStyleTask)
  );
}
