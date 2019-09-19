const gulp = require('gulp');
const livereload = require('gulp-livereload');

const { compileStyle } = require('./style');

exports.watchStyles = (styleFolders) => {
  console.log('start watch gulpfile');
  const files = styleFolders.map((f) => `${f}/**/*.scss`);
  const styleFunc = () => compileStyle(files);
  styleFunc.displayName = 'regenerate scss';
  return gulp.watch(
    files,
    gulp.series(styleFunc)
  );
};

exports.watchLive = () => {
  console.log('start livereload');
  livereload.listen();
};
