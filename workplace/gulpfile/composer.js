import createDebug from 'debug';
const debug = createDebug('workplace:composer');
import { spawn } from 'child_process';

const runCommand = (command, directory, cb, ignore = true) => {
  const commandList = ignore ? [command, '--ignore-platform-reqs'] : [command];
  const cmd = spawn(
    'composer',
    commandList,
    { stdio: 'inherit', cwd: directory }
  );
  cmd.on('close', function(code) {
    if (code !== 0) {
      console.log('docker exited on close with code ' + code);
    }
    cb(code);
  });
  cmd.on('error', function(code) {
    console.log('docker exited on error with code ' + code);
    cb(code);
  });
};

export function composerCommand(cb, command, config) {
  debug(`${config.appDirectory}/composer.json`);
  if (process.env.NODE_ENV === 'production' || process.env.NODE_ENV === 'test' || !!process.env.CI) {
    cb();
    return;
  }
  runCommand(command, config.appDirectory, cb);
}
