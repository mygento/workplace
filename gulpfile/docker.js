const spawn = require('child_process').spawn;
const path = require('path');
const fs = require('fs');

const fileTemplate = (
  projectRoot,
  workplaceRoot,
  projectName,
  phpImage,
  nginxImage,
  nginxPort,
  dbImage,
  dbPort
) => JSON.stringify({
  version: '3.7',
  services: {
    php: {
      container_name: `${projectName}-php`,
      image: phpImage,
      links: ['db'],
      volumes: [`${projectRoot}:/var/www/public`],
      environment: ['PHPFPM_USER=$USERID'],
      depends_on: ['db'],
    },
    nginx: {
      container_name: `${projectName}-nginx`,
      image: nginxImage,
      environment: ['NGINX_USER=$USERID'],
      volumes: [
        `${workplaceRoot}:/etc/nginx/sites-enabled/`,
        `${projectRoot}:/var/www/public`,
      ],
      ports: [`${nginxPort}:80`],
      depends_on: ['php'],
    },
    db: {
      container_name: `${projectName}-db`,
      image: dbImage,
      command: 'mysqld --character-set-server=utf8 '
      + ' --collation-server=utf8_general_ci'
      + ' --innodb_file_per_table',
      ports: [`${dbPort}:3306`],
      environment: [
        'MYSQL_ROOT_PASSWORD=mygento',
        'MYSQL_DATABASE=magento',
        'MYSQL_USER=mygento',
        'MYSQL_PASSWORD=mygento',
      ],
      volumes: ['m2-db-data:/var/lib/mysql'],
    },
  },
  volumes: {
    'm2-db-data': {},
  },
}, null, 2);

exports.composeStart = (cb, config) => {
  process.env.COMPOSE_PROJECT_NAME = config.projectName;
  process.env.USERID = require('os').userInfo().uid;
  fs.writeFileSync(
    'docker-compose.json',
    fileTemplate(
      config.appDirectory,
      path.resolve('../nginx/'),
      config.projectName,
      config.php,
      config.nginx.image,
      config.nginx.port,
      config.mysql.image,
      config.mysql.port
    )
  );
  const cmd = spawn(
    'docker-compose',
    ['-f','docker-compose.json', 'up', '-d'],
    { stdio: 'inherit' }
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
