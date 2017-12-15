var gulp  = require('gulp');
var watch = require('gulp-watch');

var source = './src';
var destination = './application';

gulp.task('watch-folder', function() {
  gulp.src([
     `${source}/**/*`,
     `!${source}/node_modules`,
     `!${source}/vendor`,
     `!${source}/setup`,
     `!${source}/dev`
    ], {base: source})
    .pipe(watch(source, {base: source}))
    .pipe(gulp.dest(destination));
});

gulp.task('default', ['watch-folder']);
