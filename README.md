# Workplace
Combined && Fast && Easy solution to create Magento workplace

## Config Schema

in ```package.json``` create ```workplace``` key
based on full list of posibilities.

Here full list of default options:

```js
{
  ...
  "workplace": {
    "type": "magento2", // magento2 | magento | symfony
    "php": {
      "image": "mygento/php:7.2-full" // php docker image
    },
    "nginx": {
      "image": "luckyraul/nginx:backports", // nginx docker image
      "port": 8081 // nginx listen port
    },
    "mysql": {
      "image": "mygento/mysql:5.7", // docker mysql image
      "port": 3306 //mysql listen port
    },
    "varnish": {
      "image": null, //docker varnish image
      "port": 8082 //varnish listen port
    },
    "redis": {
      "image": null, //docker redis image
      "port": 6379 //docker redis image
    },
    "elasticsearch": {
      "image": null, //docker elasticsearch image
      "port": 5900 //elasticsearch listen port
    },
    "clickhouse": {
      "image": null, //docker clickhouse image
      "port": 8123 //clickhouse listen port
    },
    "magento2": {
      "theme": [
        "app/design/frontend/XXX/yyy/web" // relative path to theme web folder
      ],
      "lint": [], // additional globs to lint
    },
    "livereload": true,
  }
  ...
}
```

Override by custom config in ```.workplace/config.local.json```

```js
{
  ...
  "php": {
    "image": "mygento/php:7.2-debug"
  },
  "nginx": {
    "port": 8082
  },
  ...
}
```
