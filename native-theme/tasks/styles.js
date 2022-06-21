import pkg from 'gulp';
import createDebug from 'debug';
import gdebug from 'gulp-debug';
import { obj as noop } from 'through2';
import dartSass from 'sass';
import gulpSass from 'gulp-sass';
import { join } from 'path';
import notify from 'gulp-notify';
import postcss from 'gulp-postcss';
import autoprefixer from 'autoprefixer';
import cssnano from 'cssnano';
import sourcemaps from 'gulp-sourcemaps';
import sassInheritance from 'gulp-sass-parent';

import sassime from 'sassime';
// import normalize from 'node-normalize-scss';

const { src, lastRun, dest } = pkg;
const debug = createDebug('workplace:styles');
const notProduction = process.env.NODE_ENV !== 'production';

export function compileStyle(files, task, themeFolders, cssOptions = {})
{
  debug('compile Style', files);
  const postcss_plugins = [
    autoprefixer(),
    cssnano(cssOptions)
  ];

  const includeOpt = {
    includePaths:
      [
      //normalize.includePaths,
        sassime.includePaths,
      ]
  };

  const sass = gulpSass(dartSass);

  return src(
    files,
    {
      allowEmpty: true,
      since: lastRun(task),
    }).pipe(gdebug({ title: 'style' }))
    .pipe(notProduction ? sourcemaps.init() : noop())
    .pipe(notProduction ? sassInheritance({ dir: themeFolders.map(t => `${t}/web/scss`) }) : noop())
    .pipe(gdebug({ title: 'style with parent' }))
    .pipe(sass.sync(includeOpt).on('error', sass.logError))
    .pipe(postcss(postcss_plugins))
    .pipe(notProduction ? sourcemaps.write('.') : noop())
    .pipe(gdebug({ title: 'write style' }))
    .pipe(dest(function(file) {
      return join(file.base, '../css');
    }))
    .pipe(notProduction ? notify({
      message: 'Styles complete',
      onLast: true,
    }) : noop());
}
