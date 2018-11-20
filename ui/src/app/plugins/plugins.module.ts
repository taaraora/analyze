import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { PluginsRoutingModule } from './plugins-routing.module';
import { PluginsComponent } from './plugins.component';

@NgModule({
  declarations: [PluginsComponent],
  imports: [
    CommonModule,
    PluginsRoutingModule
  ]
})
export class PluginsModule { }
