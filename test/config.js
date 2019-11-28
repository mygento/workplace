const { describe, it } = require('mocha');

const { assert } = require('chai');
const mergeConfig = require('../gulpfile/config').mergeConfig;

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
          php: 'mygento/php:7.2-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: true
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
          php: 'mygento/php:7.2-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: true
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
          php: 'mygento/php:7.2-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: true
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
          php: 'mygento/php:7.2-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: true
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
          php: 'mygento/php:7.2-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: true
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
          php: 'mygento/php:7.2-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: true
        }
      );
    });

    it('custom php image', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', php: 'mygento/php:7.3-full' },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8081 },
          php: 'mygento/php:7.3-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: true
        }
      );
    });

    it('custom livereload', function() {
      const config = mergeConfig({
        workplace: { type: 'magento2', livereload: false },
      }, '');
      assert.deepEqual(
        config,
        {
          appDirectory: '',
          mysql: { image: 'mygento/mysql:5.7', port: 3306 },
          nginx: { image: 'luckyraul/nginx:backports', port: 8081 },
          php: 'mygento/php:7.2-full',
          projectName: 'workplace',
          type: 'magento2',
          magento2: {
            theme: []
          },
          livereload: false
        }
      );
    });
  });
});
