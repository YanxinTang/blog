import gulp from 'gulp';
import gulpif from 'gulp-if';
import sass from 'gulp-sass';
import useref from 'gulp-useref';
import uglify from 'gulp-uglify';
import minifyCss from 'gulp-clean-css';
import htmlmin from 'gulp-htmlmin';

const copyStaticFile = () => {
  return gulp.src('src/static/images/**/*')
    .pipe(gulp.dest('dist/static/images'));
}

const buildTmpl = () => {
  return gulp.src('src/**/*.tmpl')
    .pipe(useref({searchPath: 'src'}))
    .pipe(gulpif('*.js', uglify()))
    .pipe(gulpif('*.css', sass()))
    .pipe(gulpif('*.css', minifyCss()))
    .pipe(gulpif('*.tmpl', htmlmin({ collapseWhitespace: true })))
    .pipe(gulp.dest('dist'))
}

const build = gulp.parallel(buildTmpl, copyStaticFile);

export {
  build
}