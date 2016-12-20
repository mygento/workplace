var gulp  = require('gulp');
var watch = require('gulp-watch');

var source = './src';
var destination = './application/public';

gulp.task('watch-folder', function() {
  gulp.src(source + '/**/*', {base: source})
    .pipe(watch(source, {base: source}))
    .pipe(gulp.dest(destination));
});

gulp.task('default', ['watch-folder']);
