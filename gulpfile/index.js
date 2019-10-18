const gulp = require('gulp');

const {
  watchStylesTask, watchLintStylesTask, watchLintJsTask,
  styleTask,
  dockerTask, liveTask,
  lintJsTask, lintStyleTask,
  fixStyleTask, fixJsTask,
} = require('./tasks');

const style = require('./style');

const result = Object.assign({
  someMagic: function(cb) {
    cb();
  },
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
    // dockerTask,
    lintJsTask,
    lintStyleTask
  ),
}, style);

module.exports = result;
