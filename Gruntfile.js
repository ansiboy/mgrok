/* 1. 需要安装 go-bindata，运行下面的命令：
 * go get -u github.com/jteeuwen/go-bindata/...
 */

const path = require("path");
var cert_name = '*.mgrok.cn';

module.exports = function (grunt) {


    let is_win = /^win/.test(process.platform);         


    let build_tag = grunt.option('build');
    if(build_tag!='debug' && build_tag != 'release') {
        console.log(`
usage:
grunt --build=debug
grunt --build=release
        `);
    }

    let working_path = path.resolve("./");
    let set_gopath_command = setEnvVariableCommand("GOPATH",working_path);
    let output_path = `bin/${build_tag}`;
    
    grunt.initConfig({
        copy: {
            cert: {
                files: build_tag == 'debug' ?
                [
                    { expand: true, cwd: `cert/${cert_name}_debug/`, src: `ngrokroot.crt`, dest: 'assets/client/tls/' },
                    { expand: true, cwd: `cert/${cert_name}_debug/`, src: `snakeoil.crt`, dest: `assets/server/tls/` },
                    { expand: true, cwd: `cert/${cert_name}_debug/`, src: `snakeoil.key`, dest: `assets/server/tls/` }
                ] :
                [
                    { expand: true, cwd: `cert/${cert_name}/`, src: `ngrokroot.crt`, dest: 'assets/client/tls/' },
                    { expand: true, cwd: `cert/${cert_name}/`, src: `snakeoil.crt`, dest: `assets/server/tls/` },
                    { expand: true, cwd: `cert/${cert_name}/`, src: `snakeoil.key`, dest: `assets/server/tls/` }
                ] 
            },
            ngrok_config: {
                files: [
                    { expand: true, cwd: `src/ngrok/main/ngrok/`, src: '.ngrok', dest: output_path }
                ]
            }
        },
        shell: {
            options: {
                stderr: true
            },
            build: {
                command: [
                    setEnvVariableCommand("GOPATH", working_path),
                    `go-bindata -nomemcopy -pkg=assets -tags=${build_tag} -debug=${build_tag == 'debug' ? 'true' : 'false'} -o=src/ngrok/client/assets/assets_${build_tag}.go assets/client/...`,
                    `go build -o ${output_path}/ngrok${is_win?'.exe':''} -tags "${build_tag}"  src/ngrok/main/ngrok/ngrok.go`,
                    `go-bindata -nomemcopy -pkg=assets -tags=${build_tag} -debug=${build_tag == 'debug' ? 'true' : 'false'} -o=src/ngrok/server/assets/assets_${build_tag}.go assets/server/...`,
                    `go build -o ${output_path}/ngrokd${is_win?'.exe':''} -tags "${build_tag}"  src/ngrok/main/ngrokd/ngrokd.go`
                    
                ].join("&&")
            },
            all_client:{
                command: [
                    `sudo GOOS=linux GOARCH=amd64 make release-client`,
                    `sudo GOOS=linux GOARCH=386 make release-client`,
                    `sudo GOOS=linux GOARCH=arm make release-client`,
                    `sudo GOOS=windows GOARCH=amd64 make release-client`,
                    `sudo GOOS=windows GOARCH=386 make release-client`,
                    `sudo GOOS=darwin GOARCH=amd64 make release-client`,
                    `sudo GOOS=darwin GOARCH=386 make release-client`,
                ].join('&&')
            }
        }
    });

    
    function setEnvVariableCommand(name, value){
        return is_win ? `SET ${name}=${value}` : `${name}=${value}` ;
     }

    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-shell');

    let tasks = ['copy', 'shell:build'];
    if(build_tag == 'release') {
        tasks.push('shell:all_client')
    }
    grunt.registerTask('default', tasks);
}

