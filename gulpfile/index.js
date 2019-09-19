const path = require('path');
const gulp = require('gulp');
const fs = require('fs');
const appDirectory = fs.realpathSync(process.env.PWD);
const resolveApp = relativePath => path.resolve(appDirectory, relativePath);
const config = require(resolveApp('package.json'));

console.log('gulp');
console.log(config.workplace);

const { watchStyles, watchLive } = require('./watch');
const { compose } = require('./docker');
const style = require('./style');
const watchTask = () => watchStyles(
  config.workplace.theme.map(f => resolveApp(f))
);
watchTask.displayName = 'watch';
const liveTask = () => watchLive();
liveTask.displayName = 'livereload';
const dockerTask = (cb) => compose(cb);
dockerTask.displayName = 'docker';
const result = Object.assign({
  someMagic: function(cb) {
    cb();
  },
  start: gulp.parallel(
    watchTask,
    liveTask,
    dockerTask
  ),
}, style);

module.exports = result;
