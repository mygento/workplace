var gulp  = require('gulp');
var watch = require('gulp-watch');
var debug = require('gulp-debug');

var source = './src';
var destination = './application';
var paths = [
  `${source}/**/*`,
  `!${source}/.git/**/*`,
  `!${source}/node_modules/**/*`,
  `!${source}/vendor/**/*`,
  `!${source}/setup/**/*`,
  `!${source}/dev/**/*`,
  `!${source}/composer.lock`,
];

gulp.task('watch-folder', function() {
  gulp.src(paths, { base: source})
    .pipe(debug())
    .pipe(watch(paths, { base: source, verbose: true }))
    .pipe(gulp.dest(destination));
});

gulp.task('default', ['watch-folder']);
