const path = require('path');
const fs = require('fs');
const debug = require('debug')('workplace:tasks');

const { PROJECT_FILE, fileExists } = require('./gitignore');
const { mergeConfig } = require('./config');

const appDirectory = fs.realpathSync(process.env.PWD);
const resolveApp = relativePath => path.resolve(appDirectory, relativePath);
const config = require(resolveApp('package.json'));
const overrideFile = resolveApp(path.join(PROJECT_FILE,'config.local.json'));
const override = fileExists(overrideFile) ? require(overrideFile) : {};

const { watchLintStyles, watchLintJs } = require('./watch');
const { lintStyle, lintJs, fixStyle, fixJs } = require('./lint');
const { composeCommand } = require('./docker');
const { composerCommand } = require('./composer');

// Global Config
const workplaceConfig = mergeConfig(config, appDirectory, override);

const styleGlobs = (type) => {
  switch (type) {
    case 'magento2': {
      if (!(workplaceConfig.magento2 || {}).style) {
        return [];
      }
      return ((workplaceConfig.magento2 || {}).theme || [])
        .map(f => resolveApp(`${f}/web/scss`));
    }
  }
  return [];
};

const jsFolders = (type) => {
  switch (type) {
    case 'magento2': {
      const theme = ((workplaceConfig.magento2 || {}).theme || []);
      const lint = ((workplaceConfig.magento2 || {}).lint || []);
      return [
        resolveApp('app/code'),
        `!${resolveApp('app/code')}/**/web/js/vendor`,
        ...theme.map(f => resolveApp(f)),
        ...theme.map(f => `!${resolveApp(`${f}/web/js/mage/`)}`),
        ...theme.map(f => `!${resolveApp(`${f}/web/js/vendor/`)}`),
        ...lint.map(f => f.indexOf('!') !== 0 ? resolveApp(f) : `!${resolveApp(f.substring(1))}`),
      ];
    }
  }
  return [];
};

// Theme config
const styleGlob = styleGlobs(workplaceConfig.type);
const lintJsGlob = [...(jsFolders(workplaceConfig.type) || [])];
const lintStyleGlob = [...styleGlob, ...styleGlob.map(f => `!${f}/vendor`)];

// DEBUG
debug('real config', workplaceConfig);
debug('style Glob', styleGlob);
debug('lintJs Glob', lintJsGlob);
debug('lintStyle Glob', lintStyleGlob);

// Tasks
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
exports.watchLintStylesTask = watchLintStylesTask;
exports.watchLintJsTask = watchLintJsTask;
exports.fixStyleTask = fixStyleTask;
exports.fixJsTask = fixJsTask;
exports.composerTask = composerInstallTask;
