const path = require('path');
const fs = require('fs');
const appDirectory = fs.realpathSync(process.env.PWD);
const resolveApp = relativePath => path.resolve(appDirectory, relativePath);
const config = require(resolveApp('package.json'));

const { mergeConfig } = require('./config');
const { watchStyles,
  watchLintStyles, watchLintJs,
  watchLive,
} = require('./watch');
const { lintStyle, lintJs, fixStyle, fixJs } = require('./lint');
const { compileStyle } = require('./style');
const { composeCommand } = require('./docker');
const { composerCommand } = require('./composer');

// Global Config
const workplaceConfig = mergeConfig(config, appDirectory);

const styleGlobs = (type) => {
  switch (type) {
  case 'magento2': {
    return ((workplaceConfig.magento2 || {}).theme || [])
      .map(f => resolveApp(`${f}/scss`));
  }
  }
  return [];
};

const jsFolders = (type) => {
  switch (type) {
  case 'magento2': {
    const theme = ((workplaceConfig.magento2 || {}).theme || []);
    return [
      resolveApp('app/code'),
      `!${resolveApp('app/code')}/**/web/js/vendor`,
      theme.map(f => resolveApp(f)),
      theme.map(f => `!${resolveApp(`${f}/js/mage/`)}`),
      theme.map(f => `!${resolveApp(`${f}/js/vendor/`)}`),
    ];
  }
  }
  return [];
};

const themeFolders = (type) => {
  switch (type) {
  case 'magento2': {
    return ((workplaceConfig.magento2 || {}).theme || []);
  }
  }
  return [];
};

// Theme config
const styleGlob = styleGlobs(workplaceConfig.type);
const lintJsGlob = [...(jsFolders(workplaceConfig.type) || [])];
const lintStyleGlob = [...styleGlob, ...styleGlob.map(f => `!${f}/vendor`)];
const themeFolder = themeFolders(workplaceConfig.type) || [];

// DEBUG
console.log('real config', workplaceConfig);
console.log('style Glob', styleGlob);
console.log('lintJs Glob', lintJsGlob);
console.log('lintStyleGlob', lintStyleGlob);

// Tasks
const liveTask = () => watchLive();
liveTask.displayName = 'livereload';

const watchStylesTask = (cb) => {
  if (styleGlob.length === 0) {
    return cb();
  }
  return watchStyles(
    styleGlob.map((f) => `${f}/**/*.scss`)
  );
};
watchStylesTask.displayName = 'watch style';

const watchLintStylesTask = (cb) => {
  if (lintStyleGlob.length === 0) {
    return cb();
  }
  return watchLintStyles(
    lintStyleGlob.map((f) => `${f}/**/*.scss`)
  );
};
watchLintStylesTask.displayName = 'watch lint scss';

const watchLintJsTask = (cb) => {
  if (lintJsGlob.length === 0) {
    return cb();
  }
  return watchLintJs(
    lintJsGlob.map((f) => `${f}/**/*.js`)
  );
};
watchLintJsTask.displayName = 'watch lint js';

const styleTask = (cb) => {
  if (styleGlob.length === 0) {
    return cb();
  }
  return compileStyle(
    styleGlob.map((f) => `${f}/**/*.scss`),
    styleTask,
    themeFolder.map(f => resolveApp(f))
  );
};
styleTask.displayName = 'scss';

const lintJsTask = (cb) => {
  if (lintJsGlob.length === 0) {
    return cb();
  }
  return lintJs(lintJsGlob.map((f) => `${f}/**/*.js`));
};
lintJsTask.displayName = 'js lint';

const lintStyleTask = (cb) => {
  if (lintStyleGlob.length === 0) {
    return cb();
  }
  return lintStyle(
    lintStyleGlob.map((f) => `${f}/**/*.scss`)
  );
};
lintStyleTask.displayName = 'style lint';

const fixStyleTask = (cb) => {
  if (lintStyleGlob.length === 0) {
    return cb();
  }
  return fixStyle(
    lintStyleGlob.map((f) => `${f}/**/*.scss`)
  );
};
fixStyle.displayName = 'style lint autofix';

const fixJsTask = (cb) => {
  if (lintJsGlob.length === 0) {
    return cb();
  }
  return fixJs(
    lintJsGlob.map((f) => `${f}/**/*.js`)
  );
};
fixJsTask.displayName = 'js lint autofix';

// Docker
const startDockerTask = (cb) => composeCommand(
  cb, ['up', '-d'], workplaceConfig
);
startDockerTask.displayName = 'docker';
const stopDockerTask  = (cb) => composeCommand(cb, ['stop'], workplaceConfig);
stopDockerTask.displayName = 'docker';
const rmDockerTask  = (cb) => composeCommand(cb, ['rm'], workplaceConfig);
rmDockerTask.displayName = 'docker';

const composerInstallTask = (cb) => composerCommand(
  cb, 'install', workplaceConfig
);

exports.lintStyleTask = lintStyleTask;
exports.lintJsTask = lintJsTask;
exports.startDockerTask = startDockerTask;
exports.stopDockerTask = stopDockerTask;
exports.rmDockerTask = rmDockerTask;
exports.watchStylesTask = watchStylesTask;
exports.watchLintStylesTask = watchLintStylesTask;
exports.watchLintJsTask = watchLintJsTask;
exports.liveTask = liveTask;
exports.styleTask = styleTask;
exports.fixStyleTask = fixStyleTask;
exports.fixJsTask = fixJsTask;
exports.composerTask = composerInstallTask;
