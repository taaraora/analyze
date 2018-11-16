import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { HttpClient }                           from "@angular/common/http";
import { map }                                  from "rxjs/operators";
import { Observable }                           from "rxjs";

@Component({
  selector: 'app-checks',
  templateUrl: './checks.component.html',
  styleUrls: ['./checks.component.scss'],
  encapsulation: ViewEncapsulation.None,
})
export class ChecksComponent implements OnInit {

  checks$: Observable<any[]>;

  constructor(
    private http: HttpClient,
  ) {
  }

  ngOnInit() {
    const apiV1Check = 'http://ec2-54-241-150-199.us-west-1.compute.amazonaws.com:30759/api/v1/check';
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
    );
  }

  isObject(val) {
    return typeof val === 'object';
  }
}
