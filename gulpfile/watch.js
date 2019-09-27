exports.watchStyles = (files) => {
  const gulp = require('gulp');
  console.log('start watch gulpfile');
  const { styleTask } = require('./tasks');
  return gulp.watch(
    files,
    gulp.series(styleTask)
  );
};

exports.watchLive = () => {
  const livereload = require('gulp-livereload');
  console.log('start livereload');
  livereload.listen();
};
