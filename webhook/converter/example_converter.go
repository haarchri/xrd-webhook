package converter

import (
	"fmt"
	"strings"

	"k8s.io/klog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func convertExampleCRD(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	klog.V(2).Info("converting crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}

	switch Object.GetAPIVersion() {
	case "conversion.haarchri.io/v1beta1":
		switch toVersion {
		case "conversion.haarchri.io/v1":
			hostPort, ok, _ := unstructured.NestedString(convertedObject.Object, "spec", "hostPort")
			if ok {
				delete(convertedObject.Object, "spec")
				parts := strings.Split(hostPort, ":")
				if len(parts) != 2 {
					return nil, statusErrorWithMessage("invalid hostPort value `%v`", hostPort)
				}
				host := parts[0]
				port := parts[1]
				unstructured.SetNestedField(convertedObject.Object, host, "spec", "host")
				unstructured.SetNestedField(convertedObject.Object, port, "spec", "port")
			}
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}
	case "conversion.haarchri.io/v1":
		switch toVersion {
		case "conversion.haarchri.io/v1beta1":
			host, hasHost, _ := unstructured.NestedString(convertedObject.Object, "spec", "host")
			port, hasPort, _ := unstructured.NestedString(convertedObject.Object, "spec", "port")
			if hasHost || hasPort {
				if !hasHost {
					host = ""
				}
				if !hasPort {
					port = ""
				}
				hostPort := fmt.Sprintf("%s:%s", host, port)
				unstructured.SetNestedField(convertedObject.Object, hostPort, "spec", "hostPort")
			}
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}
	default:
		return nil, statusErrorWithMessage("unexpected conversion version %q", fromVersion)
	}
	return convertedObject, statusSucceed()
}
