import createDebug from 'debug';
const debug = createDebug('workplace:config');

const getServiceConfig = (config, override, service, workplaceConfig) => {
  if (
    !Object.prototype.hasOwnProperty.call(config, service) &&
    !Object.prototype.hasOwnProperty.call(override, service)) {
    return;
  }
  const { image: image1 = false, port: port1 = false  } = override[service] || {};
  const { image: image2 = false, port: port2 = false  } = config[service] || {};
  return workplaceConfig[service] = { image: image1 ? image1 : image2, port: port1 ? port1: port2 };
};

const mergeConfig = (config, appDirectory, override = {}) => {
  debug('project config', config.workplace);
  debug('override config', override);

  if (!config.workplace) {
    throw new Error('Empty workplace config');
  }
  const workplaceConfig = Object.assign({ appDirectory: appDirectory }, {
    projectName: config.name ? config.name : 'workplace',
    type: config.workplace.type || 'magento2',
    php: Object.assign({},
      {
        image: 'ghcr.io/mygento/php:7.3-full',
      },
      config.workplace.php,
      override.php
    ),
    nginx: Object.assign({},
      {
        image: 'luckyraul/nginx:backports',
        port: 8081,
      },
      config.workplace.nginx,
      override.nginx
    ),
    mysql: Object.assign({},
      {
        image: 'mygento/mysql:5.7',
        port: 3306,
      },
      config.workplace.mysql,
      override.mysql
    )
  });

  if (workplaceConfig.type === 'magento2') {
    workplaceConfig.magento2 = Object.assign({}, { style: true, theme: [], lint: [] }, config.workplace.magento2);
  }

  // console.log(config.workplace);

  getServiceConfig(config.workplace, override, 'redis', workplaceConfig);
  getServiceConfig(config.workplace, override, 'elasticsearch', workplaceConfig);
  getServiceConfig(config.workplace, override, 'varnish', workplaceConfig);
  getServiceConfig(config.workplace, override, 'clickhouse', workplaceConfig);
  getServiceConfig(config.workplace, override, 'rabbitmq', workplaceConfig);

  return workplaceConfig;
};

export { mergeConfig };
