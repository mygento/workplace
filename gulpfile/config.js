const getServiceConfig = (config, service, workplaceConfig) => {
  if (!Object.prototype.hasOwnProperty.call(config, service)) {
    return;
  }
  const { image = false, port = false  } = config[service];
  return workplaceConfig[service] = { image, port };
};

const mergeConfig = (config, appDirectory) => {
  if (!config.workplace) {
    throw new Error('Empty workplace config');
  }
  const workplaceConfig = Object.assign({ appDirectory: appDirectory }, {
    projectName: config.name ? config.name : 'workplace',
    type: config.workplace.type || 'magento2',
    php: config.workplace.php || 'mygento/php:7.2-full',
    nginx: Object.assign({
      image: 'luckyraul/nginx:backports',
      port: 8081,
    }, config.workplace.nginx || {}),
    mysql: Object.assign({
      image: 'mygento/mysql:5.7',
      port: 3306,
    }, config.workplace.mysql || {}),
    livereload: config.workplace.livereload !== undefined ? config.workplace.livereload : true
  });

  if (workplaceConfig.type === 'magento2') {
    workplaceConfig.magento2 = workplaceConfig.magento2 || {};
    workplaceConfig.magento2.theme = workplaceConfig.magento2.theme || [];
  }

  if (workplaceConfig.type === 'magento1') {
    workplaceConfig.magento1 = workplaceConfig.magento1 || {};
    workplaceConfig.magento1.src = workplaceConfig.magento1.src || 'src';
    workplaceConfig.magento1.dest = workplaceConfig.magento1.dest || 'public';
  }

  getServiceConfig(config.workplace, 'redis', workplaceConfig);
  getServiceConfig(config.workplace, 'elasticsearch', workplaceConfig);
  getServiceConfig(config.workplace, 'varnish', workplaceConfig);
  getServiceConfig(config.workplace, 'clickhouse', workplaceConfig);

  // console.log('real', workplaceConfig);
  return workplaceConfig;
};

exports.mergeConfig = mergeConfig;
