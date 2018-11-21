(window["webpackJsonp"] = window["webpackJsonp"] || []).push([["plugins-plugins-module"],{

/***/ "./src/app/plugins/plugins-routing.module.ts":
/*!***************************************************!*\
  !*** ./src/app/plugins/plugins-routing.module.ts ***!
  \***************************************************/
/*! exports provided: PluginsRoutingModule */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "PluginsRoutingModule", function() { return PluginsRoutingModule; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_router__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/router */ "./node_modules/@angular/router/fesm5/router.js");
/* harmony import */ var src_app_plugins_plugins_component__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! src/app/plugins/plugins.component */ "./src/app/plugins/plugins.component.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};



var routes = [
    {
        path: '',
        component: src_app_plugins_plugins_component__WEBPACK_IMPORTED_MODULE_2__["PluginsComponent"],
    }
];
var PluginsRoutingModule = /** @class */ (function () {
    function PluginsRoutingModule() {
    }
    PluginsRoutingModule = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["NgModule"])({
            imports: [_angular_router__WEBPACK_IMPORTED_MODULE_1__["RouterModule"].forChild(routes)],
            exports: [_angular_router__WEBPACK_IMPORTED_MODULE_1__["RouterModule"]]
        })
    ], PluginsRoutingModule);
    return PluginsRoutingModule;
}());



/***/ }),

/***/ "./src/app/plugins/plugins.component.html":
/*!************************************************!*\
  !*** ./src/app/plugins/plugins.component.html ***!
  \************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = "<div *ngFor=\"let plugin of (plugins | async)\">\n\n  <mat-card class=\"plugin transparent\">\n    <mat-card-title>\n      {{ plugin.name }}\n    </mat-card-title>\n    <mat-card-content>\n      {{ plugin.description }}\n    </mat-card-content>\n  </mat-card>\n</div>\n"

/***/ }),

/***/ "./src/app/plugins/plugins.component.scss":
/*!************************************************!*\
  !*** ./src/app/plugins/plugins.component.scss ***!
  \************************************************/
/*! no static exports found */
/***/ (function(module, exports) {

module.exports = ".plugin {\n  margin-bottom: 20px; }\n\n/*# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJzb3VyY2VzIjpbIi91c3Ivc3JjL2FwcC9zcmMvYXBwL3BsdWdpbnMvcGx1Z2lucy5jb21wb25lbnQuc2NzcyJdLCJuYW1lcyI6W10sIm1hcHBpbmdzIjoiQUFBQTtFQUNFLG9CQUFtQixFQUNwQiIsImZpbGUiOiJzcmMvYXBwL3BsdWdpbnMvcGx1Z2lucy5jb21wb25lbnQuc2NzcyIsInNvdXJjZXNDb250ZW50IjpbIi5wbHVnaW4ge1xuICBtYXJnaW4tYm90dG9tOiAyMHB4O1xufVxuIl19 */"

/***/ }),

/***/ "./src/app/plugins/plugins.component.ts":
/*!**********************************************!*\
  !*** ./src/app/plugins/plugins.component.ts ***!
  \**********************************************/
/*! exports provided: PluginsComponent */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "PluginsComponent", function() { return PluginsComponent; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var src_app_plugins_plugins_service__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! src/app/plugins/plugins.service */ "./src/app/plugins/plugins.service.ts");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var PluginsComponent = /** @class */ (function () {
    function PluginsComponent(pluginsService) {
        this.pluginsService = pluginsService;
    }
    PluginsComponent.prototype.ngOnInit = function () {
        this.plugins = this.pluginsService.all();
    };
    PluginsComponent = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Component"])({
            selector: 'app-plugins',
            template: __webpack_require__(/*! ./plugins.component.html */ "./src/app/plugins/plugins.component.html"),
            styles: [__webpack_require__(/*! ./plugins.component.scss */ "./src/app/plugins/plugins.component.scss")]
        }),
        __metadata("design:paramtypes", [src_app_plugins_plugins_service__WEBPACK_IMPORTED_MODULE_1__["PluginsService"]])
    ], PluginsComponent);
    return PluginsComponent;
}());



/***/ }),

/***/ "./src/app/plugins/plugins.module.ts":
/*!*******************************************!*\
  !*** ./src/app/plugins/plugins.module.ts ***!
  \*******************************************/
/*! exports provided: PluginsModule */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "PluginsModule", function() { return PluginsModule; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common */ "./node_modules/@angular/common/fesm5/common.js");
/* harmony import */ var _plugins_routing_module__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./plugins-routing.module */ "./src/app/plugins/plugins-routing.module.ts");
/* harmony import */ var _plugins_component__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./plugins.component */ "./src/app/plugins/plugins.component.ts");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
/* harmony import */ var src_app_plugins_plugins_service__WEBPACK_IMPORTED_MODULE_5__ = __webpack_require__(/*! src/app/plugins/plugins.service */ "./src/app/plugins/plugins.service.ts");
/* harmony import */ var _angular_material__WEBPACK_IMPORTED_MODULE_6__ = __webpack_require__(/*! @angular/material */ "./node_modules/@angular/material/esm5/material.es5.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};







var PluginsModule = /** @class */ (function () {
    function PluginsModule() {
    }
    PluginsModule = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["NgModule"])({
            declarations: [_plugins_component__WEBPACK_IMPORTED_MODULE_3__["PluginsComponent"]],
            imports: [
                _angular_common__WEBPACK_IMPORTED_MODULE_1__["CommonModule"],
                _plugins_routing_module__WEBPACK_IMPORTED_MODULE_2__["PluginsRoutingModule"],
                _angular_common_http__WEBPACK_IMPORTED_MODULE_4__["HttpClientModule"],
                _angular_material__WEBPACK_IMPORTED_MODULE_6__["MatCardModule"],
            ],
            providers: [
                src_app_plugins_plugins_service__WEBPACK_IMPORTED_MODULE_5__["PluginsService"],
            ]
        })
    ], PluginsModule);
    return PluginsModule;
}());



/***/ }),

/***/ "./src/app/plugins/plugins.service.ts":
/*!********************************************!*\
  !*** ./src/app/plugins/plugins.service.ts ***!
  \********************************************/
/*! exports provided: PluginsService */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
__webpack_require__.r(__webpack_exports__);
/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, "PluginsService", function() { return PluginsService; });
/* harmony import */ var _angular_core__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! @angular/core */ "./node_modules/@angular/core/fesm5/core.js");
/* harmony import */ var _angular_common_http__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! @angular/common/http */ "./node_modules/@angular/common/fesm5/http.js");
var __decorate = (undefined && undefined.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (undefined && undefined.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};


var PluginsService = /** @class */ (function () {
    function PluginsService(http) {
        this.http = http;
    }
    PluginsService.prototype.all = function () {
        return this.http.get('/api/v1//plugin');
    };
    PluginsService = __decorate([
        Object(_angular_core__WEBPACK_IMPORTED_MODULE_0__["Injectable"])(),
        __metadata("design:paramtypes", [_angular_common_http__WEBPACK_IMPORTED_MODULE_1__["HttpClient"]])
    ], PluginsService);
    return PluginsService;
}());



/***/ })

}]);
//# sourceMappingURL=plugins-plugins-module.js.map