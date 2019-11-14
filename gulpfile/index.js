const gulp = require('gulp');

const {
  watchStylesTask, watchLintStylesTask, watchLintJsTask,
  styleTask,
  startDockerTask, liveTask, stopDockerTask, rmDockerTask,
  lintJsTask, lintStyleTask,
  fixStyleTask, fixJsTask,
  composerTask,
} = require('./tasks');

const style = require('./style');

const result = Object.assign({
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
  test: gulp.parallel(
    fixJsTask,
    fixStyleTask
  ),
  start: gulp.parallel(
    styleTask,
    watchStylesTask,
    watchLintStylesTask,
    watchLintJsTask,
    liveTask,
    startDockerTask,
    lintJsTask,
    lintStyleTask
  ),
}, style);

module.exports = result;
