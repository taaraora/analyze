import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { PluginsRoutingModule } from './plugins-routing.module';
import { PluginsComponent }     from './plugins.component';
import { HttpClientModule }     from "@angular/common/http";
import { PluginsService }       from "src/app/plugins/plugins.service";
import { MatCardModule }        from "@angular/material";

@NgModule({
  declarations: [PluginsComponent],
  imports: [
    CommonModule,
    PluginsRoutingModule,
    HttpClientModule,
    MatCardModule,
  ],
  providers: [
    PluginsService,
  ]
})
export class PluginsModule { }
