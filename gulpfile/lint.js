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

exports.fixStyle = (files) => {
  console.log('lint style autofix', files);
  const gulpStylelint = require('gulp-stylelint');
  return gulp.src(files, { allowEmpty: true, base: '.' })
    .pipe(debug({ title: 'lint Style' }))
    .pipe(gulpStylelint({
      // debug: true,
      failAfterError: true,
      fix: true,
      reporters: [
        { formatter: 'verbose', console: true },
        { formatter: 'string', console: true },
      ],
    }))
    .pipe(gulp.dest('.'));
};

exports.lintJs = (files, jsOptions = {}) => {
  console.log('lint js', files);
  const eslint = require('gulp-eslint');
  return gulp.src(files, { allowEmpty: true })
    .pipe(debug({ title: 'lint Js' }))
    .pipe(eslint({
      baseConfig: {
        extends: ['mygento'],
      },
      rules: {
        semi: ['error', 'always'],
      },
    }))
    .pipe(eslint.format())
    .pipe(eslint.failAfterError());
};

exports.fixJs = (files) => {
  console.log('lint js autofix', files);
  const eslint = require('gulp-eslint');
  return gulp.src(files, { allowEmpty: true, base: '.' })
    .pipe(debug({ title: 'lint Js' }))
    .pipe(eslint({
      fix: true,
      baseConfig: {
        extends: ['mygento'],
      },
      rules: {
        semi: ['error', 'always'],
      },
    }))
    .pipe(eslint.format())
    // .pipe(eslint.failAfterError())
    .pipe(gulp.dest('.'));
};
