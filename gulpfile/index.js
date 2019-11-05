const gulp = require('gulp');

const {
  watchStylesTask, watchLintStylesTask, watchLintJsTask,
  styleTask,
  startDockerTask, liveTask, stopDockerTask, rmDockerTask,
  lintJsTask, lintStyleTask,
  fixStyleTask, fixJsTask,
} = require('./tasks');

const style = require('./style');

const result = Object.assign({
  someMagic: function(cb) {
    cb();
  },
  stop: gulp.parallel(
    stopDockerTask
  ),
  delete: gulp.parallel(
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
