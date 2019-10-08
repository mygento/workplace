const gulp = require('gulp');

const {
  watchTask, liveTask, styleTask,
  dockerTask,
  lintJsTask, lintStyleTask,
  fixStyleTask,
} = require('./tasks');

const style = require('./style');

const result = Object.assign({
  someMagic: function(cb) {
    cb();
  },
  test: gulp.parallel(
    // lintJsTask,
    fixStyleTask
  ),
  start: gulp.parallel(
    styleTask,
    watchTask,
    liveTask,
    // dockerTask,
    // lintJsTask,
    lintStyleTask
  ),
}, style);

module.exports = result;
