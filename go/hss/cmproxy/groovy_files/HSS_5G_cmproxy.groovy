pipelineJob('HSS_5G_cmproxy') {
    concurrentBuild(false)
    parameters {
        stringParam {
            name('GERRIT_BRANCH')
            description('Use this parameter to select your repository BRANCH.')
            defaultValue('master')
        }
        stringParam {
            name('GERRIT_REFSPEC')
            description('Use same value as GERRIT_BRANCH.')
            defaultValue('master')
        }
    }

    triggers {
        gerrit {
            events {
                patchsetCreated()
                changeMerged()
            }

            project('plain:HSS/5G/cmproxy', ["ANT:**"])

            configure { gerritTrigger ->
                new groovy.util.Node(gerritTrigger, 'serverName', 'hss')
            }

            buildSuccessful(null, null)
        }
    }

    definition {
        cpsScm {
            scm {
                git {
                    remote {
                        url ('ssh://gerrit.ericsson.se:29418/HSS/5G/cmproxy')
                        credentials ('gerritpk-hss')
                        branch('${GERRIT_BRANCH}')
                        refspec('${GERRIT_REFSPEC}')
                    }
                    extensions {
                        choosingStrategy {
                            gerritTrigger()
                        }
                    }
                }
                lightweight (false)
                scriptPath ("Jenkinsfile")
            }
        }
    }
}
