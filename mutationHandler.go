package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

func mutationHandler(c *gin.Context) {
	var admissionResponse *v1beta1.AdmissionResponse
	admissionReview := v1beta1.AdmissionReview{}
	var err error
	err = c.ShouldBindJSON(&admissionReview)
	if err != nil {
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
	}

	req := admissionReview.Request
	var pod corev1.Pod

	if err := json.Unmarshal(req.Object.Raw, &pod); err != nil {
		glog.Errorf("Could not unmarshal raw object: %v", err)
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
		fmt.Println("error2")
		c.AbortWithStatusJSON(http.StatusInternalServerError,admissionResponse)
		return
	}

	patchBytes, err := createPatch(&pod)
	if err != nil {
		admissionResponse = &v1beta1.AdmissionResponse{
			Result: &metav1.Status{
				Message: err.Error(),
			},
		}
		fmt.Println("error3")
		c.AbortWithStatusJSON(http.StatusInternalServerError,admissionResponse)
		return
	}

	fmt.Println("continue")
	fmt.Println(string(patchBytes))
	admissionResponse = &v1beta1.AdmissionResponse{
		Allowed: true,
		Patch:   patchBytes,
		PatchType: func() *v1beta1.PatchType {
			pt := v1beta1.PatchTypeJSONPatch
			return &pt
		}(),
	}

	admissionReview.Response = admissionResponse
	admissionReview.Response.UID = admissionReview.Request.UID

	result,_  := json.Marshal(admissionReview)
	fmt.Println(string(result))
	c.JSON(http.StatusOK,admissionReview)
}

func createPatch(pod *corev1.Pod) ([]byte, error) {
	var patches []patchOperation
	for idx, _ := range pod.Spec.Containers {
		// 如果是 corev1.EnvVar 要加 /- ，如果不是 不需要加
		patches = append(patches, addEnv(fmt.Sprintf("/spec/containers/%d/env", idx))...)
	}
	return json.Marshal(patches)
}

func addEnv(basePath string) (patch []patchOperation) {
	patch = []patchOperation{
		{Op: "add", Path: basePath, Value: []corev1.EnvVar{{Name: "CLUSTER_NAME", Value: "aks-test-01", ValueFrom: nil}}},
		{Op: "add", Path: basePath + "/-", Value: corev1.EnvVar{Name: "TZ", Value: "Asia/Shanghai", ValueFrom: nil}},
	}
	return patch
}
