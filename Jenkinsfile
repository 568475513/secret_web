// 流水线文件，除注释表明处外勿擅自修改
// 任何问题，联系基础架构@Sparkli

def getGitBranchName() {
return scm.branches[0].name
}
def getRepositoryUrl() {
return scm.userRemoteConfigs[0].url
}

pipeline {
  environment {
    NAMESPACE = 'alive-platform' // 替换服务所在的项目名
    CONFIGFILE = '.env.docker' // 替换为配置文件相对路径+文件名，如config/.env.xxx
    REGISTRY = 'registry.xiaoe-tools.com'
    REGISTRYNAMESPACE = 'dev'
    // GOPROXY = 'https://goproxy.cn,direct'
    BRANCH = getGitBranchName()
    GITURL = getRepositoryUrl()
    GOPROXY='https://goproxy.cn,http://goproxy.xiaoe-tools.com,direct'
  }
  agent {
    node {
      label 'go'
    }
  }

  options {
    disableConcurrentBuilds()
    skipDefaultCheckout()
    buildDiscarder(
        logRotator(numToKeepStr: '10',daysToKeepStr: '7')
    )
    timeout(time: 1, unit: 'HOURS')
  }

  stages {
    stage('git pull') {
      agent none
      steps {
        container('go') {
            checkout([
            $class: 'GitSCM',
            branches: [[name: "${branch}"]],
            extensions: [[$class: 'CloneOption', depth: 1, noTags: false, reference: '', shallow: true, timeout: 10]],
            userRemoteConfigs: [[credentialsId: 'gitlab', url: "${giturl}"]]
            ])
        }
      }
    }

//     stage('单元测试') {
//           agent none
//           steps {
//             container('go') {
//               sh 'go test -work -timeout 3s ./...'
//             }
//           }
//     }

    stage('docker build && login && push') {
      agent none
      steps {
        container('go') {
          withCredentials([usernamePassword(credentialsId : 'registryhub' ,passwordVariable : 'registrypwd' ,usernameVariable : 'registryid' ,)]) {
            sh '''export COMMITID=$(git rev-parse --short HEAD)
export BUILDDATE=$(date "+%Y%m%d")
export SERVICE=$(basename -s .git `git config --get remote.origin.url`)
export SERVICE=${SERVICE//\\_/-}

docker build --network host -t $REGISTRY/$REGISTRYNAMESPACE/$SERVICE:$BUILDDATE-$COMMITID  .
docker login -u $registryid -p $registrypwd $REGISTRY

docker push $REGISTRY/$REGISTRYNAMESPACE/$SERVICE:$BUILDDATE-$COMMITID'''
          }
        }
      }
    }

    stage('k8s deploy') {
      agent none
      steps {
        container('go') {
        withCredentials([kubeconfigFile(credentialsId : 'kubeconfig' ,variable : 'KUBECONFIG' ,)]) {
          sh '''export COMMITID=$(git rev-parse --short HEAD)
export BUILDDATE=$(date "+%Y%m%d")
export SERVICE=$(basename -s .git `git config --get remote.origin.url`)
export BRANCH=${BRANCH//\\//-}
export BRANCH=${BRANCH//\\_/-}
export SERVICE=${SERVICE//\\_/-}
export CONFIGMAPHASH=$(md5sum $(echo $CONFIGFILE) | awk '{ print $1 }')

CMNAME=$SERVICE-$BRANCH-$CONFIGMAPHASH
CMFULLNAME=`kubectl get cm $(echo $CMNAME) -n $(echo $NAMESPACE) -o name --ignore-not-found=true`
if [ "$CMFULLNAME" != "configmap/${CMNAME}" ];then
    kubectl create cm $(echo $CMNAME) --from-file=$(echo $CONFIGFILE) --namespace=$(echo $NAMESPACE) ;
fi

envsubst < deploy.yaml | kubectl apply -f -'''
        }
       }
      }
    }

  }
}
