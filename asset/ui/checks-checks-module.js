(window["webpackJsonp"] = window["webpackJsonp"] || []).push([["checks-checks-module"],{

/***/ "./src/app/checks/checks-routing.module.ts":
/*!*************************************************!*\
  !*** ./src/app/checks/checks-routing.module.ts ***!
  \*************************************************/
/*! exports provided: ChecksRoutingModule */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChecksRoutingModule", function() { return ChecksRoutingModule; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/router */ "./node_modules/@angular/router/fesm5/router.js");
/* harmony import */ var src_app_checks_checks_component__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! src/app/checks/checks.component */ "./src/app/checks/checks.component.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};



var routes = [
    {
        path: '',
        component: src_app_checks_checks_component__WEBPACK_IMPORTED_MODULE_2__["ChecksComponent"]
    },
];
var ChecksRoutingModule = /** @class */ (function () {
    function ChecksRoutingModule() {
    }
    ChecksRoutingModule = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["NgModule"])({
            imports: [_angular_router__WEBPACK_IMPORTED_MODULE_1__["RouterModule"].forChild(routes)],
            exports: [_angular_router__WEBPACK_IMPORTED_MODULE_1__["RouterModule"]]
        })
    ], ChecksRoutingModule);
    return ChecksRoutingModule;
}());



/***/ }),

/***/ "./src/app/checks/checks.component.html":
/*!**********************************************!*\
  !*** ./src/app/checks/checks.component.html ***!
  \**********************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<mat-card class=\"check-item transparent\" *ngFor=\"let check of (checks$ | async); let checkIndex = index\">\n  <mat-card-title>\n    <div> {{ check.name }}</div>\n    <section class=\"check-date\" [ngClass]=\"{\n      'green': check.checkStatus == 'GREEN',\n      'yellow': check.checkStatus == 'YELLOW',\n      'red': check.checkStatus == 'RED'\n    }\">\n      <div>\n        <div class=\"status-time\"> {{ check.completedAt | date: 'hh:mm:ss' }}</div>\n        <div class=\"status-time\"> {{ check.completedAt | date: 'MMM d, y'}}</div>\n      </div>\n      &nbsp;\n      <mat-icon [ngSwitch]=\"check.checkStatus\">\n        <span *ngSwitchCase=\"'RED'\">block</span>\n        <span *ngSwitchCase=\"'YELLOW'\">warning</span>\n        <span *ngSwitchCase=\"'GREEN'\">check_circle</span>\n      </mat-icon>\n    </section>\n  </mat-card-title>\n  <br>\n  <!-- TODO: remove harcode -->\n\n  <div *ngIf=\"check.description as data\">\n    <mat-accordion>\n      <mat-expansion-panel>\n        <mat-expansion-panel-header>\n          <mat-panel-title>\n            Details\n          </mat-panel-title>\n        </mat-expansion-panel-header>\n        <div class=\"scrollable-area\">\n          <ng-container *ngIf=\"checkIndex === 0\">\n            <table>\n              <tbody>\n              <tr>\n                <th>\n                  <div>Node Name</div>\n                </th>\n                <th>\n                  <div>Pod Resource Requirements</div>\n                </th>\n              </tr>\n              <tr *ngFor=\"let node of data.nodesResourceRequirements\">\n                <td>\n                  <div>{{ node.nodeName}}</div>\n                </td>\n                <td>\n                  <table>\n                    <tbody>\n                    <tr>\n                      <th>\n                        <div>Pod Name</div>\n                      </th>\n                      <th>\n                        <div>Containers Resource Requirements</div>\n                      </th>\n                    </tr>\n                    <tr *ngFor=\"let pod of node.podResourceRequirements\">\n                      <td>\n                        <div> {{pod.podName}}</div>\n                      </td>\n                      <td>\n                        <table>\n                          <tbody>\n                          <tr>\n                            <th>\n                              <div>Container Name</div>\n                            </th>\n                            <th>\n                              <div>Container Image</div>\n                            </th>\n                            <th>\n                              <div>Requests</div>\n                            </th>\n                            <th>\n                              <div>Limits</div>\n                            </th>\n                          </tr>\n                          <tr *ngFor=\"let container of pod.containersResourceRequirements\">\n                            <td>\n                              <div>{{ container.containerName }}</div>\n                            </td>\n                            <td>\n                              <div>{{ container.containerImage }}</div>\n                            </td>\n                            <td>\n                              <table class=\"requests\">\n                                <tbody>\n                                <tr>\n                                  <th>\n                                    <div>Ram</div>\n                                  </th>\n                                  <td>\n                                    <div [ngClass]=\"{warn: !container.requests.ram }\">\n                                      {{ container.requests.ram || IS_NOT_SET_MSG }}\n                                    </div>\n                                  </td>\n                                </tr>\n                                <tr>\n                                  <th>\n                                    <div>Cpu</div>\n                                  </th>\n                                  <td>\n                                    <div [ngClass]=\"{warn: !container.requests.cpu }\">\n                                      {{ container.requests.cpu || IS_NOT_SET_MSG }}\n                                    </div>\n                                  </td>\n                                </tr>\n                                </tbody>\n                              </table>\n                            </td>\n                            <td>\n                              <table class=\"limits\">\n                                <tbody>\n                                <tr>\n                                  <th>\n                                    <div>Ram</div>\n                                  </th>\n                                  <td>\n                                    <div [ngClass]=\"{warn: !container.limits.ram}\">\n                                      {{ container.limits.ram || IS_NOT_SET_MSG }}\n                                    </div>\n                                  </td>\n                                </tr>\n                                <tr>\n                                  <th>\n                                    <div>Cpu</div>\n                                  </th>\n                                  <td>\n                                    <div [ngClass]=\"{warn: !container.limits.cpu}\">\n                                      {{ container.limits.cpu || IS_NOT_SET_MSG }}\n                                    </div>\n                                  </td>\n                                </tr>\n                                </tbody>\n                              </table>\n                            </td>\n                          </tr>\n                          </tbody>\n                        </table>\n                      </td>\n                    </tr>\n                    </tbody>\n                  </table>\n                </td>\n              </tr>\n              </tbody>\n            </table>\n          </ng-container>\n          <ng-container *ngIf=\"checkIndex === 1\">\n            <table>\n              <tr>\n                <th>Region/Zone</th>\n                <th>Instance ID</th>\n                <th>RAM requested (GIB)</th>\n                <th>RAM not requested (GIB)</th>\n                <th>Total RAM (GIB)</th>\n                <th>Recommended to sunset</th>\n              </tr>\n\n              <tr *ngFor=\"let node of data\">\n                <td>\n                  {{ node.kube.region}}\n                </td>\n                <td>\n                  {{ node.cloudProvider.instanceId}}\n                </td>\n                <td>\n                  {{ node.kube.memoryRequests / 1000000000}}\n                </td>\n                <td>\n                  {{ (node.kube.allocatableMemory - node.kube.memoryRequests) / 1000000000 }}\n                </td>\n                <td>\n                  {{ node.kube.allocatableMemory / 1000000000}}\n                </td>\n                <td>\n                  {{node.kube.isRecommendedToSunset ? 'Yes' : 'No'}}\n                </td>\n              </tr>\n            </table>\n          </ng-container>\n        </div>\n      </mat-expansion-panel>\n    </mat-accordion>\n\n  </div>\n\n\n  <mat-tab-group>\n    <mat-tab *ngFor=\"let check of check.possibleActions as action\">\n      <ng-template mat-tab-label>\n        {{ check.name }}\n      </ng-template>\n      <div class=\"description\">\n        {{ check.description }}\n      </div>\n      <mat-chip-list class=\"action\">\n        <mat-chip class=\"red\">RUN</mat-chip>\n      </mat-chip-list>\n    </mat-tab>\n\n  </mat-tab-group>\n\n</mat-card>\n\n"

/***/ }),

/***/ "./src/app/checks/checks.component.scss":
/*!**********************************************!*\
  !*** ./src/app/checks/checks.component.scss ***!
  \**********************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ".scrollable-area {\n  max-height: 500px;\n  overflow-y: scroll;\n  margin-top: 12px; }\n\n.mat-tab-header .mat-ink-bar {\n  background-color: transparent !important;\n  border-bottom: 2px dotted white;\n  width: auto; }\n\n.mat-expansion-panel-header {\n  border: 1px solid rgba(255, 255, 255, 0.2); }\n\n.mat-tab-body-content {\n  padding: 20px 0; }\n\n.mat-card-title {\n  display: flex !important;\n  justify-content: space-between; }\n\n.description {\n  padding-bottom: 24px; }\n\n.action .red {\n  background: red; }\n\n.check-date {\n  font-size: 12px;\n  display: flex; }\n\n.check-date.red {\n    color: red; }\n\n.check-date.yellow {\n    color: yellow; }\n\n.check-date.green {\n    color: green; }\n\n.check-item {\n  margin: 20px; }\n\n.mat-expansion-panel {\n  background: transparent; }\n\ntable {\n  width: 100%;\n  height: 100%; }\n\ntable th {\n    background-color: rgba(255, 0, 166, 0.25);\n    border: none;\n    min-height: 45px; }\n\ntable tr {\n    background-color: rgba(255, 255, 255, 0.1); }\n\ntable td {\n    vertical-align: top;\n    padding: 5px; }\n\ntable td ~ table {\n      padding: 0; }\n\n.requests .warn {\n  color: red; }\n\n.limits .warn {\n  color: yellow; }\n\n/*# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIi91c3Ivc3JjL2FwcC9zcmMvYXBwL2NoZWNrcy9jaGVja3MuY29tcG9uZW50LnNjc3MiLCIvdXNyL3NyYy9hcHAvc3JjL2Fzc2V0cy9zdHlsZXMvY29sb3JzLnNjc3MiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IkFBSUE7RUFDRSxrQkFBaUI7RUFDakIsbUJBQWtCO0VBQ2xCLGlCQUFnQixFQUNqQjs7QUFFRDtFQUVJLHlDQUF3QztFQUN4QyxnQ0NYUztFRFlULFlBQVcsRUFDWjs7QUFHSDtFQUNFLDJDQUEwQyxFQUMzQzs7QUFFRDtFQUNFLGdCQUFlLEVBQ2hCOztBQUVEO0VBQ0UseUJBQXdCO0VBQ3hCLCtCQUE4QixFQUMvQjs7QUFFRDtFQUNFLHFCQUFvQixFQUNyQjs7QUFFRDtFQUVJLGdCQUFlLEVBQ2hCOztBQUdIO0VBQ0UsZ0JBQWU7RUFDZixjQUFhLEVBYWQ7O0FBZkQ7SUFLSSxXQUFVLEVBQ1g7O0FBTkg7SUFTSSxjQUFhLEVBQ2Q7O0FBVkg7SUFhSSxhQUFZLEVBQ2I7O0FBR0g7RUFDRSxhQUFZLEVBQ2I7O0FBVUQ7RUFDRSx3QkFBdUIsRUFDeEI7O0FBRUQ7RUFDRSxZQUFXO0VBQ1gsYUFBWSxFQW9CYjs7QUF0QkQ7SUFLSSwwQ0FBeUM7SUFDekMsYUFBWTtJQUNaLGlCQUFnQixFQUNqQjs7QUFSSDtJQVdJLDJDQUEwQyxFQUMzQzs7QUFaSDtJQWVJLG9CQUFtQjtJQUNuQixhQUFZLEVBS2I7O0FBckJIO01BbUJNLFdBQVUsRUFDWDs7QUFLTDtFQUVJLFdBQVUsRUFDWDs7QUFHSDtFQUVJLGNBQWEsRUFDZCIsImZpbGUiOiJzcmMvYXBwL2NoZWNrcy9jaGVja3MuY29tcG9uZW50LnNjc3MiLCJzb3VyY2VzQ29udGVudCI6WyJAaW1wb3J0IFwifnNyYy9hc3NldHMvc3R5bGVzL2NvbG9yc1wiO1xuXG4kdGFibGUtYm9yZGVyOiAxcHggc29saWQgd2hpdGU7XG5cbi5zY3JvbGxhYmxlLWFyZWEge1xuICBtYXgtaGVpZ2h0OiA1MDBweDtcbiAgb3ZlcmZsb3cteTogc2Nyb2xsO1xuICBtYXJnaW4tdG9wOiAxMnB4O1xufVxuXG4ubWF0LXRhYi1oZWFkZXIge1xuICAubWF0LWluay1iYXIge1xuICAgIGJhY2tncm91bmQtY29sb3I6IHRyYW5zcGFyZW50ICFpbXBvcnRhbnQ7XG4gICAgYm9yZGVyLWJvdHRvbTogMnB4IGRvdHRlZCAkd2hpdGU7XG4gICAgd2lkdGg6IGF1dG87XG4gIH1cbn1cblxuLm1hdC1leHBhbnNpb24tcGFuZWwtaGVhZGVyIHtcbiAgYm9yZGVyOiAxcHggc29saWQgcmdiYSgyNTUsIDI1NSwgMjU1LCAwLjIpO1xufVxuXG4ubWF0LXRhYi1ib2R5LWNvbnRlbnQge1xuICBwYWRkaW5nOiAyMHB4IDA7XG59XG5cbi5tYXQtY2FyZC10aXRsZSB7XG4gIGRpc3BsYXk6IGZsZXggIWltcG9ydGFudDtcbiAganVzdGlmeS1jb250ZW50OiBzcGFjZS1iZXR3ZWVuO1xufVxuXG4uZGVzY3JpcHRpb24ge1xuICBwYWRkaW5nLWJvdHRvbTogMjRweDtcbn1cblxuLmFjdGlvbiB7XG4gIC5yZWQge1xuICAgIGJhY2tncm91bmQ6IHJlZDtcbiAgfVxufVxuXG4uY2hlY2stZGF0ZSB7XG4gIGZvbnQtc2l6ZTogMTJweDtcbiAgZGlzcGxheTogZmxleDtcblxuICAmLnJlZCB7XG4gICAgY29sb3I6IHJlZDtcbiAgfVxuXG4gICYueWVsbG93IHtcbiAgICBjb2xvcjogeWVsbG93O1xuICB9XG5cbiAgJi5ncmVlbiB7XG4gICAgY29sb3I6IGdyZWVuO1xuICB9XG59XG5cbi5jaGVjay1pdGVtIHtcbiAgbWFyZ2luOiAyMHB4O1xufVxuXG5AbWl4aW4gYm9yZGVyLWJvdHRvbSB7XG4gIGJvcmRlci1ib3R0b206ICR0YWJsZS1ib3JkZXI7XG59XG5cbkBtaXhpbiBib3JkZXItcmlnaHQge1xuICBib3JkZXItcmlnaHQ6ICR0YWJsZS1ib3JkZXI7XG59XG5cbi5tYXQtZXhwYW5zaW9uLXBhbmVsIHtcbiAgYmFja2dyb3VuZDogdHJhbnNwYXJlbnQ7XG59XG5cbnRhYmxlIHtcbiAgd2lkdGg6IDEwMCU7XG4gIGhlaWdodDogMTAwJTtcblxuICB0aCB7XG4gICAgYmFja2dyb3VuZC1jb2xvcjogcmdiYSgyNTUsIDAsIDE2NiwgMC4yNSk7XG4gICAgYm9yZGVyOiBub25lO1xuICAgIG1pbi1oZWlnaHQ6IDQ1cHg7XG4gIH1cblxuICB0ciB7XG4gICAgYmFja2dyb3VuZC1jb2xvcjogcmdiYSgyNTUsIDI1NSwgMjU1LCAwLjEpO1xuICB9XG5cbiAgdGQge1xuICAgIHZlcnRpY2FsLWFsaWduOiB0b3A7XG4gICAgcGFkZGluZzogNXB4O1xuXG4gICAgJn50YWJsZSB7XG4gICAgICBwYWRkaW5nOiAwO1xuICAgIH1cbiAgfVxufVxuXG5cbi5yZXF1ZXN0cyB7XG4gIC53YXJuIHtcbiAgICBjb2xvcjogcmVkO1xuICB9XG59XG5cbi5saW1pdHMge1xuICAud2FybiB7XG4gICAgY29sb3I6IHllbGxvdztcbiAgfVxufVxuIiwiJGJsYWNrLXRyYW5zcGFyZW50OiByZ2JhKDAsIDAsIDAsIDAuNyk7XG4kc2ctcGluazogI0QxMDA4ODtcbiR3aGl0ZTogd2hpdGU7XG4iXX0= */"

/***/ }),

/***/ "./src/app/checks/checks.component.ts":
/*!********************************************!*\
  !*** ./src/app/checks/checks.component.ts ***!
  \********************************************/
/*! exports provided: ChecksComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChecksComponent", function() { return ChecksComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var rxjs_operators__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! rxjs/operators */ "./node_modules/rxjs/_esm5/operators/index.js");
var __assign = (undefined && undefined.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};



var ChecksComponent = /** @class */ (function () {
    function ChecksComponent(http) {
        this.http = http;
        this.IS_NOT_SET_MSG = 'Is not set';
    }
    ChecksComponent.prototype.ngOnInit = function () {
        var apiV1Check = '/api/v1/check';
        var mapJson = Object(rxjs_operators__WEBPACK_IMPORTED_MODULE_2__["map"])(function (checks) { return checks.map(function (check) {
            try {
                var description = JSON.parse(check.description);
                return __assign({}, check, { description: description });
            }
            catch (e) {
                //  is not json
                return check;
            }
        }); });
        this.checks$ = this.http.get(apiV1Check).pipe(mapJson);
    };
    ChecksComponent.prototype.isObject = function (val) {
        return typeof val === 'object';
    };
    ChecksComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-checks',
            template: __webpack_require__(/*! ./checks.component.html */ "./src/app/checks/checks.component.html"),
            styles: [__webpack_require__(/*! ./checks.component.scss */ "./src/app/checks/checks.component.scss")],
            encapsulation: _angular_core__WEBPACK_IMPORTED_MODULE_0__["ViewEncapsulation"].None,
        }),
        __metadata("design:paramtypes", [_angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpClient"]])
    ], ChecksComponent);
    return ChecksComponent;
}());



/***/ }),

/***/ "./src/app/checks/checks.module.ts":
/*!*****************************************!*\
  !*** ./src/app/checks/checks.module.ts ***!
  \*****************************************/
/*! exports provided: ChecksModule */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "ChecksModule", function() { return ChecksModule; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common */ "./node_modules/@angular/common/fesm5/common.js");
/* harmony import */ var _checks_routing_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./checks-routing.module */ "./src/app/checks/checks-routing.module.ts");
/* harmony import */ var _checks_component__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./checks.component */ "./src/app/checks/checks.component.ts");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var _angular_material__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! @angular/material */ "./node_modules/@angular/material/esm5/material.es5.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};






var ChecksModule = /** @class */ (function () {
    function ChecksModule() {
    }
    ChecksModule = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["NgModule"])({
            declarations: [_checks_component__WEBPACK_IMPORTED_MODULE_3__["ChecksComponent"]],
            imports: [
                _angular_common__WEBPACK_IMPORTED_MODULE_1__["CommonModule"],
                _checks_routing_module__WEBPACK_IMPORTED_MODULE_2__["ChecksRoutingModule"],
                _angular_common_http__WEBPACK_IMPORTED_MODULE_4__["HttpClientModule"],
                _angular_material__WEBPACK_IMPORTED_MODULE_5__["MatCardModule"],
                _angular_material__WEBPACK_IMPORTED_MODULE_5__["MatTabsModule"],
                _angular_material__WEBPACK_IMPORTED_MODULE_5__["MatExpansionModule"],
                _angular_material__WEBPACK_IMPORTED_MODULE_5__["MatIconModule"],
                _angular_material__WEBPACK_IMPORTED_MODULE_5__["MatChipsModule"],
            ],
        })
    ], ChecksModule);
    return ChecksModule;
}());



/***/ })

}]);
//# sourceMappingURL=checks-checks-module.js.map