import pkg from 'gulp';
const { series, parallel } = pkg;

import {
  watchLintJsTask,
  startDockerTask, stopDockerTask, rmDockerTask,
  lintJsTask, fixJsTask,
  composerInstallTask
} from './tasks.js';

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

export const lint =  parallel(
  lintJsTask
);

export const test = parallel(
  fixJsTask
);

export const start = parallel(
  watchLintJsTask,
  startDockerTask,
  lintJsTask
);
