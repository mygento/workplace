const debug = require('debug')('workplace:composer');
const spawn = require('child_process').spawn;

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

exports.composerCommand = (cb, command, config) => {
  debug(`${config.appDirectory}/composer.json`);
  if (process.env.NODE_ENV === 'production' || process.env.NODE_ENV === 'test') {
    cb();
    return;
  }
  runCommand(command, config.appDirectory, cb);
};
