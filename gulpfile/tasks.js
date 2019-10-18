const { appConfig, workplaceConfig, resolveApp } = require('./config');

const { watchStyles, watchLive } = require('./watch');
const { lintStyle, lintJs, fixStyle, fixJs } = require('./lint');
const { compileStyle } = require('./style');
const { composeStart } = require('./docker');

const styleGlob = appConfig.workplace.theme.map(f => resolveApp(`${f}/scss`));
console.log('styleGlob', styleGlob);

const lintJsGlob = [
  [
    resolveApp('app/code'),
    '!node_modules',
    `!${resolveApp('app/code')}/**/web/js/vendor/`,
  ],
  ...appConfig.workplace.theme.map(f => resolveApp(f)),
  ...appConfig.workplace.theme.map(f => `!${resolveApp(`${f}/web/mage/`)}`),
];
console.log('lintJsGlob', lintJsGlob);

const lintStyleGlob = [...styleGlob, ...styleGlob.map(f => `!${f}/vendor`)];
console.log('lintStyleGlob', lintStyleGlob);

// Tasks
const watchTask = () => watchStyles(
  styleGlob.map((f) => `${f}/**/*.scss`),
  workplaceConfig
);
watchTask.displayName = 'watch';

const liveTask = () => watchLive();
liveTask.displayName = 'livereload';

const dockerTask = (cb) => composeStart(cb, workplaceConfig);
dockerTask.displayName = 'docker';

const styleTask = () => compileStyle(styleGlob.map((f) => `${f}/**/*.scss`));
styleTask.displayName = 'scss';

console.log(lintJsGlob.map((f) => `${f}/**/*.js`));
const lintJsTask = () => lintJs(lintJsGlob.map((f) => `${f}/**/*.js`));
lintJsTask.displayName = 'js lint';

const lintStyleTask = () => lintStyle(
  lintStyleGlob.map((f) => `${f}/**/*.scss`)
);
lintStyleTask.displayName = 'style lint';

const fixStyleTask = () => fixStyle(
  lintStyleGlob.map((f) => `${f}/**/*.scss`)
);
fixStyle.displayName = 'style lint autofix';

const fixJsTask = () => fixJs(
  lintJsGlob.map((f) => `${f}/**/*.js`)
);
fixJsTask.displayName = 'js lint autofix';

exports.lintStyleTask = lintStyleTask;
exports.lintJsTask = lintJsTask;
exports.dockerTask = dockerTask;
exports.watchTask = watchTask;
exports.liveTask = liveTask;
exports.styleTask = styleTask;
exports.fixStyleTask = fixStyleTask;
exports.fixJsTask = fixJsTask;
