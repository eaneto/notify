pipeline {
    agent {
      label "master"
    }
    parameters{
        string(name: "VERSION")
    }
    stages {
        stage('Clone') {
            steps {
                git branch: "main", url:"https://github.com/eaneto/notify"
            }
        }
        stage('Release'){
            steps {
                sh "echo ${params.VERSION}"
                sh "git tag ${params.VERSION}"
                sh "git push origin ${params.VERSION}"
            }
        }
    }
}
