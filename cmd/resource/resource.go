// Copyright 2020 angmas
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resource

import (
	"errors"
	"fmt"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	iotSvc := iot.New(req.Session)
	s3Svc := s3.New(req.Session)

	testKey := req.LogicalResourceID
	_, err := s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: currentModel.Bucket,
		Key:    &testKey,
		Body:   strings.NewReader("test"),
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			fmt.Printf("%v", aerr)
		}
		response := handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Bucket is not accessible",
		}
		return response, nil
	}
	_, err = s3Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: currentModel.Bucket,
		Key:    &testKey,
	})
	active := false
	if currentModel.Status != nil && strings.Compare(*currentModel.Status, "ACTIVE") == 0 {
		active = true
	}
	res, err := iotSvc.CreateKeysAndCertificate(&iot.CreateKeysAndCertificateInput{
		SetAsActive: &active,
	})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			fmt.Printf("%v", aerr)
		}
		response := handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         fmt.Sprintf("Failed: %s", aerr.Error()),
		}
		return response, nil
	}

	var key string
	key = fmt.Sprintf("%s.pem", *res.CertificateId)
	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: currentModel.Bucket,
		Key:    &key,
		Body:   strings.NewReader(*res.CertificatePem),
	})
	key = fmt.Sprintf("%s.key", *res.CertificateId)
	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: currentModel.Bucket,
		Key:    &key,
		Body:   strings.NewReader(*res.KeyPair.PrivateKey),
	})
	currentModel.Arn = res.CertificateArn
	currentModel.Id = res.CertificateId

	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Created certificate",
		ResourceModel:   currentModel,
	}
	return response, nil
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "Read complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: Read")
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)
	iotSvc := iot.New(req.Session)

	_, err := iotSvc.UpdateCertificate(&iot.UpdateCertificateInput{
		CertificateId: currentModel.Id,
		NewStatus:     currentModel.Status,
	})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			fmt.Printf("%v", aerr)
		}
		response := handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         aerr.Error(),
		}
		return response, nil

	}
	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Certificate updated successfully",
		ResourceModel:   currentModel,
	}
	return response, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)
	iotSvc := iot.New(req.Session)
	s3Svc := s3.New(req.Session)

	inactiveStatus := "INACTIVE"
	iotSvc.UpdateCertificate(&iot.UpdateCertificateInput{
		CertificateId: currentModel.Id,
		NewStatus:     &inactiveStatus,
	})

	var key string
	key = fmt.Sprintf("%s.pem", *currentModel.Id)

	_, err := s3Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: currentModel.Bucket,
		Key:    &key,
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			fmt.Printf("%v", aerr)
		}
		response := handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Bucket is not accessible",
		}
		return response, nil

	}

	key = fmt.Sprintf("%s.key", *currentModel.Id)

	_, err = s3Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: currentModel.Bucket,
		Key:    &key,
	})
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			fmt.Printf("%v", aerr)
		}
		response := handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Bucket is not accessible",
		}
		return response, nil

	}

	_, err = iotSvc.DeleteCertificate(&iot.DeleteCertificateInput{
		CertificateId: currentModel.Id,
	})

	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			fmt.Printf("%v", aerr)
		}
		response := handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         fmt.Sprintf("Failed deleting the certificate: %s", aerr.Error()),
		}
		return response, nil

	}

	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Certificate delete successfully",
		ResourceModel:   currentModel,
	}
	return response, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// Add your code here:
	// * Make API calls (use req.Session)
	// * Mutate the model
	// * Check/set any callback context (req.CallbackContext / response.CallbackContext)

	/*
	   // Construct a new handler.ProgressEvent and return it
	   response := handler.ProgressEvent{
	       OperationStatus: handler.Success,
	       Message: "List complete",
	       ResourceModel: currentModel,
	   }

	   return response, nil
	*/

	// Not implemented, return an empty handler.ProgressEvent
	// and an error
	return handler.ProgressEvent{}, errors.New("Not implemented: List")
}
