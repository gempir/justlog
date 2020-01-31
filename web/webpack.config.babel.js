const path = require("path");
const webpack = require("webpack");
const PnpWebpackPlugin = require(`pnp-webpack-plugin`);

module.exports = (env, options) => ({

	entry: './src/index.jsx',
	output: {
		path: path.resolve(__dirname, 'public'),
		filename: 'bundle.js',
		publicPath: "/",
	},
	module: {
		rules: [
			{
				test: /\.jsx$/,
				exclude: /node_modules/,
				use: {
					loader: "babel-loader"
				}
			},
			{
				test: /\.s?css$/,
				use: ["style-loader", "css-loader", "sass-loader"]
			},
		]
	},
	resolve: {
		extensions: ['.js', '.jsx'],
		plugins: [
			PnpWebpackPlugin,
		],
	},
	plugins: [
		new webpack.DefinePlugin({
			'process.env': {
				'apiBaseUrl': options.mode === 'development' ? '"http://localhost:8025"' : '""',
			}
		}),
	],
	resolveLoader: {
		plugins: [
			PnpWebpackPlugin.moduleLoader(module),
		],
	},
	stats: {
		// Config for minimal console.log mess.
		assets: false,
		colors: true,
		version: false,
		hash: false,
		timings: true,
		chunks: false,
		chunkModules: false,
		entrypoints: false,
		modules: false,
	},
});