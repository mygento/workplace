const path = require('path');
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

exports.resolveApp = resolveApp;
exports.appConfig = config;
exports.workplaceConfig = workplaceConfig;
