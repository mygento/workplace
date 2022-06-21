import pkg from 'gulp';
import { obj as noop } from 'through2';
import createDebug from 'debug';
import gdebug from 'gulp-debug';
import gulpStylelint from 'gulp-stylelint';

const debug = createDebug('workplace:lint');
const debugMode = !!process.env.DEBUG;
const { src, dest } = pkg;

export function lintStyle(files, fail = false) {
  debug('lint style', files);

  return src(files, { allowEmpty: true })
    .pipe(debugMode ? gdebug({ title: 'lint Style' }) : noop())
    .pipe(gulpStylelint({
      // debug: true,
      failAfterError: fail,
      fix: false,
      reporters: [
        { formatter: 'verbose', console: true },
        { formatter: 'string', console: true },
      ],
    }));
}

export function fixStyle(files) {
  debug('lint style autofix', files);
  return src(files, { allowEmpty: true, base: '.' })
    .pipe(debugMode ? gdebug({ title: 'lint Style' }) : noop())
    .pipe(gulpStylelint({
      // debug: true,
      failAfterError: true,
      fix: true,
      reporters: [
        { formatter: 'verbose', console: true },
        { formatter: 'string', console: true },
      ],
    }))
    .pipe(dest('.'));
}
