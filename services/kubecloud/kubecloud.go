package kubecloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"seaman/services"
	"time"

	"github.com/golang/glog"

	"seaman/keyword"
	"seaman/utils/config"
)

var (
	clusterMap = make(map[int64]map[string]string)
)

type KubecloudRsp struct {
	IsSuccess bool        `json:"IsSuccess"`
	Data      interface{} `json:"Data"`
}

type ClusterListRep struct {
	IsSuccess bool      `json:"IsSuccess"`
	Data      []Cluster `json:"Data"`
}

type ClusterRep struct {
	IsSuccess bool    `json:"IsSuccess"`
	Data      Cluster `json:"Data"`
}

type Cluster struct {
	ID                     int         `json:"id"`
	Name                   string      `json:"name"`
	ClusterID              string      `json:"cluster_id"`
	Tenant                 string      `json:"tenant"`
	Env                    string      `json:"env"`
	DisplayName            string      `json:"display_name"`
	Registry               string      `json:"registry"`
	ImagePullAddr          string      `json:"image_pull_addr"`
	DockerVersion          string      `json:"docker_version"`
	NetworkPlugin          string      `json:"network_plugin"`
	DomainSuffix           string      `json:"domain_suffix"`
	Certificate            string      `json:"certificate"`
	PrometheusAddr         string      `json:"prometheus_addr"`
	Status                 string      `json:"status"`
	Creator                string      `json:"creator"`
	KubeVersion            string      `json:"kube_version"`
	LoadbalancerDomainName string      `json:"loadbalancer_domain_name"`
	LoadbalancerIP         string      `json:"loadbalancer_ip"`
	LoadbalancerPort       string      `json:"loadbalancer_port"`
	KubeServiceAddress     string      `json:"kube_service_address"`
	KubePodSubnet          string      `json:"kube_pod_subnet"`
	TillerHost             string      `json:"tiller_host"`
	PromRuleIndex          int         `json:"prom_rule_index"`
	LabelPrefix            string      `json:"label_prefix"`
	Deleted                int         `json:"deleted"`
	DeleteAt               time.Time   `json:"delete_at"`
	Master                 interface{} `json:"master"`
	LbConfig               struct {
		DomainName string `json:"domain_name"`
		IP         string `json:"ip"`
		Port       string `json:"port"`
	} `json:"lb_config"`
	CreateAt         string `json:"create_at"`
	UpdateAt         string `json:"update_at"`
	ConfigRepo       string `json:"config_repo"`
	ConfigRepoBranch string `json:"config_repo_branch"`
	LastCommitId     string `json:"last_commit_id"`
}

type Kubecloud struct {
	base  string
	token string
}

func NewKubecloud() *Kubecloud {
	kubecloudConfig := config.GetConfig().Kubecloud
	return &Kubecloud{
		base:  kubecloudConfig.URL + "/zcloud/api/v3",
		token: kubecloudConfig.Token,
	}
}

func (k *Kubecloud) GetClusterInfo(cluster string, companyId int64) (*ClusterRep, error) {
	clusterId, err := k.getClusterId(companyId, cluster)
	if err != nil {
		glog.Errorf("get cluster info error: %s", err.Error())
		return nil, err
	}
	api := fmt.Sprintf("/clusters/%v", clusterId)
	respData, code, err := k.request("GET", api, nil, keyword.UserAdmin, companyId)
	if err != nil {
		glog.Errorf("get cluster info error: %s", err.Error())
		return nil, err
	}
	if code != http.StatusOK {
		glog.Errorf("get cluster info error: %s", string(respData))
		return nil, fmt.Errorf("failed to get cluster, status code %v", code)
	}
	var resp ClusterRep
	if err := json.Unmarshal(respData, &resp); err != nil {
		glog.Errorf("get cluster info error: %s", err.Error())
		return nil, err
	}
	return &resp, nil
}

func (k *Kubecloud) UpdateCluster(cluster string, companyId int64, body Cluster) error {
	clusterId, err := k.getClusterId(companyId, cluster)
	if err != nil {
		glog.Errorf("get cluster info error: %s", err.Error())
		return err
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}
	payload := bytes.NewBuffer(jsonData)
	api := fmt.Sprintf("/clusters/%v", clusterId)
	respData, code, err := k.request("PUT", api, payload, keyword.UserAdmin, companyId)
	if err != nil {
		glog.Errorf("get cluster info error: %s", err.Error())
		return err
	}
	if code != http.StatusOK {
		glog.Errorf("get cluster info error: %s", string(respData))
		return fmt.Errorf("failed to get cluster, status code %v", code)
	}
	return nil
}

func (k *Kubecloud) GetClusterList(companyId int64) ([]Cluster, error) {
	api := fmt.Sprintf("/clusters")
	respData, code, err := k.request("GET", api, nil, keyword.UserAdmin, companyId)
	if err != nil {
		glog.Errorf("get cluster list error: %s", err.Error())
		return nil, err
	}
	if code != http.StatusOK {
		glog.Errorf("get cluster list error: %s", string(respData))
		return nil, fmt.Errorf("failed to get cluster list, status code %v", code)
	}
	var resp ClusterListRep
	if err := json.Unmarshal(respData, &resp); err != nil {
		glog.Errorf("get cluster info error: %s", err.Error())
		return nil, err
	}
	return resp.Data, nil
}

func (k *Kubecloud) getClusterId(companyId int64, cluster string) (string, error) {
	if _, ok := clusterMap[companyId]; !ok {
		clusterMap[companyId] = make(map[string]string)
	}
	if _, ok := clusterMap[companyId][cluster]; !ok {
		clusters, err := k.GetClusterList(companyId)
		if err != nil {
			return "", err
		}
		for _, icluster := range clusters {
			clusterMap[companyId][icluster.Name] = icluster.ClusterID
		}
	}
	if _, ok := clusterMap[companyId][cluster]; ok {
		return clusterMap[companyId][cluster], nil
	}
	return "", fmt.Errorf("Can't find clusterId of cluster %v in tenant %v", cluster, companyId)
}

func (k *Kubecloud) request(method, api string, body io.Reader, user string, companyId int64) ([]byte, int, error) {
	url := k.base + api
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, -1, err
	}
	if companyId != 0 {
		cookie1 := &http.Cookie{Name: "zcloud_companyid", Value: fmt.Sprintf("%v", companyId), HttpOnly: true}
		req.AddCookie(cookie1)
	}
	if user == "" {
		user = "admin"
	}
	cookie2 := &http.Cookie{Name: "zcloud_username", Value: user, HttpOnly: true}
	req.AddCookie(cookie2)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+k.token)
	resp, err := services.HttpClient.Do(req)
	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, -1, err
	}
	code := resp.StatusCode
	return data, code, err
}
