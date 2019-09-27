const notProduction = process.env.NODE_ENV !== 'production';

const gulp = require('gulp');
const debug = require('gulp-debug');
const noop = require('through2').obj;
const notify = require('gulp-notify');
const livereload = require('gulp-livereload');

exports.compileStyle = (files, cssOptions = {}) => {
  console.log('compile Style');
  const path = require('path');
  const cssnano = require('gulp-cssnano');
  const sourcemaps = require('gulp-sourcemaps');
  const sass = require('gulp-sass');

  const includeOpt = { includePaths: [
    require('node-normalize-scss').includePaths,
    require('sassime').includePaths,
  ],
  };

  return gulp.src(files, { allowEmpty: true })
    .pipe(debug())
    .pipe(notProduction ? sourcemaps.init() : noop())
    // .pipe(notProduction ? sassInheritance({ dir: scss_folder }) : noop())
    .pipe(sass(includeOpt).on('error', sass.logError))
    .pipe(cssnano(cssOptions))
    .pipe(notProduction ? sourcemaps.write('.') : noop())
    .pipe(gulp.dest(function(file) {
      return path.join(file.base, '../css');
    }))
    .pipe(notProduction ? livereload()  : noop())
    .pipe(notProduction ? notify({
      message: 'Styles complete',
      onLast: true,
    }) : noop());
};
