const path = require('path');
const glob = require('glob');
const HTMLWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

const entry = glob.sync('./src/pages/**/*.js').reduce((pre, filepath) => {
  const reg = /.*(\/.*)\1.js/g;
  const filename = path.basename(filepath);
  const name = path.parse(filename).name;
  
  if (name === 'index' || reg.test(filepath)) {
    return Object.assign(pre, {[name]: filepath});
  }
  return pre;
}, {});

const pages = glob.sync('./src/pages/**/*.tmpl').map((filepath) => {
  const parentsPath = /.*src\/(?<parents>.*)\/.*.tmpl/g.exec(filepath).groups.parents;
  const pathNodes = parentsPath.split('/');
  const pageChunks = pathNodes.reduce((chunks, node) => {
    if (entry.hasOwnProperty(node)) {
      return [...chunks, node]
    }
    return chunks;
  }, [])

  const distFilename = pathNodes.map(node => node.charAt(0).toUpperCase() + node.slice(1)).join('');

  return new HTMLWebpackPlugin({
    filename: `views/${distFilename}.tmpl`,
    template: filepath,
    chunks: ['vendors', 'common', 'manifest', 'base', ...pageChunks],
    chunksSortMode: 'manual',
  });
});

module.exports = {
  mode: 'development',
  entry,
  output: {
    path: path.resolve(__dirname, 'dist'),
    chunkFilename: 'static/js/[name].js',
    filename: 'static/js/[name].js',
    publicPath: '/',
  },
  plugins: [
    ...pages,
    new CleanWebpackPlugin(),
    new MiniCssExtractPlugin({
      filename: 'static/css/[name].css',
      chunkFilename: 'static/css/chunk.[id].css',
    }),
    new CopyWebpackPlugin([
      { from: 'src/asserts', to: 'static' },
      { from: 'src/layout/**/*.tmpl', to: 'views', flatten: true, },
    ], { copyUnmodified: true })
  ],
  module: {
    rules: [
      {
        test: /\.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          { loader: 'css-loader', options: { importLoaders: 1 } },
          {
            loader: 'postcss-loader', 
            options: {
              plugins: [
                require('autoprefixer')
              ]
            }
          },
        ],
      },
      {
        test: /\.scss$/,
        use: [
          MiniCssExtractPlugin.loader,
          { loader: 'css-loader', options: { importLoaders: 2 } },
          {
            loader: 'postcss-loader', 
            options: {
              plugins: [
                require('autoprefixer')
              ]
            }
          },
          'sass-loader',
        ],
      },
    ],
  },

  optimization: {
    runtimeChunk: {
      name: 'manifest'
    },
    splitChunks: {
      cacheGroups: {
        vendors: {
          test: /[\\/]node_modules[\\/]/,
          chunks: 'all',
          name: 'vendors',
          filename: 'static/js/vendors.js',
          priority: 2,
          reuseExistingChunk: true
        },
        common: {
          test: /\.m?js$/,
          chunks: 'all',
          name: 'common',
          filename: 'static/js/common.js',
          minSize: 0,
          minChunks: 2,
          priority: 1,
          reuseExistingChunk: true
        },
      }
    }
  }
}