@Library("PipelineGlobalLibrary") _
import com.ericsson.ci.mason.Mason2

Mason2.ciPipeline(this) {
    cloud("kubernetes")
    languages('go')
    unit('nef')
        parallel(false)
        checkout  {
            gerrit('checkout-nef-nrfintegration') {
                project('5gcicd/chassis-projects')
                credentials('userpwd-adp')
                branch('c417bb8d3ad542e10942a2401c4c5c16379a54bc')
                refspec('refs/changes/87/6006487/1')
                checkoutdir {
                    subdir = 'chassis-projects'
                }
            }
        }
        ut {
            go("go-test") {
                container('nef-go-build', 'selidocker.lmera.ericsson.se/proj-nef/base-image-redhat-golang-builder:1.12.5')
                envVars(['GO111MODULE':'on'])
                wsVolumeMount(mountPath: '/go/src/gerrit.ericsson.se/nef', subPath:'chassis-projects/go/nef/src')
                wsVolumeMount(mountPath: '/go/src/gerrit.ericsson.se/nef', subPath:'chassis-projects/go/nef/go.mod')
                privileged(true)
                gotestparameters('-v -cover')
                // dir('udm/dns')
                // excludepaths(['pkg/mgmclient', 'pkg/ntfhandler', 'pkg/localcache', 'nrfdns'])
            }
        }
        build {
            go("test4-step7-go-build-files") {
                // dir('udm/dns')
                wsVolumeMount(mountPath: '/go/src/gerrit.ericsson.se/nef', subPath:'chassis-projects/go/nef/src')
                wsVolumeMount(mountPath: '/go/src/gerrit.ericsson.se/nef', subPath:'chassis-projects/go/nef/go.mod')
                privileged(true)
                
                buildFlags("""-a -installsuffix cgo -ldflags '-extldflags "-static"'""")
                // excludepaths(['pkg/mgmclient', 'pkg/ntfhandler', 'pkg/localcache', 'nrfdns'])
                wsVolumeMount('/go/src/gerrit.ericsson.se')

            }
        }
        coverity {
            go("test-go-coverity") {
                wsVolumeMount('/go/src/gerrit.ericsson.se')
                dir('udm/dns')
                excludepaths(['pkg/mgmclient', 'pkg/ntfhandler', 'pkg/localcache', 'nrfdns'])
                precommand ('echo ============ Hello this is precommand =============')
            }
        }
        sonar {
            scanner("test06-step4-scanner") {
                abortOnError(true)
                wsVolumeMount('/go/src/gerrit.ericsson.se')
                version('fantastic-version')
                url("https://udm5g-qg-test.lmera.ericsson.se")
                project("perfectProject")
                dir('udm/dns')
                token('sonar-token-adp')
                exclusionPath(['pkg/mgmclient', 'pkg/ntfhandler', 'pkg/localcache', 'nrfdns'])
                // src('java/maven')
            }
        }
}