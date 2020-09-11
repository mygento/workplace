const notProduction = process.env.NODE_ENV !== 'production';

const debug = require('debug')('workplace:style');

const gulp = require('gulp');
const gdebug = require('gulp-debug');
const noop = require('through2').obj;
const notify = require('gulp-notify');
const livereload = require('gulp-livereload');

exports.compileStyle = (files, task, themeFolders, cssOptions = {}) => {
  debug('compile Style', files);
  const path = require('path');
  const postcss = require('gulp-postcss');
  const autoprefixer = require('autoprefixer');
  const cssnano = require('cssnano');
  const postcss_plugins = [
    autoprefixer(),
    cssnano(cssOptions)
  ];
  const sourcemaps = require('gulp-sourcemaps');
  const sass = require('gulp-sass');
  const sassInheritance = require('gulp-sass-parent');

  const includeOpt = { includePaths:
    [
      require('node-normalize-scss').includePaths,
      require('sassime').includePaths,
    ]
  };

  return gulp.src(
    files, {
      allowEmpty: true,
      since: gulp.lastRun(task),
    })
    .pipe(gdebug({ title: 'style' }))
    .pipe(notProduction ? sourcemaps.init() : noop())
    .pipe(notProduction ? sassInheritance({ dir: themeFolders.map(t => `${t}/scss`) }) : noop())
    .pipe(gdebug({ title: 'style with parent' }))
    .pipe(sass(includeOpt).on('error', sass.logError))
    .pipe(postcss(postcss_plugins))
    .pipe(notProduction ? sourcemaps.write('.') : noop())
    //.pipe(gdebug({ title: 'write style' }))
    .pipe(gulp.dest(function(file) {
      return path.join(file.base, '../css');
    }))
    .pipe(notProduction ? livereload()  : noop())
    .pipe(notProduction ? notify({
      message: 'Styles complete',
      onLast: true,
    }) : noop());
};
