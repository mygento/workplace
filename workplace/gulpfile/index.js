import pkg from 'gulp';
const { series, parallel } = pkg;

import {
  watchLintJsTask,
  startDockerTask, stopDockerTask, rmDockerTask,
  lintJsTask, fixJsTask,
  composerInstallTask
} from './tasks.js';

const loadModule = async (modulePath) => {
  try {
    return await import(modulePath)
  } catch (e) {
    return {};
  }
};

const {
  styleTask, watchStylesTask,
  lintStyleTask, fixStyleTask,
  watchLintStylesTask
} = await loadModule('@mygento/frontend-native');

const fallback = (cb) => cb();
fallback.displayName = 'do nothing';

export const install  = series(
  composerInstallTask
);

export const stop = series(
  stopDockerTask
);

export const remove = series(
  stopDockerTask,
  rmDockerTask
);

export const build =  parallel(
  styleTask ? styleTask : fallback
);

export const lint =  parallel(
  lintJsTask,
  lintStyleTask ? lintStyleTask: fallback
);

export const test = parallel(
  fixJsTask,
  fixStyleTask? fixStyleTask : fallback
);

export const start = parallel(
  startDockerTask,
  lintJsTask,
  watchLintJsTask,
  styleTask ? styleTask : fallback,
  watchStylesTask ? watchStylesTask : fallback,
  lintStyleTask ? lintStyleTask: fallback,
  watchLintStylesTask ? watchLintStylesTask : fallback
);
