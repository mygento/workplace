var gulp  = require('gulp');
var watch = require('gulp-watch');
var debug = require('gulp-debug');

var source = './application/src';
var destination = './application/public';

gulp.task('watch-folder', function() {
  gulp.src(source + '/**/*', {base: source})
    .pipe(debug())
    .pipe(watch(source, {base: source, verbose: true}))
    .pipe(gulp.dest(destination));
});

gulp.task('default', ['watch-folder']);
