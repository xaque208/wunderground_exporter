pipeline {
  agent any

  triggers { pollSCM('*/5 * * * *') }

  environment {
    GOPATH = "${JENKINS_HOME}/workspace/go"
    PATH="${GOPATH}/bin:$PATH"
  }


  stages {
    stage('build') {

      steps {
        ws("${GOPATH}/src/github.com/xaque208/${JOB_NAME}") {
          checkout scm
          sh 'make clean'
          sh 'make build'
        }
      }
    }

    stage('publish') {
      when {
        branch 'master'
      }

      steps {
        ws("${GOPATH}/src/github.com/xaque208/${JOB_NAME}") {
          sh 'make publish'
        }
      }
    }
  }
}
