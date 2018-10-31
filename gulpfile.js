const gulp  = require('gulp');
const watch = require('gulp-watch');
const debug = require('gulp-debug');

const source = './src';
const destination = './application';
const paths = [
  `${source}/**/*`,
  `!${source}/.*/**/*`,
  `!${source}/node_modules/**/*`,
  `!${source}/vendor/**/*`,
  `!${source}/setup/**/*`,
  `!${source}/dev/**/*`,
  `!${source}/composer.lock`,
  `!${source}/app/etc/di.xml`,
];

gulp.task('watch-folder', () => gulp.src(paths, { base: source})
    .pipe(debug())
    .pipe(watch(paths, { base: source, verbose: true }))
    .pipe(gulp.dest(destination))
);

gulp.task('default', gulp.series('watch-folder'));
