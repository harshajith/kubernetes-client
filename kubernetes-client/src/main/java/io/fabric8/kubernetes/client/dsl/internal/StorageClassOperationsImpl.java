/**
 * Copyright (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package io.fabric8.kubernetes.client.dsl.internal;

import io.fabric8.kubernetes.api.model.DoneableStorageClass;
import io.fabric8.kubernetes.api.model.StorageClass;
import io.fabric8.kubernetes.api.model.StorageClassList;
import io.fabric8.kubernetes.client.Config;
import io.fabric8.kubernetes.client.dsl.Resource;
import io.fabric8.kubernetes.client.dsl.base.HasMetadataOperation;
import okhttp3.OkHttpClient;

import java.util.Map;
import java.util.TreeMap;

public class StorageClassOperationsImpl extends HasMetadataOperation<StorageClass, StorageClassList, DoneableStorageClass, Resource<StorageClass, DoneableStorageClass>> {
  public StorageClassOperationsImpl(OkHttpClient client, Config config) {
    this(client, config, null, null, null, true, null, null, false, -1, new TreeMap<String, String>(), new TreeMap<String, String>(), new TreeMap<String, String[]>(), new TreeMap<String, String[]>(), new TreeMap<String, String>());
  }

  public StorageClassOperationsImpl(OkHttpClient client, Config config, String apiVersion, String namespace, String name, Boolean cascading, StorageClass item, String resourceVersion, Boolean reloadingFromServer, long gracePeriodSeconds, Map<String, String> labels, Map<String, String> labelsNot, Map<String, String[]> labelsIn, Map<String, String[]> labelsNotIn, Map<String, String> fields) {
    super(client, config, "storage.k8s.io", apiVersion, "storageclasses", namespace, name, cascading, item, resourceVersion, reloadingFromServer, gracePeriodSeconds, labels, labelsNot, labelsIn, labelsNotIn, fields);
  }

  @Override
  public boolean isResourceNamespaced() {
    return false;
  }
}
