const gulp = require('gulp');
const debug = require('gulp-debug');

exports.Sync = (files, dest, task) => {
  console.log('start file sync', files);
  return gulp.src(files, { allowEmpty: true, since: gulp.lastRun(task) })
    .pipe(debug({ title: 'file sync' }))
    .pipe(gulp.dest(dest));
};
