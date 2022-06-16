const debug = require('debug')('workplace:watch');

exports.watchStyles = (files) => {
  const gulp = require('gulp');
  debug('start watch styles');
  const { styleTask } = require('./tasks');
  return gulp.watch(
    files,
    gulp.series(styleTask)
  );
};

exports.watchLintStyles = (files) => {
  const gulp = require('gulp');
  debug('start watch lint styles');
  const { lintStyleTask } = require('./tasks');
  return gulp.watch(
    files,
    gulp.series(lintStyleTask)
  );
};

exports.watchLintJs = (files) => {
  const gulp = require('gulp');
  debug('start watch lint js', files);
  const { lintJsTask } = require('./tasks');
  return gulp.watch(
    files,
    gulp.series(lintJsTask)
  );
};
