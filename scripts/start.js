'use strict';

process.env.BABEL_ENV = 'development';
process.env.NODE_ENV = 'development';

process.on('unhandledRejection', err => {
  throw err;
});
console.log('start');

const gulp = require('gulp');
const { someMagic, watchCode } = require('../gulpfile');
console.log('watch');
watchCode(['.']);
