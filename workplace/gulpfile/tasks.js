import { resolve, join } from 'path';
import { realpathSync, readFileSync } from 'fs';
import { PROJECT_FILE, fileExists } from './gitignore.js';
import { mergeConfig } from './config.js';
import { watchLintJs } from './watch.js';
import { lintJs, fixJs } from './lint.js';
import { composeCommand } from './docker.js';
import { composerCommand } from './composer.js';
import createDebug from 'debug';

const debug = createDebug('workplace:tasks');
const appDirectory = realpathSync(process.env.PWD);
const resolveApp = relativePath => resolve(appDirectory, relativePath);
const jsonFile = (file) => JSON.parse(
  readFileSync(
    new URL(file, import.meta.url)
  )
);

const config = jsonFile(resolveApp('package.json'));
const overrideFile = resolveApp(join(PROJECT_FILE,'config.local.json'));
const override = fileExists(overrideFile) ? jsonFile(overrideFile) : {};

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

  process.chdir(appDirectory);

  return lintJs(lintJsGlob.map((f) => `${f}/**/*.js`));
};
lintJsTask.displayName = 'js lint';

const fixJsTask = (cb) => {
  if (lintJsGlob.length === 0) {
    return cb();
  }

  process.chdir(appDirectory);

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
composerInstallTask.displayName = 'composer';

export {
  lintJsTask,
  startDockerTask,
  stopDockerTask,
  rmDockerTask,
  watchLintJsTask,
  fixJsTask,
  composerInstallTask
 };
