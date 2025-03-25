/*
Copyright 2018 Google LLC

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
package integration

import (
	"fmt"
	"os/exec"

	integration_util "github.com/soy-kyle/kritis/pkg/kritis/integration_util"
	v1 "k8s.io/api/core/v1"
)

func podLogs(pod string, ns *v1.Namespace) string {
	cmd := exec.Command("kubectl", "logs", pod, "-n", ns.Name)
	out, err := integration_util.RunCmdOut(cmd)
	if err != nil {
		return fmt.Sprintf("unable to get pod logs for %q in %s: %v", pod, ns.Name, err)
	}
	return string(out)
}

func kritisLogs(ns *v1.Namespace) string {
	cmd := exec.Command("kubectl", "logs", "-l", "label=kritis-validation-hook", "-n", ns.Name, "--tail=100")
	out, err := integration_util.RunCmdOut(cmd)
	if err != nil {
		return fmt.Sprintf("failed to get kritis-validation-hook logs: %v", err)
	}

	cmd = exec.Command("kubectl", "get", "imagesecuritypolicy", "--namespace", ns.Name)
	out2, err := integration_util.RunCmdOut(cmd)
	if err != nil {
		return fmt.Sprintf("failed to get isp: %v", err)
	}
	return string(out) + "\n" + string(out2)
}
