import os from 'os';
import { spawn } from 'child_process';
import { join, resolve } from 'path';
import { writeFileSync, existsSync, mkdirSync } from 'fs';

import { PROJECT_FILE, updateGitignore } from './gitignore.js';

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
  projectName, name, image, port = false, portDefault, env = [], volume = []
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
    } : {},
    env.length > 0 ? {
      environment: env
    } : {}
  ),
} : {};

const optionalVolume = (image, name) => image ? { [name]: {} } : {};

const fileTemplate = (
  projectType,
  projectRoot,
  workplaceRoot,
  projectName,
  [phpImage, _phpPort, phpEnv = []],
  [nginxImage,  nginxPort, nginxEnv = []],
  [dbImage, dbPort, dbEnv = []],
  [redisImage = null, redisPort = null, redisEnv = []],
  [elasticImage = null, elasticPort = null, elasticEnv = []],
  [varnishImage = null, varnishPort = null, varnishEnv = []],
  [clickhouseImage = null, clickhousePort = null, clickhouseEnv = []],
  [rabbitmqImage = null, rabbitmqPort = null, rabbitmqEnv = []]
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
        environment: ['PHPFPM_USER=$USERID', ...phpEnv],
        depends_on: ['db'],
      },
      nginx: {
        container_name: `${projectName}-nginx`,
        image: nginxImage,
        networks: [networkName],
        environment: ['NGINX_USER=$USERID', ...nginxEnv],
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
        ports: [`${dbPort}:3306`],
        environment: [
          'MYSQL_ROOT_PASSWORD=mygento',
          `MYSQL_DATABASE=${projectType}`,
          'MYSQL_USER=mygento',
          'MYSQL_PASSWORD=mygento',
          ...dbEnv
        ],
        volumes: [`${getVolumeName(projectType, 'db')}:/var/lib/mysql`],
      },
    },
    optionalService(projectName, 'redis', redisImage, redisPort, 6379, redisEnv, []),
    optionalService(projectName, 'elastic', elasticImage, elasticPort, 9200, ['discovery.type=single-node', 'DISABLE_SECURITY_PLUGIN=true', ...elasticEnv], ['elastic:/usr/share/elasticsearch/data']),
    optionalService(projectName, 'varnish', varnishImage, varnishPort, 8081, varnishEnv, []),
    optionalService(projectName, 'clickhouse', clickhouseImage, clickhousePort, 8123, clickhouseEnv, ['clickhouse:/var/lib/clickhouse']),
    optionalService(projectName, 'rabbitmq', rabbitmqImage, rabbitmqPort, 5672, rabbitmqEnv, ['rabbitmq:/var/lib/rabbitmq'])
  ),
  volumes: Object.assign(
    {},
    { [getVolumeName(projectType, 'db')]: {} },
    optionalVolume(elasticImage, 'elastic'),
    optionalVolume(clickhouseImage, 'clickhouse'),
    optionalVolume(rabbitmqImage, 'rabbitmq')
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
  writeFileSync(
    join(config.appDirectory, PROJECT_FILE, 'docker-compose.json'),
    fileTemplate(
      config.type,
      config.appDirectory,
      resolve(`../nginx/${config.type}`),
      config.projectName,
      [config.php.image, config.php.port, config.php.env],
      [config.nginx.image,  config.nginx.port, config.nginx.env],
      [config.mysql.image,  config.mysql.port, config.mysql.env],
      getServiceConfig(config, 'redis'),
      getServiceConfig(config, 'elasticsearch'),
      getServiceConfig(config, 'varnish'),
      getServiceConfig(config, 'clickhouse'),
      getServiceConfig(config, 'rabbitmq')
    )
  );

  const cmd = spawn(
    'docker-compose',
    ['-f','docker-compose.json', ...command],
    { stdio: 'inherit', cwd: join(config.appDirectory, PROJECT_FILE) }
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

export function composeCommand(cb, command, config) {
  if (!projectTypes.includes(config.type)) {
    return cb();
  }
  updateGitignore(config.appDirectory);
  if (!existsSync(join(config.appDirectory, PROJECT_FILE))) {
    mkdirSync(join(config.appDirectory, PROJECT_FILE));
  }

  process.env.COMPOSE_PROJECT_NAME = config.projectName;
  process.env.USERID = os.userInfo().uid;

  writeFileSync(
    join(config.appDirectory, PROJECT_FILE, 'docker-compose.json'),
    fileTemplate(
      config.type,
      config.appDirectory,
      resolve(`../nginx/${config.type}`),
      config.projectName,
      [config.php.image, config.php.port, config.php.env],
      [config.nginx.image,  config.nginx.port, config.nginx.env],
      [config.mysql.image,  config.mysql.port, config.mysql.env],
      getServiceConfig(config, 'redis'),
      getServiceConfig(config, 'elasticsearch'),
      getServiceConfig(config, 'varnish'),
      getServiceConfig(config, 'clickhouse'),
      getServiceConfig(config, 'rabbitmq')
    )
  );

  runCommand(command, config, cb);
}
