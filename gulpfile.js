const gulp  = require('gulp');
const watch = require('gulp-watch');
const debug = require('gulp-debug');

const source = './src';
const destination = './application';
const paths = [
  `${source}/**/*`,
  `!${source}/.*/**/*`,
  `${source}/.php_cs`,
  `!${source}/node_modules/**/*`,
  `!${source}/dev/**/*`,
  `!${source}/generated/**/*`,
  `!${source}/lib/**/*`,
  `!${source}/setup/**/*`,
  `!${source}/phpserver/**/*`,
  `!${source}/vendor/**/*`,
  `!${source}/composer.lock`,
  `!${source}/app/etc/di.xml`,
];

gulp.task('watch-folder', () => gulp.src(paths, { base: source, allowEmpty: true })
    .pipe(debug())
    .pipe(watch(paths, { base: source, verbose: true, allowEmpty: true }))
    .pipe(gulp.dest(destination))
);

gulp.task('default', gulp.series('watch-folder'));
