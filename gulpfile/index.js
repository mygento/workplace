const gulp = require('gulp');

const {
  watchStylesTask, watchLintStylesTask, watchLintJsTask,
  styleTask,
  startDockerTask, liveTask, stopDockerTask, rmDockerTask,
  lintJsTask, lintStyleTask,
  fixStyleTask, fixJsTask,
  composerTask,
  watchSyncTask, syncTask
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
  lint: gulp.parallel(
    lintJsTask,
    lintStyleTask
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
    lintStyleTask,
    syncTask,
    watchSyncTask
  ),
}, style);

module.exports = result;
