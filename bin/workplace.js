#!/usr/bin/env node
'use strict';

process.on('unhandledRejection', err => {
  throw err;
});

const spawn = require('cross-spawn');
const args = process.argv.slice(2);

const scriptIndex = args.findIndex(
  x => x === 'build' || x === 'start' || x === 'test'
);
const script = scriptIndex === -1 ? args[0] : args[scriptIndex];
// const nodeArgs = scriptIndex > 0 ? args.slice(0, scriptIndex) : [];

switch (script) {
  case 'build':
  case 'test':
  case 'start': {
    const result = spawn.sync(
      'gulp',
      ['-f', require.resolve('../gulpfile'), script],
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
    break;
  }
  default:
    console.log('Unknown script "' + script + '".');
    break;
}
