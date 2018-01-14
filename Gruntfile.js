/* 1. 需要安装 go-bindata，运行下面的命令：
 * go get -u github.com/jteeuwen/go-bindata/...
 */

const path = require("path");
var cert_name = '*.mgrok.cn';
// var build_tag = 'debug'
module.exports = function (grunt) {

    let is_win = /^win/.test(process.platform);

    let build_tag = grunt.option('build');
    if (build_tag != 'debug' && build_tag != 'release') {
        console.log(`
usage:
grunt --build=debug
grunt --build=release
        `);
        return;
    }

    let working_path = path.resolve("./");
    let set_gopath_command = is_win ? `SET GOPATH=${working_path}` : `GOPATH=${working_path}`;

    console.log(set_gopath_command);

    grunt.initConfig({
        copy: {
            cert: {
                files: [
                    { expand: true, cwd: `cert/${cert_name}/`, src: `ngrokroot.crt`, dest: 'assets/client/tls/' },
                    { expand: true, cwd: `cert/${cert_name}/`, src: `snakeoil.crt`, dest: `assets/server/tls/` },
                    { expand: true, cwd: `cert/${cert_name}/`, src: `snakeoil.key`, dest: `assets/server/tls/` }
                ]
            },
            ngrok_config: {
                files: [
                    { expand: true, cwd: `src/ngrok/main/ngrok/`, src: '.ngrok', dest: 'bin' }
                ]
            }
        },
        shell: {
            options: {
                stderr: true
            },
            build_client: {
                command: `
                    ${set_gopath_command}
                    go-bindata -nomemcopy -pkg=assets -tags=${build_tag} -debug=${build_tag == 'debug' ? 'true' : 'false'} -o=src/ngrok/client/assets/assets_${build_tag}.go assets/client/...
                    go build -o bin/ngrok -tags "${build_tag}"  src/ngrok/main/ngrok/ngrok.go
                `
            },
            build_server: {
                command: `
                    ${set_gopath_command}
                    go-bindata -nomemcopy -pkg=assets -tags=${build_tag} -debug=${build_tag == 'debug' ? 'true' : 'false'} -o=src/ngrok/server/assets/assets_${build_tag}.go assets/server/...
                    go build -o bin/ngrokd -tags "${build_tag}"  src/ngrok/main/ngrokd/ngrokd.go
                `
            },
            build_ngrok: {
                command: [
                    set_gopath_command,
                    `go-bindata -nomemcopy -pkg=assets -tags=${build_tag} -debug=${build_tag == 'debug' ? 'true' : 'false'} -o=src/ngrok/server/assets/assets_${build_tag}.go assets/server/...`,
                    `go build -o bin/ngrok -tags "${build_tag}"  src/ngrok/main/ngrok/ngrok.go`
                ].join('&&')
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-shell');
    grunt.registerTask('default', ['copy', 'shell']);
}