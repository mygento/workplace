import { describe, it } from 'mocha';

import { assert } from 'chai';
import { mergeConfig } from '../gulpfile/config.js';

describe('Config', function() {
  describe('parse', function() {
    it('error on empty config', function() {
      assert.throws(() => {
        mergeConfig({}, '');
      }, Error, 'Empty workplace config');
    });

    it('magento2 default', function() {
      const config = mergeConfig({ workplace: { type: 'magento2' } }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports',port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom mysql port', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', mysql: { port: 3307 } },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3307 },
          nginx: { image: 'luckyraul/nginx:backports',port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom php port', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', php: { port: 9999 } },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports',port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full', port: 9999 },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom php env', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', php: { env: ['XDEBUG_SESSION=ide'] } },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports',port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full', env: ['XDEBUG_SESSION=ide'] },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom mysql image', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', mysql: { image: 'mygento/mysql:5.6' } },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.6', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports',port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom nginx image', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', nginx: { image: 'luckyraul/nginx' } },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx',port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom nginx port', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', nginx: { port: 8082 } },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8082 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom nginx image and port', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', nginx: {
          image: 'luckyraul/nginx', port: 8082 },
        },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx', port: 8082 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('custom php image', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', php: { image: 'mygento/php:7.3-full' } },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8081 },
          php: { image: 'mygento/php:7.3-full' },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('redis', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', redis: { image: 'redis' } },
      }, '');
      console.log(config);
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          redis: {
            image: 'redis',
            port: false,
            environment: [],
          },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('elastic', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', elasticsearch: { image: 'elastic', port: 9300,  environment: ['XXX=1'] } },
      }, '');
      console.log(config);
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          elasticsearch: {
            image: 'elastic',
            port: 9300,
            environment: ['XXX=1'],
          },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('varnish', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', varnish: { image: 'varnish' } },
      }, '');
      console.log(config);
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          varnish: {
            image: 'varnish',
            port: false,
            environment: [],
          },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });

    it('clickhouse', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', clickhouse: { image: 'clickhouse', port: 123 } },
      }, '');
      console.log(config);
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8081 },
          php: { image: 'ghcr.io/mygento/php:7.3-full' },
          clickhouse: {
            image: 'clickhouse',
            port: 123,
            environment: [],
          },
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: [],
            lint: [],
            style: true
          }
        }
      );
    });
  });
});
