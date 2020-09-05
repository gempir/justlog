const path = require("path");
const webpack = require("webpack");
const fs = require("fs");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = (env, options) => {
	const plugins = [
		new webpack.DefinePlugin({
			'process.env': {
				'apiBaseUrl': options.mode === 'development' ? '"http://localhost:8025"' : '""',
			}
		}),
		new HtmlWebpackPlugin({
			template: path.resolve(__dirname, 'src', 'index.html')
		})
	];

	return {
		entry: './src/index.jsx',
		output: {
			path: path.resolve(__dirname, 'public'),
			filename: 'bundle.[hash].js',
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
		},
		plugins: plugins,
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
		}
	}
};