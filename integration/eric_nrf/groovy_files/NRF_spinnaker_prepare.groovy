pipelineJob('_NRF_spinnaker_prepare') {

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
            description('Helm chart name part of requirements.yaml for eric_nrf')
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
            name('GERRIT_BRANCH')
            description("Branch that will be used to clone Jenkinsfile from 5gcicd/chassis-projects repository")
            defaultValue('master')
        }
        stringParam {
            name('GERRIT_REFSPEC')
            description("""Refspec for 5gcicd/chassis-projects repository. This parameter takes prevalence over
            the CHART_* parameters. It will download the specified refspec to prepare a new version of
            eric_nrf helm chart.
            This parameter is also used to clone the Jenkinsfile that will run the job""")
            defaultValue(' ')
        }
        stringParam {
            name('IHC_AUTO')
            description("""Allows you to select the eric_nrf Helm Chart
            image to use, in case you don't want to use the latest version or
            you need to use a custom version. Please specify full URL including Docker tag""")
            defaultValue('armdocker.rnd.ericsson.se/proj-5g-cicd-dev/adp-int-helm-chart-auto:experimental')
        }
        stringParam {
            name('CLOUD')
            description('To choose a different cluster from default one.<br><br>')
            defaultValue('hoff005')
        }
        stringParam {
            name('RELEASE')
            description('Create a new release version, default false')
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
                    }
                    branch('${GERRIT_BRANCH}')
                    extensions {
                        wipeOutWorkspace()
                    }
                }
                scriptPath('integration/eric_nrf/Jenkinsfile.prepare')
            }
        }
    }
}
