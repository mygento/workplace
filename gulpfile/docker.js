const spawn = require('child_process').spawn;
const path = require('path');
const fs = require('fs');

const getVolumeName = (projectType) => `${projectType}-db-data`;

const projectTypes = ['magento2', 'symfony'];

const fileTemplate = (
  projectType,
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
      volumes: [`${projectRoot}:/var/www/${projectType}`],
      environment: ['PHPFPM_USER=$USERID'],
      depends_on: ['db'],
    },
    nginx: {
      container_name: `${projectName}-nginx`,
      image: nginxImage,
      environment: ['NGINX_USER=$USERID'],
      volumes: [
        `${workplaceRoot}:/etc/nginx/sites-enabled/`,
        `${projectRoot}:/var/www/${projectType}`,
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
        `MYSQL_DATABASE=${projectType}`,
        'MYSQL_USER=mygento',
        'MYSQL_PASSWORD=mygento',
      ],
      volumes: [`${getVolumeName(projectType)}:/var/lib/mysql`],
    },
  },
  volumes: {
    [getVolumeName(projectType)]: {},
  },
}, null, 2);

exports.composeStart = (cb, config) => {
  if (!projectTypes.includes(config.type)) {
    return cb();
  }

  process.env.COMPOSE_PROJECT_NAME = config.projectName;
  process.env.USERID = require('os').userInfo().uid;

  fs.writeFileSync(
    'docker-compose.json',
    fileTemplate(
      config.type,
      config.appDirectory,
      path.resolve(`../nginx/${config.type}`),
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
