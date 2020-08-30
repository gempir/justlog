const path = require("path");
const webpack = require("webpack");
const fs = require("fs");
const PnpWebpackPlugin = require(`pnp-webpack-plugin`);

module.exports = (env, options) => {
	const plugins = [
		new webpack.DefinePlugin({
			'process.env': {
				'apiBaseUrl': options.mode === 'development' ? '"http://localhost:8025"' : '""',
			}
		}),
	];

	if (options.mode === "production") {
		plugins.push(new HashApplier());
	}

	return {
		entry: './src/index.jsx',
		output: {
			path: path.resolve(__dirname, 'public'),
			filename: options.mode === "production" ? 'bundle.[hash].js' : "bundle.js",
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
		plugins: plugins,
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
		}
	}
};

class HashApplier {
	apply(compiler) {
		compiler.hooks.done.tap('hash-applier', data => {
			const fileContents = fs.readFileSync(__dirname + "/public/index.html", "utf8");
			const newFileContents = fileContents.replace(
				/<!-- webpack-bundle-start -->(.*)?<!-- webpack-bundle-end -->/,
				`<!-- webpack-bundle-start --><script src="/bundle/bundle.${data.hash}.js"></script><!-- webpack-bundle-end -->`
			);

			fs.writeFileSync(__dirname + "/public/index.html", newFileContents);
		});
	}
}