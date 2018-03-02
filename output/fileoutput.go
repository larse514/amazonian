package output

import "io/ioutil"

//Output struct
type Output struct {
	FileName    string
	ServiceName string
	ClusterName string
	ServiceURL  string
	VPCId       string
	VPCName     string
}

//WriteOutputFile writes output to a file
func WriteOutputFile(output Output) error {
	content := []byte("ServiceName=" + output.ServiceName + "\nClusterName=" + output.ClusterName + "\nURL=" + output.ServiceURL + "\nVPCId=" + output.VPCId + "VPCName=" + output.VPCName)
	return ioutil.WriteFile(output.FileName, content, 0644)
}
