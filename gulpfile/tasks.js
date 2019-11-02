const { workplaceConfig, resolveApp } = require('./config');

const { watchStyles,
  watchLintStyles, watchLintJs,
  watchLive,
} = require('./watch');
const { lintStyle, lintJs, fixStyle, fixJs } = require('./lint');
const { compileStyle } = require('./style');
const { composeStart, composeStop } = require('./docker');

const styleGlob = workplaceConfig.theme.map(f => resolveApp(`${f}/scss`));
console.log('styleGlob', styleGlob);

const jsFolders = {
  magento2: [
    resolveApp('app/code'),
    `!${resolveApp('app/code')}/**/web/js/vendor`,
    ...workplaceConfig.theme.map(f => resolveApp(f)),
    ...workplaceConfig.theme.map(f => `!${resolveApp(`${f}/js/mage/`)}`),
    ...workplaceConfig.theme.map(f => `!${resolveApp(`${f}/js/vendor/`)}`),
  ],
};

const lintJsGlob = [
  ...(jsFolders[workplaceConfig.type] || []),
];
console.log('lintJsGlob', lintJsGlob);

const lintStyleGlob = [...styleGlob, ...styleGlob.map(f => `!${f}/vendor`)];
console.log('lintStyleGlob', lintStyleGlob);

// Tasks
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

const liveTask = () => watchLive();
liveTask.displayName = 'livereload';

const styleTask = (cb) => {
  if (styleGlob.length === 0) {
    return cb();
  }
  return compileStyle(
    styleGlob.map((f) => `${f}/**/*.scss`),
    styleTask,
    workplaceConfig.theme.map(f => resolveApp(f))
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
const dockerTask = (cb) => composeStart(cb, workplaceConfig);
dockerTask.displayName = 'docker';
const stopDockerTask  = (cb) => composeStop(cb, workplaceConfig);
stopDockerTask.displayName = 'docker';

exports.lintStyleTask = lintStyleTask;
exports.lintJsTask = lintJsTask;
exports.dockerTask = dockerTask;
exports.stopDockerTask = stopDockerTask;
exports.watchStylesTask = watchStylesTask;
exports.watchLintStylesTask = watchLintStylesTask;
exports.watchLintJsTask = watchLintJsTask;
exports.liveTask = liveTask;
exports.styleTask = styleTask;
exports.fixStyleTask = fixStyleTask;
exports.fixJsTask = fixJsTask;
