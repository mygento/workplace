const notProduction = process.env.NODE_ENV !== 'production';

const gulp = require('gulp');
const noop = require('through2').obj;

const styleLintOptions = require('./scssLint');

let notify;

if (notProduction) {
  notify = require('gulp-notify');
}

exports.compileScss = function(
  sourcePath, targetPath, css_options = {}, task = 'scss'
) {
  const includeOpt = { includePaths: [
    require('node-normalize-scss').includePaths,
    require('sassime').includePaths,
  ],
  };
  const sourcemaps = require('gulp-sourcemaps');
  const sass = require('gulp-sass');
  const cssnano = require('gulp-cssnano');
  return gulp.src([`${sourcePath}/**/*.scss`], { since: gulp.lastRun(task) })
    .pipe(notProduction ? sourcemaps.init() : noop())
    // .pipe(notProduction ? sassInheritance({ dir: scss_folder }) : noop())
    .pipe(sass(includeOpt).on('error', sass.logError))
    .pipe(cssnano(css_options))
    .pipe(notProduction ? sourcemaps.write('.') : noop())
    .pipe(gulp.dest(targetPath))
    .pipe(notProduction ? notify({
      message: 'Styles complete',
      onLast: true,
    }) : noop());
};

exports.lintJs = function(codeFolder, themeFolder) {
  const eslint = require('gulp-eslint');
  return gulp.src([
    `${codeFolder}/**/*.js`,
    `${themeFolder}/**/*.js`,
    '!node_modules/**',
    `!${codeFolder}/**/web/js/vendor/**/*.js`,
    `!${codeFolder}/web/js/vendor/**/*.js`,
    `!${themeFolder}/web/mage/**/*.js`,
  ]).pipe(eslint())
    .pipe(eslint.format())
    .pipe(eslint.failAfterError());
};

exports.lintScss = function(themeFolder) {
  const gulpStylelint = require('gulp-stylelint');
  return gulp.src(
    [`${themeFolder}/**/*.scss`, `!${themeFolder}/vendor/**/*.scss`]
  ).pipe(gulpStylelint({
    config: styleLintOptions,
    reporters: [
      { formatter: 'string', console: true },
    ],
  }));
  // .pipe(scsslint({
  //   maxBuffer: 307200,
  // }))
  // .pipe(scsslint.failReporter('E'));
};
