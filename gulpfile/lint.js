const gulp = require('gulp');
const debug = require('gulp-debug');

exports.lintStyle = (files, fail = false) => {
  console.log('lint style', files);
  const gulpStylelint = require('gulp-stylelint');
  return gulp.src(files, { allowEmpty: true })
    .pipe(debug({ title: 'lint Style' }))
    .pipe(gulpStylelint({
      // debug: true,
      failAfterError: fail,
      fix: false,
      reporters: [
        { formatter: 'verbose', console: true },
        { formatter: 'string', console: true },
      ],
    }));
};

exports.lintJs = (files, jsOptions = {}) => {
  console.log('lint js', files);
  const eslint = require('gulp-eslint');
  return gulp.src(files, { allowEmpty: true })
    .pipe(debug({ title: 'lint Js' }))
    .pipe(eslint())
    .pipe(eslint.format())
    .pipe(eslint.failAfterError());
};
