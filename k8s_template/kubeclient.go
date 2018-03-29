package kubeclient

import (
	"blockchain-mgr/httpclient"
	"blockchain-mgr/httpclient/encryption"
	"blockchain-mgr/logger"
	"blockchain-mgr/util"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"fmt"
)

var GKubeClient *KubeClient = nil

type KubeClient struct {
	ApiServer       string
	CerFilePath     string
	KeyFilePath     string
	PemFilePasswd   string
	RootKey         string
	CommonSharedKey string
}

type KubeResource struct {
	Kind		string `json:"kind"`
	Metadata 	KubeResourceMetadata `json:"metadata"`
}

type KubeResourceMetadata struct {
	Name string `json:"name"`
}
/*
 @ Description: Initialize HTTPS client and ApiServer
 @ param apiServer - CFE APIServer IP-Address:port
 @ param cerFilePath - Certificate file path
 @ param keyFilePath - Key file path
 @ param pemFilePasswd - password
 @ param rootKey - Root Key file
 @ param commonSharedKey - Common shared key file
 @ return success: nil
*/
func New(apiServer, cerFilePath, keyFilePath, pemFilePasswd, rootKey, commonSharedKey string) KubeClient {
	return KubeClient{
		ApiServer:       apiServer,
		CerFilePath:     cerFilePath,
		KeyFilePath:     keyFilePath,
		PemFilePasswd:   pemFilePasswd,
		RootKey:         rootKey,
		CommonSharedKey: commonSharedKey,
	}
}

// InitHttp is private function used to init https client.
func (kc *KubeClient) InitHttp() (err error) {
	if http_client.IsInitHttpsClient() {
		return nil
	}
	err = encryption.EncryptionInit()
	if err != nil {
		logger.Info("EncryptionInit: ", err)
		return err
	}
	err = http_client.HttpsClientInit(kc.KeyFilePath, kc.PemFilePasswd, kc.CerFilePath)
	if err != nil {
		logger.Info("HttpsClientInit: ", err, kc.KeyFilePath, kc.PemFilePasswd, kc.CerFilePath)
		return err
	}
	return nil
}

func (kc *KubeClient) InitUrl(namespace string, rsType string) (string, error) {
	var path string
	err := kc.InitHttp()
	if err != nil {
		logger.Info("InitUrl: InitHttp client error.")
		return "", err
	}
	if namespace == "" {
		namespace = "default"
	}
	switch rsType {
	case "Deployment":
		path = "/apis/extensions/v1beta1/namespaces/" + namespace + "/deployments"
	case "Service":
		path = "/api/v1/namespaces/" + namespace + "/services"
	case "Ingress":
		path = "/apis/extensions/v1beta1/namespaces/" + namespace + "/ingresses"
	case "ConfigMap":
		path = "/api/v1/namespaces/" + namespace + "/configmaps"
	case "Secret":
		path = "/api/v1/namespaces/" + namespace + "/secrets"
	case "Job":
		path = "/apis/batch/v1/namespaces/" + namespace + "/jobs"
	case "LimitRange":
		path = "/api/v1/namespaces/" + namespace + "/limitranges"
	case "Node":
		path = "/api/v1/namespaces/" + namespace + "/nodes"
	case "Namespace":
		path = "/api/v1/namespaces"
	case "ProcessLifecycle":
		path = "/apis/paas/v1alpha1/namespaces/" + namespace + "/processlifecycles"
	case "StatefulSet":
		path = "/apis/apps/v1beta1/namespaces/" + namespace + "/statefulsets"
	case "DaemonSet":
		path = "/apis/extensions/v1beta1/namespaces/" + namespace + "/daemonsets"
	case "Endpoints":
		path = "/api/v1/namespaces/" + namespace + "/endpoints"
	case "ReplicaSet":
		path = "/apis/extensions/v1beta1/namespaces/" + namespace + "/replicasets"
	case "Pod":
		path = "/api/v1/namespaces/" + namespace + "/pods"
	case "PersistentVolume":
		path = "/api/v1/namespaces/" + namespace + "/persistentvolumes"
	case "PersistentVolumeClaim":
		path = "/api/v1/namespaces/" + namespace + "/persistentvolumeclaims"
	default:
		err = errors.New("unsupported k8s yaml type: " + rsType)
		logger.Info("CreateKubeResource: ", err.Error())
		return "", err
	}
	logger.Info("path=", path)
	return path, err
}

func (kc *KubeClient) InitUrlWithBody(namespace string, body string) (string, error) {
	var path string
	err := kc.InitHttp()
	if err != nil {
		logger.Info("InitUrl: InitHttp client error.")
		return "", err
	}
	//检验参数
	if namespace == "" {
		namespace = "default"
	}
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return "", err
	}
	rsType := js.Get("kind").MustString()
	apiVersion := js.Get("apiVersion").MustString()

	switch rsType {
	case "Deployment":
		path = "/apis/" + apiVersion + "/namespaces/" + namespace + "/deployments"
	case "Service":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/services"
	case "Ingress":
		path = "/apis/" + apiVersion + "/namespaces/" + namespace + "/ingresses"
	case "ConfigMap":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/configmaps"
	case "Secret":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/secrets"
	case "Job":
		path = "/apis/" + apiVersion + "/namespaces/" + namespace + "/jobs"
	case "LimitRange":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/limitranges"
	case "Node":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/nodes"
	case "Namespace":
		path = "/api/" + apiVersion + "/namespaces"
	case "ProcessLifecycle":
		path = "/apis/" + apiVersion + "/namespaces/" + namespace + "/processlifecycles"
	case "StatefulSet":
		path = "/apis/" + apiVersion + "/namespaces/" + namespace + "/statefulsets"
	case "DaemonSet":
		path = "/apis/" + apiVersion + "/namespaces/" + namespace + "/daemonsets"
	case "Endpoints":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/endpoints"
	case "ReplicaSet":
		path = "/apis/" + apiVersion + "/namespaces" + namespace + "/replicasets"
	case "Pod":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/pods"
	case "PersistentVolume":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/persistentvolumes"
	case "PersistentVolumeClaim":
		path = "/api/" + apiVersion + "/namespaces/" + namespace + "/persistentvolumeclaims"
	default:
		err = errors.New("unsupported k8s yaml type: " + rsType)
		logger.Info("CreateKubeResource: ", err.Error())
		return "", err
	}
	logger.Info("path=", path)
	return path, err
}

func setHeaders(token, clusterID string) map[string]string {
	headers := make(map[string]string)

	if token != "" {
		headers["X-Auth-Token"] = token
	}
	if clusterID != "" {
		headers["X-Cluster-ID"] = clusterID
	}
	return headers
}

func setParams() map[string]string {
	params := make(map[string]string)

	return params
}

func (kc *KubeClient) CreateKubeResource(namespace, body, token, clusterID string) error {
	var path string

	path, err := kc.InitUrlWithBody(namespace, body)
	if err != nil {
		logger.Info("Failed to init url for https, error: %v", err)
		return err
	}
	headers := setHeaders(token, clusterID)
	params := setParams()
	count := 3
	for i := 0; i < count; i++ {
		resp, code, err := http_client.PostReq(kc.ApiServer, path, []byte(body), headers, params)
		if err != nil {
			logger.Info("CreateKubeResource: ", string(resp), err.Error())
			if code >= 500 && i < count-1 {
				time.Sleep(2 * time.Second)
				continue
			}
			return err
		} else {
			break
		}
	}
	return nil
}

func (kc *KubeClient) UpdateKubeResource(namespace, body, resName, token, clusterID string) error {
	var path string

	path, err := kc.InitUrlWithBody(namespace, body)
	if err != nil {
		fmt.Printf("Failed to init url for https, error: %v", err)
		return err
	}
	path = path + "/" + resName
	headers := setHeaders(token, clusterID)
	params := setParams()
	count := 3
	for i := 0; i < count; i++ {
		resp, err := http_client.PutReq(kc.ApiServer, path, []byte(body), headers, params)
		if err != nil {
			logger.Error("UpdateKubeResource: ", string(resp), err.Error())
			return err
		} else {
			break
		}
	}
	return nil
}

func (kc *KubeClient) UpdateKubeResourceByMethod(namespace,  rsType, resourceName, body, token, clusterID, httpMethod string) error {
	var path string
	var resp []byte
	rsPath, err := kc.InitUrl(namespace, rsType)
	if err != nil {
		return err
	}
	headers := setHeaders(token, clusterID)
	params := setParams()

	path = rsPath + "/" + resourceName

	switch httpMethod {
	case util.HTTP_PUT:
		resp, err = http_client.PutReq(kc.ApiServer, path, []byte(body), headers, params);
	case util.HTTP_PATCH:
		resp, err = http_client.PatchReq(kc.ApiServer, path, []byte(body), headers, params);
	default:
		err = fmt.Errorf("not support http method type: %s", httpMethod)
	}
	//logger.Info("request info, api addr:", kc.ApiServer)
	//logger.Info("request info, path:", path)
	//logger.Info("request info, body:", body)
	//logger.Info("request info, headers:", headers)
	//logger.Info("request info, params:", params)
	if err != nil {
		err = errors.New("UpdateKubeResource failed: " + string(resp) + err.Error())
		return err
	}
	return nil
}

func (kc *KubeClient) DeleteKubeResource(namespace, rsType, resourceName, token, clusterID string) error {
	var path string

	rsPath, err := kc.InitUrl(namespace, rsType)
	if err != nil {
		return err
	}
	headers := setHeaders(token, clusterID)
	params := setParams()

	path = rsPath + "/" + resourceName
	logger.Info("path=", path)

	if resp, err := http_client.DeleteReq(kc.ApiServer, path, []byte(""), headers, params); err != nil {
		logger.Info("DeleteKubeResource: ", string(resp), err.Error())
		// 如果错误是NOT_FOUND，则返回nil
		if err.Error() == http_client.NOT_FOUND {
			return nil
		}
		return err
	}
	return nil
}

func (kc *KubeClient) IsResourceNameRunning(namespace, resourceName, resourceType, token, clusterID string, replica int) (isRunning bool, err error) {
	//use statefulSetName plus replica to get pod name
	var count = 0
	isRunning = true
	resp, err := kc.QueryResourceExist(namespace, "", token, clusterID, util.K8S_POD)

	//解析json
	respJson, _ := simplejson.NewJson([]byte(resp))
	if respJson == nil {
		logger.Info("QueryPod result and respose result is nil")
		return false, err
	}
	pods, _ := respJson.Get("items").Array()

	//获取status phase
	for _, pod := range pods {
		//就在这里对pod进行类型判断
		podData, _ := pod.(map[string]interface{})

		//获取Pod的name值
		podMetadata := podData["metadata"].(map[string]interface{})
		podName := podMetadata["name"]

		//根据`statefulSetName`过滤pod
		if strings.Contains(podName.(string), resourceName) {
			if resourceType == util.RESOURCE_DEPLOYMENT {
				if !strings.HasPrefix(podName.(string),resourceName+"-"){
					continue
				}
			}else if resourceType == util.RESOURCE_STATEFULSET {
				//判断podName是否为资源对应的pod
				podNameSuffix := strings.TrimPrefix(podName.(string), resourceName+"-")
				if _, err := strconv.Atoi(podNameSuffix); err != nil {
					continue
				}
			}else {
				continue
			}

			count++
			//获取Pod的phase值
			podStatus := podData["status"].(map[string]interface{})
			podphase := podStatus["phase"]

			if podphase != "Running" {
				isRunning = false
				return isRunning, err
			}
		}

	}
	if count != replica {
		isRunning = false
		return isRunning, err
	}

	return isRunning, err
}

func (kc *KubeClient) QueryResourceExist(namespace, resName, token, clusterID string, resType util.K8sResType) (resp []byte, err error) {
	rsPath, err := kc.InitUrl(namespace, string(resType))
	if err != nil {
		logger.Info("Failed to init url for https, error: %v", err)
		return
	}
	path := rsPath
	if resName != "" {
		path = path + "/" + resName
	}
	headers := setHeaders(token, clusterID)
	params := setParams()

	if resp, err = http_client.GetReq(kc.ApiServer, path, headers, params); err != nil {
		if err.Error() == http_client.NOT_FOUND {
			logger.Info(resType, resName, http.StatusNotFound, err.Error())
			return
		}
		logger.Info(resType, resName, http.StatusInternalServerError, err.Error())
		return
	}
	return
}

// later will add secrets
func isSupportUpgrade(t string) bool {
	switch t {
	case "Deployment":
		return true
	case "StatefulSet":
		return true
	default:
		return false
	}
}

func (kc *KubeClient) UpgradeKubeResource(namespace, body, token, clusterID string) error {
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return err
	}

	rsType := js.Get("kind").MustString()
	if !isSupportUpgrade(rsType) {
		logger.Info("don't support %s upgrade.", rsType)
		return nil
	}

	resouceName := js.Get("metadata").Get("name").MustString()
	rsPath, err := kc.InitUrlWithBody(namespace, body)
	if err != nil {
		return err
	}
	headers := setHeaders(token, clusterID)
	params := setParams()
	path := rsPath + "/" + resouceName
	logger.Info("upgrade path=", path)
	if resp, err := http_client.PatchReq(kc.ApiServer, path, []byte(body), headers, params); err != nil {
		logger.Info("UpgradeKubeResource: ", string(resp), err.Error())
		return err
	}
	return nil
}

// PatchNodeLabel patch label on nodes
func (kc *KubeClient) PatchNodeLabel(namespace, nodename, body, token, clusterID string) error {
	var path = ""
	nodePath, err := kc.InitUrl(namespace, "Node")
	if err != nil {
		return err
	}
	headers := setHeaders(token, clusterID)
	params := setParams()
	path = nodePath + "/" + nodename

	if resp, err := http_client.PatchReq(kc.ApiServer, path, []byte(body), headers, params); err != nil {
		logger.Info("PatchNodeLabel: ", string(resp), err.Error())
		return err
	}
	return nil
}

//Patch annotations on nodes
func (kc *KubeClient) PatchNodeAnnotations(namespace, nodename, body, token, clusterID string) error {
	var path = ""
	nodePath, err := kc.InitUrl(namespace, "Node")
	if err != nil {
		return err
	}
	headers := setHeaders(token, clusterID)
	params := setParams()
	path = nodePath + "/" + nodename

	if resp, err := http_client.PatchReq(kc.ApiServer, path, []byte(body), headers, params); err != nil {
		logger.Info("PatchNodeAnnotations: ", string(resp), err.Error())
		return err
	}
	return nil
}
