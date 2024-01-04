def versionTag
pipeline {
  environment {
    gitBranch = "${env.BRANCH_NAME}"
    gitURL = "git@github.com:Memphisdev/memphis-functions.go.git"
    repoUrlPrefix = "memphisos"
  }

  agent {
    label 'memphis-jenkins-small-fleet-agent'
  }

  stages {
    stage ('Connect GIT repository') {
      steps {
        git credentialsId: 'main-github', url: "git@github.com:Memphisdev/memphis-functions.go.git", branch: "${env.gitBranch}" 
      }
    }

    stage('Define version - BETA') {
      when {branch 'change-jenkins-agent'}
      steps {
        script {
          versionTag = readFile('./version-beta.conf')
        }
      }
    }
    stage('Define version - LATEST') {
      when {branch 'latest'}
      steps {
        script {
          versionTag = readFile('./version.conf')
        }
      }
    }
        

    stage('Install GoLang') {
      steps {
        sh """
          wget -q https://go.dev/dl/go1.20.12.linux-amd64.tar.gz
          sudo  tar -C /usr/local -xzf go1.20.12.linux-amd64.tar.gz
        """
      }
    }

    stage('Deploy GO Functions SDK') {
      steps {
        withCredentials([sshUserPrivateKey(keyFileVariable:'check',credentialsId: 'main-github')]) {
          sh """
            git tag v$versionTag
            GIT_SSH_COMMAND='ssh -i $check' git push origin v$versionTag
          """
        }
        sh """
          GOPROXY=proxy.golang.org /usr/local/go/bin/go list -m github.com/memphisdev/memphis-functions.go@v$versionTag
        """
      }
    }

    stage('Checkout to version branch'){
      when {branch 'latest'}
      steps {
        withCredentials([sshUserPrivateKey(keyFileVariable:'check',credentialsId: 'main-github')]) {
          sh "git reset --hard origin/latest"
          sh "GIT_SSH_COMMAND='ssh -i $check'  git checkout -b $versionTag"
          sh "GIT_SSH_COMMAND='ssh -i $check' git push --set-upstream origin $versionTag"
        }
      }
    }

  }
   
  post {
    always {
      cleanWs()
    }
    success {
      notifySuccessful()
    }
    failure {
      notifyFailed()
    }
  }
}
def notifySuccessful() {
    emailext (
        subject: "SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
        body: """SUCCESSFUL: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]':
        Check console output and connection attributes at ${env.BUILD_URL}""",
        recipientProviders: [requestor()]
    )
}
def notifyFailed() {
    emailext (
        subject: "FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'",
        body: """FAILED: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]':
        Check console output at ${env.BUILD_URL}""",
        recipientProviders: [requestor()]
    )
}
  
