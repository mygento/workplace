var gulp  = require('gulp');
var watch = require('gulp-watch');
var debug = require('gulp-debug');

var source = './src';
var destination = './application';

gulp.task('watch-folder', function() {
  gulp.src([
     `${source}/**/*`,
     `!${source}/node_modules`,
     `!${source}/vendor`,
     `!${source}/setup`,
     `!${source}/dev`,
     `!${source}/composer.lock`,
    ], {base: source})
    .pipe(debug())
    .pipe(watch(source, {base: source}))
    .pipe(gulp.dest(destination));
});

gulp.task('default', ['watch-folder']);
