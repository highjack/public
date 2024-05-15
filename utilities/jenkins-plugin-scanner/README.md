# Jenkins plugin scanner

## WTF
Automate download of jenkins plugin and scan them with trivy for vulnerabilties


## HTF
- Add list of plugins to jenkins_plugins.txt
- docker build -t jenkins-downloader .
- docker run jenkins-downloader
- For now you will have to use docker cp to grab the files will fix this later
