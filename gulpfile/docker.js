const spawn = require('child_process').spawn;
const path = require('path');
const fs = require('fs');

const { PROJECT_FILE, updateGitignore } = require('./gitignore');

const getVolumeName = (projectType, type) => `${projectType}-${type}`;
const networkName = 'net';

const projectTypes = ['magento1', 'magento2', 'symfony'];

const getServiceConfig = (config, service) => {
  if (!Object.prototype.hasOwnProperty.call(config, service)) {
    return [false, false];
  }
  const { image = false, port = false  } = config[service];
  return [image, port];
};

const optionalService = (
  projectName, name, image, port = false, portDefault, volume = []
) => image ? {
  [name]: Object.assign(
    {},
    {
      container_name: `${projectName}-${name}`,
      networks: [networkName],
      image: image,
    },
    port ? {
      ports: [`${port}:${portDefault}`],
    } : {},
    volume.length > 0 ? {
      volumes: volume,
    } : {}
  ),
} : {};

const optionalVolume = (image, name) => image ? { [name]: {} } : {};

const fileTemplate = (
  projectType,
  projectRoot,
  workplaceRoot,
  projectName,
  phpImage,
  [nginxImage,  nginxPort],
  [dbImage, dbPort],
  [redisImage = null, redisPort = null],
  [elasticImage = null, elasticPort = null],
  [varnishImage = null, varnishPort = null],
  [clickhouseImage = null, clickhousePort = null]
) => JSON.stringify({
  version: '3.7',
  services: Object.assign(
    {},
    {
      php: {
        container_name: `${projectName}-php`,
        image: phpImage,
        networks: [networkName],
        volumes: [`${projectRoot}:/var/www/${projectType}`],
        environment: ['PHPFPM_USER=$USERID'],
        depends_on: ['db'],
      },
      nginx: {
        container_name: `${projectName}-nginx`,
        image: nginxImage,
        networks: [networkName],
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
        networks: [networkName],
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
        volumes: [`${getVolumeName(projectType, 'db')}:/var/lib/mysql`],
      },
    },
    optionalService(projectName, 'redis', redisImage, redisPort, 6379, []),
    optionalService(projectName, 'elastic', elasticImage, elasticPort, 9200, ['elastic:/usr/share/elasticsearch/data']),
    optionalService(projectName, 'varnish', varnishImage, varnishPort, 8081, []),
    optionalService(projectName, 'clickhouse', clickhouseImage, clickhousePort, 8123, ['clickhouse:/var/lib/clickhouse'])
  ),
  volumes: Object.assign(
    {},
    { [getVolumeName(projectType, 'db')]: {} },
    optionalVolume(elasticImage, 'elastic'),
    optionalVolume(clickhouseImage, 'clickhouse')
  ),
  networks: {
    [networkName]: {
      driver: 'bridge',
      driver_opts: {
        'com.docker.network.enable_ipv6': 'false',
        'com.docker.network.bridge.name': projectName
      }
    }
  }
}, null, 2);

const runCommand = (command, config, cb) => {
  fs.writeFileSync(
    path.join(config.appDirectory, PROJECT_FILE, 'docker-compose.json'),
    fileTemplate(
      config.type,
      config.appDirectory,
      path.resolve(`../nginx/${config.type}`),
      config.projectName,
      config.php,
      [config.nginx.image,  config.nginx.port],
      [config.mysql.image,  config.mysql.port],
      getServiceConfig(config, 'redis'),
      getServiceConfig(config, 'elasticsearch'),
      getServiceConfig(config, 'varnish'),
      getServiceConfig(config, 'clickhouse')
    )
  );

  const cmd = spawn(
    'docker-compose',
    ['-f','docker-compose.json', ...command],
    { stdio: 'inherit', cwd: path.join(config.appDirectory, PROJECT_FILE) }
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

exports.composeCommand = (cb, command, config) => {
  if (!projectTypes.includes(config.type)) {
    return cb();
  }
  updateGitignore(config.appDirectory);
  if (!fs.existsSync(path.join(config.appDirectory, PROJECT_FILE))) {
    fs.mkdirSync(path.join(config.appDirectory, PROJECT_FILE));
  }

  process.env.COMPOSE_PROJECT_NAME = config.projectName;
  process.env.USERID = require('os').userInfo().uid;

  fs.writeFileSync(
    path.join(config.appDirectory, PROJECT_FILE, 'docker-compose.json'),
    fileTemplate(
      config.type,
      config.appDirectory,
      path.resolve(`../nginx/${config.type}`),
      config.projectName,
      config.php,
      [config.nginx.image,  config.nginx.port],
      [config.mysql.image,  config.mysql.port],
      getServiceConfig(config, 'redis'),
      getServiceConfig(config, 'elasticsearch'),
      getServiceConfig(config, 'varnish'),
      getServiceConfig(config, 'clickhouse')
    )
  );

  runCommand(command, config, cb);
};
