pipelineJob('_NRF_promote_microservice') {

    concurrentBuild(false)


    logRotator(-1, 15, 1, -1)
    authorization {
        permission('hudson.model.Item.Read', 'anonymous')
        permission('hudson.model.Item.Read:authenticated')
        permission('hudson.model.Item.Build:authenticated')
        permission('hudson.model.Item.Cancel:authenticated')
        permission('hudson.model.Item.Workspace:authenticated')
    }

    parameters {
        stringParam {
            name('CHART_NAME')
            description("CHART name to be processed")
            defaultValue('')
        }
        stringParam {
            name('CHART_VERSION')
            description("Chart VERSION to be processed")
            defaultValue('')
        }
        stringParam {
            name('CHART_REPO')
            description("Chart REPOSITORY to be processed")
            defaultValue('')
        }
        stringParam {
            name('GERRIT_BRANCH')
            description("Branch that will be used to checkout repository and Jenkinsfile from 5gcicd/chassis-projects repository")
            defaultValue('master')
        }
        stringParam {
            name('GERRIT_REFSPEC')
            description("""Refspec that will be used to checkout repository and Jenkinsfile from 5gcicd/chassis-projects repository""")
            defaultValue('')
        }
        stringParam {
            name('RELEASE')
            description("""Indicates if a release should be made: true or false""")
            defaultValue('false')
        }
        stringParam {
            name('VERSION_STRATEGY')
            description("""Indicates strategy to follow when a release should be made: MAJOR,MINOR,PATCH""")
            defaultValue('PATCH')
        }
    }

    definition {
        cpsScm {
            scm {
                git {
                    remote {
                        name('origin')
                        url('https://esdccci@gerritmirror-ha.lmera.ericsson.se/a/5gcicd/chassis-projects')
                        credentials('userpwd-adp')
                        refspec('${GERRIT_REFSPEC}')
                        branch('${GERRIT_BRANCH}')
                    }

                    extensions {
                        wipeOutWorkspace()
                        choosingStrategy {
                            gerritTrigger()
                        }
                    }
                }
                scriptPath('integration/eric_nrf/Jenkinsfile.promote_ms')
            }
        }
    }
}
