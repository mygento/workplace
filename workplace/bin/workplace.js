#!/usr/bin/env node
'use strict';
import { sync } from 'cross-spawn';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

process.on('unhandledRejection', err => {
  throw err;
});

const args = process.argv.slice(2);

const scriptIndex = args.findIndex(
  x => x === 'build' || x === 'start' || x === 'test'
);
const script = scriptIndex === -1 ? args[0] : args[scriptIndex];
// const nodeArgs = scriptIndex > 0 ? args.slice(0, scriptIndex) : [];

switch (script) {
  case 'install':
  case 'build':
  case 'test':
  case 'stop':
  case 'lint':
  case 'delete':
  case 'start': {
    const result = sync(
      'gulp',
      ['-f', resolve(__dirname, '../gulpfile/index.js'), script],
      { stdio: 'inherit' }
    );
    if (result.signal) {
      if (result.signal === 'SIGKILL') {
        console.log(
          'The system ran out of memory or someone called kill on the process.'
        );
      } else if (result.signal === 'SIGTERM') {
        console.log(
          'Someone might have called `kill` or `killall`'
        );
      }
      process.exit(1);
    }
    process.exit(result.status);
  }
  default:
    console.log('Unknown script "' + script + '".');
    break;
}
