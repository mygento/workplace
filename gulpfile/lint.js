const gulp = require('gulp');
const debug = require('debug')('workplace:watch');
const gdebug = require('gulp-debug');

const debugMode = !!process.env.DEBUG;
const noop = require('through2').obj;

exports.lintStyle = (files, fail = false) => {
  debug('lint style', files);
  const gulpStylelint = require('gulp-stylelint');
  return gulp.src(files, { allowEmpty: true })
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
};

exports.fixStyle = (files) => {
  debug('lint style autofix', files);
  const gulpStylelint = require('gulp-stylelint');
  return gulp.src(files, { allowEmpty: true, base: '.' })
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
    .pipe(gulp.dest('.'));
};

exports.lintJs = (files, jsOptions = {}) => {
  debug('lint js', files);
  const eslint = require('gulp-eslint');
  return gulp.src(files, { allowEmpty: true })
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
};

exports.fixJs = (files) => {
  debug('lint js autofix', files);
  const eslint = require('gulp-eslint');
  return gulp.src(files, { allowEmpty: true, base: '.' })
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
    .pipe(gulp.dest('.'));
};
