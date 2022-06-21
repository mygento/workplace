import { compileStyle }  from './styles.js';
import { lintStyle, fixStyle }  from './lint.js';
import { watchStyles, watchLintStyles }  from './watch.js';

import { styleGlob, themeFolder, resolveApp, lintStyleGlob } from '@mygento/workplace/gulpfile/tasks.js';

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


const watchStylesTask = (cb) => {
  if (styleGlob.length === 0) {
    return cb();
  }
  return watchStyles(
    styleGlob.map((f) => `${f}/**/*.scss`)
  );
};
watchStylesTask.displayName = 'watch style';

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
fixStyleTask.displayName = 'style lint autofix';

const watchLintStylesTask = (cb) => {
  if (lintStyleGlob.length === 0) {
    return cb();
  }
  return watchLintStyles(
    lintStyleGlob.map((f) => `${f}/**/*.scss`)
  );
};
watchLintStylesTask.displayName = 'watch lint scss';

export {
  styleTask,
  watchStylesTask,
  lintStyleTask,
  fixStyleTask,
  watchLintStylesTask
};
