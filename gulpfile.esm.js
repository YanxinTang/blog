import gulp from 'gulp';
import gulpif from 'gulp-if';
import sass from 'gulp-sass';
import useref from 'gulp-useref';
import uglify from 'gulp-uglify';
import minifyCss from 'gulp-clean-css';
import htmlmin from 'gulp-htmlmin';
import babel from 'gulp-babel';

const copyStaticFile = () => {
  return gulp.src('public/**/*')
    .pipe(gulp.dest('dist/static/'));
}

const buildTmpl = () => {
  return gulp.src('src/**/*.tmpl')
    .pipe(useref())
    .pipe(gulpif('*.js', babel({ presets: ['@babel/env'] })))
    .pipe(gulpif('*.js', uglify()))
    .pipe(gulpif('*.css', sass()))
    .pipe(gulpif('*.css', minifyCss()))
    .pipe(gulpif('*.tmpl', htmlmin({
      collapseWhitespace: true,
      ignoreCustomFragments: [/{{[\s\S]*?}}/],
    })))
    .pipe(gulp.dest('dist'))
}

const build = gulp.parallel(buildTmpl, copyStaticFile);

const start = () => {
  return gulp.watch(['src'], function bundleTemplate () {
    return gulp.src('src/**/*.tmpl')
      .pipe(useref())
      .pipe(gulpif('*.css', sass()))
      .pipe(gulp.dest('dist'))
  })
}

export {
  start,
  build
}