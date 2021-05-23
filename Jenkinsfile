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
        stage('Test'){
            steps {
                sh "echo ${params.VERSION}"
            }
        }
    }
}
