const gulp = require('gulp');
const debug = require('gulp-debug');

exports.lintStyle = (files, cssOptions = {}) => {
  console.log('lint style', files);
  const gulpStylelint = require('gulp-stylelint');
  return gulp.src(files, { allowEmpty: true })
    .pipe(debug())
    .pipe(gulpStylelint({
      reporters: [
        { formatter: 'string', console: true },
      ],
    }));
};

exports.lintJs = (files, jsOptions = {}) => {
  console.log('lint js', files);
  const eslint = require('gulp-eslint');
  return gulp.src(files, { allowEmpty: true })
    .pipe(debug())
    .pipe(eslint())
    .pipe(eslint.format())
    .pipe(eslint.failAfterError());
};
