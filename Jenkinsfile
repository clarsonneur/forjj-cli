pipeline {
    agent any

    stages {
        stage('Cleanup') {
            when { branch 'master' }
            steps {
                cleanWs()
            }
        }
        stage('Tests') {
            steps {
                sh('''source ./build-env.sh
                build.sh''')
            }
        }
    }
}
