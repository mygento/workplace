exports.watchStyles = (files) => {
  const gulp = require('gulp');
  console.log('start watch styles');
  const { styleTask } = require('./tasks');
  return gulp.watch(
    files,
    gulp.series(styleTask)
  );
};

exports.watchLintStyles = (files) => {
  const gulp = require('gulp');
  console.log('start watch lint styles');
  const { lintStyleTask } = require('./tasks');
  return gulp.watch(
    files,
    gulp.series(lintStyleTask)
  );
};

exports.watchLintJs = (files) => {
  const gulp = require('gulp');
  console.log('start watch lint js', files);
  const { lintJsTask } = require('./tasks');
  return gulp.watch(
    files,
    gulp.series(lintJsTask)
  );
};

exports.watchLive = () => {
  const livereload = require('gulp-livereload');
  console.log('start livereload');
  livereload.listen();
};
