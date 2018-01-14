var cert_name = '*.mgrok.cn';
// var build_tag = 'debug'
module.exports = function(grunt){
    let build_tag = grunt.option('build');
    if(build_tag != 'debug' && build_tag != 'release'){
        console.log(`
usage:
grunt --build=debug
grunt --build=release
        `);
        return;
    }
    grunt.initConfig({
        copy: {
            cert:{
                files:[
                    { expand: true, cwd:`cert/${cert_name}/`, src: `ngrokroot.crt`, dest: 'assets/client/tls/' },
                    { expand: true, cwd:`cert/${cert_name}/`, src: `snakeoil.crt`, dest: `assets/server/tls/`},
                    { expand: true, cwd:`cert/${cert_name}/`, src: `snakeoil.key`, dest: `assets/server/tls/`}
                ]
            }
        },
        shell: {
            build_client: {
                command: `
                    GOPATH=/home/maishu/projects/mgrok
                    go-bindata -nomemcopy -pkg=assets -tags=${build_tag} -debug=${build_tag == 'debug' ? 'true' : 'false'} -o=src/ngrok/client/assets/assets_${build_tag}.go assets/client/...
                    go build -o bin/ngrok -tags "${build_tag}"  src/ngrok/main/ngrok/ngrok.go
                `
            },
            build_server: {
                command: `
                    GOPATH=/home/maishu/projects/mgrok
                    go-bindata -nomemcopy -pkg=assets -tags=${build_tag} -debug=${build_tag == 'debug' ? 'true' : 'false'} -o=src/ngrok/server/assets/assets_${build_tag}.go assets/server/...
                    go build -o bin/ngrokd -tags "${build_tag}"  src/ngrok/main/ngrokd/ngrokd.go
                `
            },
            ngrok_config: {
                command: `cp src/ngrok/main/ngrok/.ngrok bin`
            }
        }
    });

    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-shell');
    grunt.registerTask('default', ['copy', 'shell']);
}