const fs = require('fs')
const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const { CleanWebpackPlugin } = require('clean-webpack-plugin')

const resolve = (...ps) => {
    return path.resolve(__dirname, ...ps)
}

// TODO ??? https://github.com/TypeStrong/fork-ts-checker-webpack-plugin

const template = (config) => {
    const htmlWebpackPlugin = config.htmlWebpackPlugin
    const js = htmlWebpackPlugin.files.js

    console.log('const htmlWebpackPlugin =', htmlWebpackPlugin)

    const inline = (chunk) => {
        const name = js.find(n => n.startsWith(chunk))
        const path = resolve('dist', name)
        const data = fs.readFileSync(path, {encoding:'utf8', flag:'r'})
        fs.unlinkSync(path)
        return `<script>${data}</script>`
    }

    const link = (chunk) => {
        const name = js.find(n => n.startsWith(chunk))
        return `</script><script defer="defer" src="${name}"></script>`
    }

    // ${htmlWebpackPlugin.tags.headTags}
    // ${htmlWebpackPlugin.tags.bodyTags}

    return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>TheArchitect</title>
    <link rel="stylesheet" type="text/css" href="https://cdn.jsdelivr.net/gh/kneedeepincode/fsex-webfont@v1.0.1/fsex300.css">
    <link href="https://fonts.googleapis.com/css?family=Inconsolata&display=swap" rel="stylesheet">
</head>
<body>
<div id="terminal"></div>
${link('welcome')}
${link('webapp')}
</body>
</html>
`
}

module.exports = {
    mode: 'production',
    entry: {
        'welcome': './src/welcome/index.ts',
        'webapp': './src/webapp/index.ts',
    },
    /*devtool: 'inline-source-map',*/
    plugins: [
        new HtmlWebpackPlugin ({
            inject: false,
            scriptLoading: 'defer',
            // template: resolve('index.html')
            templateContent: template
        }),
        new CleanWebpackPlugin()
    ],
    module: {
        rules: [
            {
                test: [/\.ts$/],
                use: {
                    loader: 'ts-loader',
                    options: {
                        logLevel: 'info',
                        context: resolve(),
                        configFile: 'tsconfig.json',
                        allowTsInNodeModules: true
                    },
                },
                exclude: /node_modules/
            },
            {
                test: /\.styl$/,
                use: [
                    {
                        loader: "style-loader" // creates style nodes from JS strings
                    },
                    {
                        loader: "css-loader" // translates CSS into CommonJS
                    },
                    {
                        loader: "stylus-loader" // compiles Stylus to CSS
                    }
                ]
            },
            {
                test: /\.css$/,
                use: [
                    {
                        loader: "style-loader"
                    },
                    {
                        loader: "css-loader"
                    }
                ]
            },
            {
                test: /\.(png|jpg|gif|woff2?|ttf)(\?.+)?$/i,
                use: [
                    {
                        loader: 'url-loader',
                        options: {
                            limit: 256 * 1024,
                        },
                    },
                ],
            }
        ]
    },
    resolve: {
        extensions: ['.ts', '.css', '.styl', '.woff']
    },
    output: {
        filename: '[name]-[contenthash].js',
        path: resolve('dist')
    }
}
