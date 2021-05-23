pipeline {
    agent {
      label "master"
    }
    parameters{
        string(name: "VERSION")
        string(name: "GIT_USERNAME")
        string(name: "GIT_PASSWORD")
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
                sh "git push https://${params.GIT_USERNAME}:${params.GIT_PASSWORD}@github.com/eaneto/notify --tags"
            }
        }
    }
}
