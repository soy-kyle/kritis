# Copyright 2018 Google, Inc. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM gcr.io/gcp-runtimes/ubuntu_16_0_4 as runtime_deps

RUN apt-get install --no-install-recommends --no-install-suggests -y \
        git

ENV HELM_VERSION v2.8.1
RUN curl -LO https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    tar -zxvf helm-${HELM_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/local/bin/helm

# ENV HELM_HOME /builder/home
# RUN mkdir -p ${HELM_HOME} && \
#     helm init --client-only --home ${HELM_HOME} && \
#     helm plugin install https://github.com/nouney/helm-gcs --version 0.2.0 --home ${HELM_HOME}

WORKDIR /go/src/github.com/soy-kyle/kritis

COPY . .

CMD ["./hack/helm-release.sh"]
