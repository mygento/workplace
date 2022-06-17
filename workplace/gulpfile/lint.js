import pkg from 'gulp';
const { src, dest } = pkg;
import gdebug from 'gulp-debug';
import { obj as noop } from 'through2';

import eslint from 'gulp-eslint';

import createDebug from 'debug';
const debug = createDebug('workplace:watch');
const debugMode = !!process.env.DEBUG;

export function lintJs(files, jsOptions = {}) {
  debug('lint js', files);

  return src(files, { allowEmpty: true })
    .pipe(debugMode ? gdebug({ title: 'lint Js' }) : noop())
    .pipe(eslint({
      baseConfig: {
        extends: ['mygento'],
      },
      rules: {
        semi: ['error', 'always'],
      },
    }))
    .pipe(eslint.format())
    .pipe(eslint.failAfterError());
}

export function fixJs(files) {
  debug('lint js autofix', files);

  return src(files, { allowEmpty: true, base: '.' })
    .pipe(debugMode ? gdebug({ title: 'lint Js' }) : noop())
    .pipe(eslint({
      fix: true,
      baseConfig: {
        extends: ['mygento'],
      },
      rules: {
        semi: ['error', 'always'],
      },
    }))
    .pipe(eslint.format())
    // .pipe(eslint.failAfterError())
    .pipe(dest('.'));
}
