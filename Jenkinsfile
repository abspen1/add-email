pipeline {
    agent any
    stages {
        stage('Run unit tets') {
            steps {
                sh 'echo "I should be running tests.."'
            }
        }
        stage('copy files to pi') {
            steps {
                sshagent(credentials : ['gt-ssh']) {
                    sh 'if [ -d "/var/jenkins_home/.ssh/known_hosts" ]; then rm -Rf /var/jenkins_home/.ssh/known_hosts; fi'
                    sh """
                        ssh -o StrictHostKeyChecking=no pi@192.168.1.30 \
                        'if [ -d "source" ]; then rm -Rf source; fi'
                    """
                    sh 'scp -o StrictHostKeyChecking=no -r source/ pi@192.168.1.30:~/'
                }
            }
        }
        stage('ssh into pi and build/deploy function') {
            steps {
                sshagent(credentials : ['gt-ssh']) {                    
                    sh """
                        ssh -o StrictHostKeyChecking=no pi@192.168.1.30 \
                        'cd source/ \
                        && faas template pull https://github.com/openfaas-incubator/golang-http-template \
                        && cat ~/faas_pass.txt | faas-cli login --password-stdin -g 127.0.0.1:31375 \
                        && faas-cli up --build-arg GO111MODULE=on -f add-email.yml'
                    """
                }
            }
        }
        stage('Check k3 cluster pods') {
            steps {
                withCredentials([file(credentialsId: 'config', variable: 'config')]) {
                    sh 'curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"'
                    sh 'chmod +x ./kubectl'
                    sh './kubectl version --client'
                    
                    sh "cat \$config >> config"
                    sh "export KUBECONFIG=config"
                    
                    sh "./kubectl get pods -n openfaas-fn"
                }
            }
        }
    }
}
