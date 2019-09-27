const gulp = require('gulp');

const {
  watchTask, liveTask, dockerTask, lintJsTask, lintStyleTask, styleTask,
} = require('./tasks');

const style = require('./style');

const result = Object.assign({
  someMagic: function(cb) {
    cb();
  },
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
