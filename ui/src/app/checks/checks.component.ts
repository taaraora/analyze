import { Component, OnInit } from '@angular/core';
import { HttpClient }        from "@angular/common/http";
import { map, catchError }   from "rxjs/operators";
import { Observable, of }    from "rxjs";

@Component({
  selector: 'app-checks',
  templateUrl: './checks.component.html',
  styleUrls: ['./checks.component.scss']
})
export class ChecksComponent implements OnInit {

  checks$: Observable<any[]>;

  constructor(
    private http: HttpClient
  ) {
  }

  ngOnInit() {
    const apiV1Check = 'http://ec2-54-241-150-199.us-west-1.compute.amazonaws.com:30759/api/v1/check';
    const mock = [
      {
        "checkStatus": "RED",
        "completedAt": "2018-11-14T11:37:39.287Z",
        "description": "{\"Nodes\":[{\"NodeName\":\"ip-172-20-1-21.us-west-1.compute.internal\",\"Pods\":[{\"PodName\":\"dmts-es-1-elasticsearch-client-848f4d5db6-bvksx\",\"Containers\":[{\"ContainerName\":\"elasticsearch\",\"ContainerImage\":\"docker.elastic.co/elasticsearch/elasticsearch-oss:6.4.2\",\"Requests\":{\"RAM\":\"512Mi\",\"CPU\":\"25m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"1\"}}]},{\"PodName\":\"dmts-es-1-elasticsearch-master-0\",\"Containers\":[{\"ContainerName\":\"elasticsearch\",\"ContainerImage\":\"docker.elastic.co/elasticsearch/elasticsearch-oss:6.4.2\",\"Requests\":{\"RAM\":\"512Mi\",\"CPU\":\"25m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"1\"}}]},{\"PodName\":\"winning-otter-analyze-6b5d658d76-6jlzb\",\"Containers\":[{\"ContainerName\":\"analyze\",\"ContainerImage\":\"supergiant/analyze:latest\",\"Requests\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"}}]},{\"PodName\":\"kube-proxy-ip-172-20-1-21.us-west-1.compute.internal\",\"Containers\":[{\"ContainerName\":\"kube-proxy\",\"ContainerImage\":\"gcr.io/google_containers/hyperkube:v1.11.1\",\"Requests\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"}}]}]},{\"NodeName\":\"ip-172-20-1-242.us-west-1.compute.internal\",\"Pods\":[{\"PodName\":\"dmts-es-1-elasticsearch-client-848f4d5db6-xll2f\",\"Containers\":[{\"ContainerName\":\"elasticsearch\",\"ContainerImage\":\"docker.elastic.co/elasticsearch/elasticsearch-oss:6.4.2\",\"Requests\":{\"RAM\":\"512Mi\",\"CPU\":\"25m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"1\"}}]},{\"PodName\":\"dmts-es-1-elasticsearch-data-0\",\"Containers\":[{\"ContainerName\":\"elasticsearch\",\"ContainerImage\":\"docker.elastic.co/elasticsearch/elasticsearch-oss:6.4.2\",\"Requests\":{\"RAM\":\"1536Mi\",\"CPU\":\"25m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"1\"}}]},{\"PodName\":\"dmts-es-1-elasticsearch-master-1\",\"Containers\":[{\"ContainerName\":\"elasticsearch\",\"ContainerImage\":\"docker.elastic.co/elasticsearch/elasticsearch-oss:6.4.2\",\"Requests\":{\"RAM\":\"512Mi\",\"CPU\":\"25m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"1\"}}]},{\"PodName\":\"kube-proxy-ip-172-20-1-242.us-west-1.compute.internal\",\"Containers\":[{\"ContainerName\":\"kube-proxy\",\"ContainerImage\":\"gcr.io/google_containers/hyperkube:v1.11.1\",\"Requests\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"}}]}]},{\"NodeName\":\"ip-172-20-1-44.us-west-1.compute.internal\",\"Pods\":[{\"PodName\":\"dmts-es-1-elasticsearch-data-1\",\"Containers\":[{\"ContainerName\":\"elasticsearch\",\"ContainerImage\":\"docker.elastic.co/elasticsearch/elasticsearch-oss:6.4.2\",\"Requests\":{\"RAM\":\"1536Mi\",\"CPU\":\"25m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"1\"}}]},{\"PodName\":\"dmts-es-1-elasticsearch-master-2\",\"Containers\":[{\"ContainerName\":\"elasticsearch\",\"ContainerImage\":\"docker.elastic.co/elasticsearch/elasticsearch-oss:6.4.2\",\"Requests\":{\"RAM\":\"512Mi\",\"CPU\":\"25m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"1\"}}]},{\"PodName\":\"etcd-5cb7f54ccb-nd2zg\",\"Containers\":[{\"ContainerName\":\"etcd\",\"ContainerImage\":\"quay.io/coreos/etcd:v3.3.5\",\"Requests\":{\"RAM\":\"256M\",\"CPU\":\"250m\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"}}]},{\"PodName\":\"coredns-77cd44d8df-bsl5j\",\"Containers\":[{\"ContainerName\":\"coredns\",\"ContainerImage\":\"k8s.gcr.io/coredns:1.1.3\",\"Requests\":{\"RAM\":\"70Mi\",\"CPU\":\"100m\"},\"Limits\":{\"RAM\":\"170Mi\",\"CPU\":\"Is Not Set.\"}}]},{\"PodName\":\"heapster-v11-2j9kb\",\"Containers\":[{\"ContainerName\":\"heapster\",\"ContainerImage\":\"gcr.io/google_containers/heapster:v1.4.0\",\"Requests\":{\"RAM\":\"212Mi\",\"CPU\":\"100m\"},\"Limits\":{\"RAM\":\"212Mi\",\"CPU\":\"100m\"}}]},{\"PodName\":\"kube-proxy-ip-172-20-1-44.us-west-1.compute.internal\",\"Containers\":[{\"ContainerName\":\"kube-proxy\",\"ContainerImage\":\"gcr.io/google_containers/hyperkube:v1.11.1\",\"Requests\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"}}]},{\"PodName\":\"monitoring-influxdb-grafana-v2-vr9jh\",\"Containers\":[{\"ContainerName\":\"influxdb\",\"ContainerImage\":\"gcr.io/google_containers/heapster_influxdb:v0.4\",\"Requests\":{\"RAM\":\"200Mi\",\"CPU\":\"100m\"},\"Limits\":{\"RAM\":\"200Mi\",\"CPU\":\"100m\"}},{\"ContainerName\":\"grafana\",\"ContainerImage\":\"beta.gcr.io/google_containers/heapster_grafana:v2.1.1\",\"Requests\":{\"RAM\":\"100Mi\",\"CPU\":\"100m\"},\"Limits\":{\"RAM\":\"100Mi\",\"CPU\":\"100m\"}}]},{\"PodName\":\"tiller-deploy-677f9cb999-jv6ct\",\"Containers\":[{\"ContainerName\":\"tiller\",\"ContainerImage\":\"gcr.io/kubernetes-helm/tiller:v2.11.0\",\"Requests\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"},\"Limits\":{\"RAM\":\"Is Not Set.\",\"CPU\":\"Is Not Set.\"}}]}]}]}",
        "executionStatus": "OK",
        "id": "Resources (CPU/RAM) requests and limits Check",
        "name": "Resources (CPU/RAM) requests and limits Check",
        "possibleActions": [
          {
            "description": "Dismiss notification",
            "id": "1"
          },
          {
            "description": "Set missing requests/limits",
            "id": "2"
          }
        ]
      },
      {
        "checkStatus": "GREEN",
        "completedAt": "2018-11-14T11:37:39.269Z",
        "description": "[{\"Region\":\"us-west-1b\",\"InstanceID\":\"i-03fb8e89232700cc3\",\"InstanceType\":\"m4.large\",\"Name\":\"ip-172-20-1-44.us-west-1.compute.internal\",\"AllocatableCpu\":2000,\"AllocatableMemory\":8260214784,\"CpuReqs\":700,\"CpuLimits\":2300,\"MemoryReqs\":3013754880,\"MemoryLimits\":715128832,\"FractionCpuReqs\":35,\"FractionCpuLimits\":114.99999999999999,\"FractionMemoryReqs\":36.48518784084925,\"FractionMemoryLimits\":8.657508923196543},{\"Region\":\"us-west-1b\",\"InstanceID\":\"i-028bc20adaf2311d6\",\"InstanceType\":\"m4.large\",\"Name\":\"ip-172-20-1-21.us-west-1.compute.internal\",\"AllocatableCpu\":2000,\"AllocatableMemory\":8260218880,\"CpuReqs\":50,\"CpuLimits\":2000,\"MemoryReqs\":1073741824,\"MemoryLimits\":0,\"FractionCpuReqs\":2.5,\"FractionCpuLimits\":100,\"FractionMemoryReqs\":12.998951233602178,\"FractionMemoryLimits\":0},{\"Region\":\"us-west-1b\",\"InstanceID\":\"i-0898d927727329231\",\"InstanceType\":\"m4.large\",\"Name\":\"ip-172-20-1-242.us-west-1.compute.internal\",\"AllocatableCpu\":2000,\"AllocatableMemory\":8260218880,\"CpuReqs\":75,\"CpuLimits\":3000,\"MemoryReqs\":2684354560,\"MemoryLimits\":0,\"FractionCpuReqs\":3.75,\"FractionCpuLimits\":150,\"FractionMemoryReqs\":32.49737808400545,\"FractionMemoryLimits\":0}]",
        "executionStatus": "OK",
        "id": "Underutilized nodes sunsetting Check",
        "name": "Underutilized nodes sunsetting Check",
        "possibleActions": [
          {
            "description": "Dismiss notification",
            "id": "1"
          },
          {
            "description": "Sunset nodes",
            "id": "2"
          }
        ]
      }
    ];
    const mapJson = map((checks: any[]) => checks.map(check => {
      try {
        const description = JSON.parse(check.description);
        return {
          ...check, description
        }
      } catch (e) {
        //  is not json
        return check;
      }
    }));
    this.checks$ = this.http.get(apiV1Check).pipe(
      mapJson,
      // TODO for testing only
      catchError(_ => {
        return of(mock).pipe(mapJson);
      })
    );
  }

  isObject(val) {
    return typeof val === 'object';
  }


}
