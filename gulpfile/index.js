const gulp = require('gulp');

const {
  watchLintStylesTask, watchLintJsTask,
  startDockerTask, stopDockerTask, rmDockerTask,
  lintJsTask, lintStyleTask,
  fixStyleTask, fixJsTask,
  composerTask
} = require('./tasks');

const result = {
  install: gulp.series(
    composerTask
  ),
  stop: gulp.series(
    stopDockerTask
  ),
  delete: gulp.series(
    stopDockerTask,
    rmDockerTask
  ),
  lint: gulp.parallel(
    lintJsTask,
    lintStyleTask
  ),
  test: gulp.parallel(
    fixJsTask,
    fixStyleTask
  ),
  start: gulp.parallel(
    watchLintStylesTask,
    watchLintJsTask,
    startDockerTask,
    lintJsTask,
    lintStyleTask
  ),
};

module.exports = result;
