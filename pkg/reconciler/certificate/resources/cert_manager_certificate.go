/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resources

import (
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/net-certmanager/pkg/reconciler/certificate/config"
	"knative.dev/networking/pkg/apis/networking"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/kmeta"
)

// MakeCertManagerCertificate creates a Cert-Manager `Certificate` for requesting a SSL certificate.
func MakeCertManagerCertificate(cmConfig *config.CertManagerConfig, knCert *v1alpha1.Certificate) *cmv1.Certificate {
	var commonName string
	if len(knCert.Spec.DNSNames) > 0 {
		commonName = knCert.Spec.DNSNames[0]
	}
	cert := &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:            knCert.Name,
			Namespace:       knCert.Namespace,
			OwnerReferences: []metav1.OwnerReference{*kmeta.NewControllerRef(knCert)},
			Annotations:     knCert.GetAnnotations(),
			Labels:          knCert.GetLabels(),
		},
		Spec: cmv1.CertificateSpec{
			CommonName: commonName,
			SecretName: knCert.Spec.SecretName,
			DNSNames:   knCert.Spec.DNSNames,
			IssuerRef:  *cmConfig.IssuerRef,
			SecretTemplate: &cmv1.CertificateSecretTemplate{
				Labels: map[string]string{
					networking.CertificateUIDLabelKey: string(knCert.GetUID()),
				}},
		},
	}
	return cert
}

// GetReadyCondition gets the ready condition of a Cert-Manager `Certificate`.
func GetReadyCondition(cmCert *cmv1.Certificate) *cmv1.CertificateCondition {
	for _, cond := range cmCert.Status.Conditions {
		if cond.Type == cmv1.CertificateConditionReady {
			return &cond
		}
	}
	return nil
}
