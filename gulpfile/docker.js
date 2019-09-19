const spawn = require('child_process').spawn;

exports.compose = (cb) => {
  const cmd = spawn('docker-compose', ['up', '-d'], { stdio: 'inherit' });
  cmd.on('close', function(code) {
    console.log('my-task exited with code ' + code);
    cb(code);
  });
  cmd.on('error', function(code) {
    console.log('my-task exited with code ' + code);
    cb(code);
  });
};
