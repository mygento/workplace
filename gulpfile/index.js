const path = require('path');
const gulp = require('gulp');
const fs = require('fs');
const appDirectory = fs.realpathSync(process.env.PWD);
const resolveApp = relativePath => path.resolve(appDirectory, relativePath);
const config = require(resolveApp('package.json'));

const workplaceConfig = Object.assign({}, {
  projectName: 'workplace',
  php: 'mygento/php:7.2-full',
  nginx: Object.assign({
    image: 'luckyraul/nginx:backports',
    port: 8081,
  }, config.nginx || {}),
  mysql: Object.assign({
    image: 'mygento/mysql',
    port: 3306,
  }, config.mysql || {}),
});

workplaceConfig.appDirectory = appDirectory;
if (config.name) {
  workplaceConfig.projectName = config.name;
}

console.log('gulp');
console.log('package.json', config.workplace);
console.log('real', workplaceConfig);

const { watchStyles, watchLive } = require('./watch');
const { composeStart } = require('./docker');
const style = require('./style');

const watchTask = () => watchStyles(
  config.workplace.theme.map(f => resolveApp(f))
);
watchTask.displayName = 'watch';

const liveTask = () => watchLive();
liveTask.displayName = 'livereload';

const dockerTask = (cb) => composeStart(cb, workplaceConfig);
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
