pipelineJob('_NRF_spinnaker_test') {

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
            description('Helm chart name that will be tested')
            defaultValue('')
        }
        stringParam {
            name('CHART_REPO')
            description('Helm chart repository URL')
            defaultValue('')
        }
        stringParam {
            name('CHART_VERSION')
            description('Helm chart version')
            defaultValue('')
        }
        stringParam {
            name('GERRIT_REFSPEC')
            description("Refspec that will clone 5gcicd/chassis-projects repository")
            defaultValue(' ')
        }
        stringParam {
            name('GERRIT_BRANCH')
            description("Branch that will be used to clone 5gcicd/chassis-projects repository")
            defaultValue('master')
        }
        stringParam {
            name('CLOUD')
            description('To choose a different cluster from default one.<br><br>')
            defaultValue('hoff005')
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
                    }
                    branch('${GERRIT_BRANCH}')
                    extensions {
                        wipeOutWorkspace()
                    }
                }
                scriptPath('integration/eric_nrf/Jenkinsfile.test')
            }
        }
    }
}
